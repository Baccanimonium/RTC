package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Doctor struct {
	Id             int    `json:"id" db:"id"`
	IdUser         string `json:"id_user" binding:"required"`
	Salary         string `json:"salary"`
	Qualifications string `json:"qualifications"`
	Contacts       string `json:"contacts"`
}

func NewDoctorPostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateDoctor(doctor Doctor) (int, error) {
	var id int

	query := fmt.Sprintf(
		"INSERT INTO %s (id_user, salary, qualifications, contacts) values ($1, $2, $3, $4) RETURNING id",
		doctorTable,
	)
	row := r.db.QueryRow(query, doctor.IdUser, doctor.Salary, doctor.Qualifications, doctor.Contacts)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
