package repository

import (
	"fmt"
	"strings"

	"github.com/Anchousfish/rest/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (l *TodoListPostgres) CreateList(userId int, list models.TodoList) (int, error) {
	tr, err := l.db.Begin()
	if err != nil {
		return 0, nil
	}
	var id int
	createListQuery := fmt.Sprintf("insert into %s(title, description) values ($1,$2) returning id", todoListsTable)
	err = tr.QueryRow(createListQuery, list.Title, list.Description).Scan(&id)
	if err != nil {
		err = tr.Rollback()
		return 0, err
	}
	createUsersListsQuery := fmt.Sprintf("insert into %s(user_id, list_id) values ($1,$2)", userListsTable)
	_, err = tr.Exec(createUsersListsQuery, userId, id)
	if err != nil {
		err = tr.Rollback()
		return 0, err
	}

	return id, tr.Commit()

}

func (l *TodoListPostgres) GetAllLists(userId int) ([]models.TodoList, error) {
	var lists []models.TodoList
	getLists := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl inner join %s tu on tl.id=tu.list_id where tu.user_id=$1", todoListsTable, userListsTable)
	if err := l.db.Select(&lists, getLists, userId); err != nil {
		return nil, err
	}
	return lists, nil

}

func (l *TodoListPostgres) GetList(listId, userId int) (models.TodoList, error) {
	var list models.TodoList
	getList := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl inner join %s tu on tl.id=tu.list_id where tu.user_id=$1 and tu.list_id=$2", todoListsTable, userListsTable)
	err := l.db.Get(&list, getList, userId, listId)
	return list, err
}

func (l *TodoListPostgres) DeleteList(listId, userId int) error {
	deleteList := fmt.Sprintf("delete from  %s tl using %s tu where tl.id=tu.list_id and tu.user_id=$1 and tu.list_id=$2", todoListsTable, userListsTable)
	_, err := l.db.Exec(deleteList, userId, listId)
	return err
}

func (l *TodoListPostgres) UpdateList(listId int, userId int, list models.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if list.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *list.Title)
		argId++
	}
	if list.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *list.Description)
		argId++
	}

	query := strings.Join(setValues, ", ")

	args = append(args, userId, listId)
	updateList := fmt.Sprintf("update %s tl set %s from %s tu  where tl.id=tu.list_id and tu.user_id=$%d and tu.list_id=$%d", todoListsTable, query, userListsTable, argId, argId+1)
	logrus.Debugf("query: %s ", updateList)
	logrus.Debugf("args: %s", args...)
	_, err := l.db.Exec(updateList, args...)
	return err
}
