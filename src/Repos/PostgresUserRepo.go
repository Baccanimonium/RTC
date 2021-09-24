package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewUserPostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) GetAllUser() ([]UserCreate, error) {
	var users []UserCreate

	query := fmt.Sprintf(`SELECT id, name, login, password, about, address, phone, surname, avatar FROM %s`, usersTable)

	err := r.db.Select(&users, query)

	return users, err
}

func (r *Postgres) GetUserById(id int) (UserCreate, error) {
	var user UserCreate

	query := fmt.Sprintf(`SELECT
		id, name, login, password, about, address, phone, surname, avatar FROM %s WHERE id = $1`,
		usersTable,
	)

	err := r.db.Get(&user, query, id)

	return user, err
}

func (r *Postgres) UpdateUser(user UserCreate, id int) (UserCreate, error) {
	var updatedUser UserCreate

	query := fmt.Sprintf(
		`UPDATE %s
		SET name=$1, login=$2, password=$3, about=$4, address=$5, phone=$6, address=$5, phone=$6, surname=$7, avatar=$8
		WHERE id=$9 RETURNING *`,
		usersTable,
	)
	err := r.db.Get(&updatedUser, query,
		user.Name, user.Login, user.Password, user.About, user.Address, user.Phone, user.Surname, user.Avatar, id,
	)

	return updatedUser, err
}

func (r *Postgres) DeleteUser(id int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1",
		usersTable,
	)
	_, err := r.db.Exec(query, id)

	return err
}
