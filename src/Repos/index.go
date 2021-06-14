package Repos

import (
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user User) (int, error)
}

type Repo struct {
	Authorization
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		Authorization: NewAuthPostgresRepo(db),
	}
}
