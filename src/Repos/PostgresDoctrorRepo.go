package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Doctor struct {
	Id             int     `json:"id" db:"id"`
	IdUser         int     `json:"id_user" binding:"required"`
	Salary         float64 `json:"salary"`
	Qualifications string  `json:"qualifications"`
	Contacts       string  `json:"contacts"`
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

func (r *Postgres) GetAllDoctor() ([]Participant, error) {
	var doctors []Participant

	query := fmt.Sprintf(`SELECT
		doc.id, doc.id_user, doc.salary, doc.qualifications, doc.contacts, us.address, us.about, us.name, us.phone
		FROM %s doc INNER JOIN %s us ON doc.id_user = us.id`,
		doctorTable, usersTable,
	)

	err := r.db.Select(&doctors, query)

	return doctors, err
}
