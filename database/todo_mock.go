package database

import (
	"github.com/albertopformoso/graphql-todo-list/model"
	"github.com/stretchr/testify/mock"
)

type MockTodoClient struct {
	mock.Mock
}

func (m *MockTodoClient) Insert(todo model.Todo) (model.Todo, error) {
	args := m.Called(todo)
	return args.Get(0).(model.Todo), args.Error(1)
}

func (m *MockTodoClient) Update(id string, update interface{}) (model.TodoUpdate, error) {
	args := m.Called(id, update)
	return args.Get(0).(model.TodoUpdate), args.Error(1)
}

func (m *MockTodoClient) Delete(id string) (model.TodoDelete, error) {
	args := m.Called(id)
	return args.Get(0).(model.TodoDelete), args.Error(1)
}

func (m *MockTodoClient) Get(id string) (model.Todo, error) {
	args := m.Called(id)
	return args.Get(0).(model.Todo), args.Error(1)
}

func (m *MockTodoClient) Search(filter interface{}) ([]model.Todo, error) {
	args := m.Called(filter)
	return args.Get(0).([]model.Todo), args.Error(1)
}