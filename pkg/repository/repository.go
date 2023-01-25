package repository

import (
	"github.com/Anchousfish/rest/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type TodoList interface {
	CreateList(userId int, list models.TodoList) (int, error)
	GetAllLists(userId int) ([]models.TodoList, error)
	GetList(listId, userId int) (models.TodoList, error)
	DeleteList(listId, userId int) error
	UpdateList(listId int, userId int, list models.UpdateListInput) error
}

type TodoItem interface {
	CreateItem(listId int, item models.TodoItem) (int, error)
	GetAllItems(userId, listId int) ([]models.TodoItem, error)
	GetItem(userId, itemId int) (models.TodoItem, error)
	DeleteItem(userId, itemId int) error
	UpdateItem(userId, itemId int, item models.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
