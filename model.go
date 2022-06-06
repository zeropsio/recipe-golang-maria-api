package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"

	_ "embed"
)

//go:embed schema.sql
var migration string

type Todo struct {
	Id        int    `json:"id" db:"id"`
	Completed bool   `json:"completed" db:"completed"`
	Text      string `json:"text" db:"text"`
}

type UpdateTodo struct {
	Completed *bool  `json:"completed" db:"completed"`
	Text      string `json:"text" db:"text"`
}

type TodoRepository struct {
	conn *sqlx.DB
}

func (t TodoRepository) FindOne(ctx context.Context, id int) (Todo, bool, error) {
	var todo Todo
	err := t.conn.GetContext(ctx, &todo, "SELECT id, completed, text FROM todos WHERE id=?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return todo, false, nil
	}
	if err != nil {
		return todo, false, err
	}
	return todo, true, nil
}

func (t TodoRepository) FindAll(ctx context.Context) ([]Todo, error) {
	var todos []Todo
	err := t.conn.SelectContext(ctx, &todos, `SELECT id, completed, text FROM todos`)
	return todos, err
}

func (t TodoRepository) Create(ctx context.Context, todo Todo) (Todo, error) {
	var id int64
	res, err := t.conn.NamedExecContext(ctx,
		"INSERT INTO todos(completed, text) VALUES (:completed, :text)", todo)
	if err != nil {
		return Todo{}, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return Todo{}, err
	}
	todo.Id = int(id)
	return todo, nil
}

func (t TodoRepository) Edit(ctx context.Context, id int, updateTodo UpdateTodo) (Todo, error) {
	oldTodo, found, err := t.FindOne(ctx, id)
	if err != nil {
		return oldTodo, err
	}
	if !found {
		return oldTodo, sql.ErrNoRows
	}
	if updateTodo.Completed != nil {
		oldTodo.Completed = *updateTodo.Completed
	}
	if updateTodo.Text != "" {
		oldTodo.Text = updateTodo.Text
	}
	_, err = t.conn.NamedExecContext(ctx, "UPDATE todos SET completed=:completed, text=:text WHERE id=:id", oldTodo)
	return oldTodo, err
}

func (t TodoRepository) Delete(ctx context.Context, id int) error {
	_, err := t.conn.ExecContext(ctx, "DELETE FROM todos WHERE id=?", id)
	return err
}

func (t TodoRepository) PrepareDatabase(ctx context.Context, dropTable bool, seeds []string) error {
	if dropTable {
		_, err := t.conn.ExecContext(ctx, "DROP TABLE IF EXISTS todos")
		if err != nil {
			return err
		}
	}
	_, err := t.conn.ExecContext(ctx, migration)

	if dropTable {
		for _, seed := range seeds {
			_, err := t.Create(ctx, Todo{Text: seed})
			if err != nil {
				return err
			}
		}
	}

	return err
}
