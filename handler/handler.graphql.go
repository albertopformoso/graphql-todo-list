package handler

import (
	"encoding/json"
	"net/http"

	"github.com/albertopformoso/graphql-todo-list/database"
	"github.com/albertopformoso/graphql-todo-list/handler/field"
	"github.com/albertopformoso/graphql-todo-list/model"
	"github.com/graphql-go/graphql"
	"github.com/labstack/echo/v4"
)

// Handler ==========

func Todos(db database.Todoer) echo.HandlerFunc {
	return func(c echo.Context) error {
		rootQuery := graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"searchTodos": searchTodos(db),
				"getTodo":     getTodo(db),
			},
		})

		rootMutation := graphql.NewObject(graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"insertTodo": insertTodo(db),
				"updateTodo": updateTodo(db),
				"deleteTodo": deleteTodo(db),
			},
		})

		schema, _ := graphql.NewSchema(graphql.SchemaConfig{
			Query:    rootQuery,
			Mutation: rootMutation,
		})

		// get query string
		requestString := c.QueryParam("q")
		// get request body
		if requestString == "" {
			var body map[string]interface{}
			c.Bind(&body)
			requestString = body["query"].(string)
		}

		res := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: requestString,
		})

		return c.JSONPretty(http.StatusOK, res, "  ")
	}
}

// Service ==========

func searchTodos(db database.Todoer) *graphql.Field {
	args := graphql.FieldConfigArgument{
		"title":     &graphql.ArgumentConfig{Type: graphql.String},
		"completed": &graphql.ArgumentConfig{Type: graphql.Boolean},
		"userId":    &graphql.ArgumentConfig{Type: graphql.Int},
	}

	return &graphql.Field{
		Name:        "Todos",
		Description: "List of To Do's",
		Type:        graphql.NewList(field.TodoType),
		Args:        args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			filter := p.Args
			return db.Search(filter)
		},
	}
}

func getTodo(db database.Todoer) *graphql.Field {
	args := graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.String},
	}

	return &graphql.Field{
		Name:        "Todo",
		Description: "Get To Do by Id",
		Type:        field.TodoType,
		Args:        args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := p.Args["id"].(string)
			return db.Get(id)
		},
	}
}

func insertTodo(db database.Todoer) *graphql.Field {
	args := graphql.FieldConfigArgument{
		"userId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"title": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"completed": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Boolean),
		},
	}

	return &graphql.Field{
		Name:        "Insert To Do",
		Type:        field.TodoType,
		Description: "Insert To Do Item",
		Args:        args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var todo model.Todo
			body, err := json.Marshal(p.Args)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(body, &todo); err != nil {
				return "", err
			}

			return db.Insert(todo)
		},
	}
}

func updateTodo(db database.Todoer) *graphql.Field {
	args := graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"userId": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"title": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"completed": &graphql.ArgumentConfig{
			Type: graphql.Boolean,
		},
	}

	return &graphql.Field{
		Name:        "Update Todo",
		Description: "Update To Do Item by ID",
		Type:        field.UpdateTodoType,
		Args:        args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			args := p.Args
			id := args["id"].(string)
			delete(args, "id")
			return db.Update(id, args)
		},
	}
}

func deleteTodo(db database.Todoer) *graphql.Field {
	args := graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	}

	return &graphql.Field{
		Name:        "Delete Todo",
		Description: "Delete Todo Item by Id",
		Type:        field.DeleteTodoType,
		Args:        args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := p.Args["id"].(string)
			return db.Delete(id)
		},
	}
}
