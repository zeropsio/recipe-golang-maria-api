package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type todoHandler struct {
	model TodoRepository
}

func (t todoHandler) getTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(400, "id should be int")
	}
	todo, found, err := t.model.FindOne(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "todo not found")
	}
	return c.JSON(http.StatusOK, todo)
}

func (t todoHandler) getTodos(c echo.Context) error {
	todos, err := t.model.FindAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, todos)
}

func (t todoHandler) createTodo(c echo.Context) error {
	var todo Todo
	err := c.Bind(&todo)
	if err != nil {
		return err
	}
	todo, err = t.model.Create(c.Request().Context(), todo)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, todo)
}

func (t todoHandler) editTodo(c echo.Context) error {
	var updateTodo UpdateTodo
	err := c.Bind(&updateTodo)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(400, "id should be int")
	}
	todo, err := t.model.Edit(c.Request().Context(), id, updateTodo)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, todo)
}

func (t todoHandler) deleteTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(400, "id should be int")
	}
	err = t.model.Delete(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}
