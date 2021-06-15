package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	About    string `json:"about"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
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
	var id int
	query := fmt.Sprintf(
		"INSERT INTO %s (name, login, password_hash, about, address, phone) values ($1, $2, $3, $4, $5, $6) RETURNING id",
		usersTable,
	)
	row := r.db.QueryRow(query, user.Name, user.Login, user.Password, user.About, user.Address, user.Phone)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
