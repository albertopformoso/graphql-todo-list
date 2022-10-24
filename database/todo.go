package database

import (
	"context"
	"encoding/json"

	"github.com/albertopformoso/graphql-todo-list/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todoer interface {
	Insert(model.Todo) (model.Todo, error)
	Update(string, interface{}) (model.TodoUpdate, error)
	Delete(string) (model.TodoDelete, error)
	Get(string) (model.Todo, error)
	Search(interface{}) ([]model.Todo, error)
}

type TodoClient struct {
	Ctx        context.Context
	Collection *mongo.Collection
}

func (c *TodoClient) Insert(docs model.Todo) (model.Todo, error) {
	var todo model.Todo

	res, err := c.Collection.InsertOne(c.Ctx, docs)
	if err != nil {
		return todo, err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()

	return c.Get(id)
}

func (c *TodoClient) Update(id string, update interface{}) (model.TodoUpdate, error) {
	result := model.TodoUpdate{
		ModifiedCount: 0,
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	todo, err := c.Get(id)
	if err != nil {
		return result, err
	}

	result.Result = todo
	var exist map[string]interface{}
	b, err := json.Marshal(todo)
	if err != nil {
		return result, err
	}
	json.Unmarshal(b, &exist)

	change := update.(map[string]interface{})
	for k := range change {
		if change[k] == exist[k] {
			delete(change, k)
		}
	}

	if len(change) == 0 {
		return result, nil
	}

	res, err := c.Collection.UpdateOne(c.Ctx, bson.M{"_id": _id}, bson.M{"$set": change})
	if err != nil {
		return result, err
	}

	newTodo, err := c.Get(id)
	if err != nil {
		return result, err
	}

	result.ModifiedCount = res.ModifiedCount
	result.Result = newTodo
	return result, nil
}

func (c *TodoClient) Delete(id string) (model.TodoDelete, error) {
	result := model.TodoDelete{
		DeletedCount: 0,
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	res, err := c.Collection.DeleteOne(c.Ctx, bson.M{"_id": _id})
	if err != nil {
		return result, err
	}

	result.DeletedCount = res.DeletedCount
	return result, nil
}

func (c *TodoClient) Get(id string) (model.Todo, error) {
	var todo model.Todo

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return todo, err
	}

	if err := c.Collection.FindOne(c.Ctx, bson.M{"_id": _id}).Decode(&todo); err != nil {
		return todo, err
	}

	todo.ID = todo.ID.(primitive.ObjectID).Hex()

	return todo, nil
}

func (c *TodoClient) Search(filter interface{}) ([]model.Todo, error) {
	var todos []model.Todo
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := c.Collection.Find(c.Ctx, filter)
	if err != nil {
		return todos, err
	}

	for cursor.Next(c.Ctx) {
		var row model.Todo
		cursor.Decode(&row)
		row.ID = row.ID.(primitive.ObjectID).Hex()
		todos = append(todos, row)
	}

	return todos, nil
}
