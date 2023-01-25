package service

import (
	"github.com/Anchousfish/rest/models"
	"github.com/Anchousfish/rest/pkg/repository"
)

type TodoItemService struct {
	itemRepo repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(item repository.TodoItem, list repository.TodoList) *TodoItemService {
	return &TodoItemService{itemRepo: item, listRepo: list}
}

func (i *TodoItemService) CreateItem(userId int, listId int, item models.TodoItem) (int, error) {
	_, err := i.listRepo.GetList(listId, userId)
	if err != nil {
		//list does not exists or does not belong to this user
		return 0, err
	}
	return i.itemRepo.CreateItem(listId, item)
}
func (i *TodoItemService) GetAllItems(uid, listId int) ([]models.TodoItem, error) {
	return i.itemRepo.GetAllItems(uid, listId)
}

func (i *TodoItemService) GetItem(userId, itemId int) (models.TodoItem, error) {
	return i.itemRepo.GetItem(userId, itemId)
}
func (i *TodoItemService) DeleteItem(userId, itemId int) error {
	return i.itemRepo.DeleteItem(userId, itemId)
}

func (i *TodoItemService) UpdateItem(userId, itemId int, item models.UpdateItemInput) error {
	return i.itemRepo.UpdateItem(userId, itemId, item)
}
