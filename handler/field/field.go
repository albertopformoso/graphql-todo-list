package field

import (
	"github.com/graphql-go/graphql"
)

var TodoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TodoType",
	Fields: graphql.Fields{
		"userId":    &graphql.Field{Type: graphql.Int},
		"id":        &graphql.Field{Type: graphql.String},
		"title":     &graphql.Field{Type: graphql.String},
		"completed": &graphql.Field{Type: graphql.Boolean},
	},
})

var UpdateTodoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UpdateTodoType",
	Fields: graphql.Fields{
		"modifiedCount": &graphql.Field{Type: graphql.Int},
		"result":        &graphql.Field{Type: TodoType},
	},
})

var DeleteTodoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DeleteTodoType",
	Fields: graphql.Fields{
		"deletedCount": &graphql.Field{Type: graphql.Int},
	},
})
