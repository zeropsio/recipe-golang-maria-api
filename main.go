package main

import (
	"context"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const DataSeed = "ZEROPS_RECIPE_DATA_SEED"
const DropTable = "ZEROPS_RECIPE_DROP_TABLE"
const DbUrl = "DB_URL"

func main() {
	ctx := context.Background()

	dbUrl, ok := os.LookupEnv(DbUrl)
	if !ok {
		panic("database url missing set " + DbUrl + " env")
	}

	conn, err := sqlx.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	seeds, err := getSeeds()
	if err != nil {
		panic(err)
	}

	dropTable, err := getDropTable()
	if err != nil {
		panic(err)
	}

	model := TodoRepository{conn}
	err = model.PrepareDatabase(ctx, dropTable, seeds)
	if err != nil {
		panic(err)
	}

	handler := todoHandler{model}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())

	cors := middleware.DefaultCORSConfig
	cors.AllowMethods = append(cors.AllowMethods, http.MethodPatch)
	e.Use(middleware.CORSWithConfig(cors))

	g := e.Group("todos")
	g.GET("", handler.getTodos)
	g.GET("/:id", handler.getTodo)
	g.POST("", handler.createTodo)
	g.PATCH("/:id", handler.editTodo)
	g.DELETE("/:id", handler.deleteTodo)

	e.Logger.Fatal(e.Start(":3000"))
}

func getSeeds() ([]string, error) {
	dbSeed, ok := os.LookupEnv(DataSeed)
	if !ok {
		dbSeed = "[]"
	}
	var seeds []string
	err := json.Unmarshal([]byte(dbSeed), &seeds)
	return seeds, err
}

func getDropTable() (bool, error) {
	dropTable, ok := os.LookupEnv(DropTable)
	if !ok {
		return false, nil
	}
	return strconv.ParseBool(dropTable)
}
