package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/albertopformoso/graphql-todo-list/config"
	"github.com/albertopformoso/graphql-todo-list/database"
	"github.com/albertopformoso/graphql-todo-list/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e = echo.New()

func main() {
	config := config.GetConfig()
	ctx := context.TODO()

	db := database.ConnectDB(ctx, config.Mongo)
	collection := db.Collection(config.Mongo.Collection)

	client := &database.TodoClient{
		Collection: collection,
		Ctx:        ctx,
	}
	fmt.Println(client)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes := e.Group("/v1")
	routes.GET("", func(c echo.Context) error {
		return c.JSONPretty(
			http.StatusOK, 
			map[string]string{"GraphQL": "To Do List"}, 
			"  ",
		)
	})
	routes.POST("/graphql", handler.Todos(client))

	e.Start(":8080")
}
