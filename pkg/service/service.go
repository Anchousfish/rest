package service

import (
	"github.com/Anchousfish/rest/models"
	"github.com/Anchousfish/rest/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (userid int, err error)
}

type TodoList interface {
	CreateList(userId int, list models.TodoList) (int, error)
	GetAllLists(userId int) ([]models.TodoList, error)
	GetList(listId, userId int) (models.TodoList, error)
	DeleteList(listId, userId int) error
	UpdateList(listId int, userId int, list models.UpdateListInput) error
}

type TodoItem interface {
	CreateItem(userId int, listId int, item models.TodoItem) (int, error)
	GetAllItems(uid, listId int) ([]models.TodoItem, error)
	GetItem(userId, itemId int) (models.TodoItem, error)
	DeleteItem(userId, itemId int) error
	UpdateItem(userId, itemId int, item models.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(r.Authorization),
		TodoList:      NewTodoListService(r.TodoList),
		TodoItem:      NewTodoItemService(r.TodoItem, r.TodoList),
	}
}
