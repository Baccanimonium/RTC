package Repos

import "github.com/jmoiron/sqlx"

type User struct {
	id       int    `json:"-" db:"id"`
	name     string `json:"name" binding:"required"`
	login    string `json:"login" binding:"required"`
	password string `json:"password" binding:"required"`
	about    string `json:"about"`
	address  string `json:"address"`
	phone    string `json:"phone"`
}

type Auth interface {
	CreateUser(user User)
}

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgresRepo(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user User) (int, error) {

	return 0, nil
}
