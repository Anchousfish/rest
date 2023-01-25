package repository

import (
	"fmt"
	"strings"

	"github.com/Anchousfish/rest/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (idb *TodoItemPostgres) CreateItem(listId int, item models.TodoItem) (int, error) {
	tr, err := idb.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "failed to create a transaction")
	}
	var itemId int
	createItemQuery := fmt.Sprintf("insert into %s(title, description) values ($1, $2) returning id", todoItemsTable)
	row := tr.QueryRow(createItemQuery, item.Title, item.Description)
	if err = row.Scan(&itemId); err != nil {
		tr.Rollback()
		return 0, errors.Wrap(err, "failed to insert new item")
	}
	createListItemQuery := fmt.Sprintf("insert into %s(list_id,item_id) values ($1, $2) returning id", listsItemsTable)
	_, err = tr.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		tr.Rollback()
		return 0, errors.Wrap(err, "failed to link list and item")
	}

	return itemId, tr.Commit()
}
func (idb *TodoItemPostgres) GetAllItems(userId, listId int) ([]models.TodoItem, error) {
	var items []models.TodoItem
	query := fmt.Sprintf(`select ti.id, ti.title, ti.description from %s ti
							inner join %s li on li.item_id=ti.id
							inner join %s ul on ul.list_id=li.list_id
							where li.list_id=$1 and ul.user_id=$2`, todoItemsTable, listsItemsTable, userListsTable)
	if err := idb.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}
	return items, nil
}

func (idb *TodoItemPostgres) GetItem(userId, itemId int) (models.TodoItem, error) {
	var item models.TodoItem
	query := fmt.Sprintf(`select ti.id, ti.title, ti.description from %s ti
	inner join %s li on li.item_id=ti.id
	inner join %s ul on ul.list_id=li.list_id
	where ti.id=$1 and ul.user_id=$2`, todoItemsTable, listsItemsTable, userListsTable)
	err := idb.db.Get(&item, query, itemId, userId)
	return item, err
}

func (idb *TodoItemPostgres) DeleteItem(userId, itemId int) error {

	query := fmt.Sprintf(`delete from %s ti using %s li, %s ul where ti.id=li.item_id and li.list_id_id=ul.list_id and  ti.id=$1 and ul.user_id=$2`, todoItemsTable, listsItemsTable, userListsTable)
	_, err := idb.db.Exec(query, itemId, userId)
	return err
}

func (idb *TodoItemPostgres) UpdateItem(userId, itemId int, item models.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if item.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *item.Title)
		argId++
	}
	if item.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *item.Description)
		argId++
	}

	if item.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *item.Done)
		argId++
	}

	query := strings.Join(setValues, ", ")

	args = append(args, itemId, userId)
	updateList := fmt.Sprintf(`update %s ti set %s from %s li, %s ul
								where ti.id=li.item_id and li.list_id=ul.list_id and ti.id=$%d and ul.user_id=$%d`,
		todoItemsTable, query, listsItemsTable, userListsTable, argId, argId+1)
	logrus.Debugf("query: %s ", updateList)
	logrus.Debugf("args: %s", args...)
	_, err := idb.db.Exec(updateList, args...)
	return err
}
