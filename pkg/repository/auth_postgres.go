package repository

import (
	"fmt"

	"github.com/Anchousfish/rest/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user models.User) (int, error) {
	id := 0
	query := fmt.Sprintf("Insert into %s(name, username, password_hash) values ($1, $2,$3) returning id", usersTable)
	row := a.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (a *AuthPostgres) GetUser(username, pass string) (models.User, error) {
	var user = &models.User{}
	query := fmt.Sprintf("select id from %s where username=$1 and password_hash=$2", usersTable)
	err := a.db.Get(user, query, username, pass)
	return *user, err
}
