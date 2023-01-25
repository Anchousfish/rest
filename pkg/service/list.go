package service

import (
	"github.com/Anchousfish/rest/models"
	"github.com/Anchousfish/rest/pkg/repository"
)

type TodoListService struct {
	repos repository.TodoList
}

func NewTodoListService(repos repository.TodoList) *TodoListService {
	return &TodoListService{repos: repos}
}

func (l *TodoListService) CreateList(userId int, list models.TodoList) (int, error) {
	return l.repos.CreateList(userId, list)
}

func (l *TodoListService) GetAllLists(userId int) ([]models.TodoList, error) {
	return l.repos.GetAllLists(userId)
}
func (l *TodoListService) GetList(listId, userId int) (models.TodoList, error) {
	return l.repos.GetList(listId, userId)
}
func (l *TodoListService) DeleteList(listId, userId int) error {
	return l.repos.DeleteList(listId, userId)
}

func (l *TodoListService) UpdateList(listId int, userId int, list models.UpdateListInput) error {
	if err := list.Validate(); err != nil {
		return err
	}
	return l.repos.UpdateList(listId, userId, list)
}
