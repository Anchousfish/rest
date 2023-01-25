package models

import "errors"

type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}

type UpdateListInput struct {
	Title       *string `json:"title" db:"title"`
	Description *string `json:"description" db:"description"`
}

func (u UpdateListInput) Validate() error {
	if u.Description == nil && u.Title == nil {
		return errors.New("update structure has no values")
	}
	return nil
}

type UpdateItemInput struct {
	Title       *string `json:"title" db:"title"`
	Description *string `json:"description" db:"description"`
	Done        *bool   `json:"done" db:"done"`
}

func (item UpdateItemInput) Validate() error {
	if item.Description == nil && item.Title == nil && item.Done == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
