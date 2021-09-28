package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Doctor struct {
	Id             int     `json:"id,omitempty" db:"id_doctor"`
	IdUser         int     `json:"id_user,omitempty" db:"id_origin" binding:"required"`
	Salary         float64 `json:"salary,omitempty"`
	Qualifications string  `json:"qualifications,omitempty"`
	Contacts       string  `json:"contacts,omitempty"`
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
		doc.id, doc.id_user as id_origin, doc.salary, doc.qualifications, doc.contacts,
		us.address, us.about, us.name, us.phone
		FROM %s doc INNER JOIN %s us ON doc.id_user = us.id`,
		doctorTable, usersTable,
	)

	err := r.db.Select(&doctors, query)

	return doctors, err
}

func (r *Postgres) GetDoctorById(id int) (Participant, error) {
	var doctors Participant

	query := fmt.Sprintf(`SELECT
		doc.id, doc.id_user as id_origin, doc.salary, doc.qualifications, doc.contacts,
		us.address, us.about, us.name, us.phone
		FROM %s doc INNER JOIN %s us ON doc.id_user = us.id WHERE doc.id = $1`,
		doctorTable, usersTable,
	)

	err := r.db.Get(&doctors, query, id)

	return doctors, err
}

func (r *Postgres) UpdateDoctor(doctor Doctor, id int) (Doctor, error) {
	var newDoctor Doctor

	query := fmt.Sprintf(
		"UPDATE %s SET salary=$1, qualifications=$2, contacts=$3 WHERE id=$4 RETURNING *",
		doctorTable,
	)
	err := r.db.Get(&newDoctor, query, doctor.Salary, doctor.Qualifications, doctor.Contacts, id)

	return newDoctor, err
}

func (r *Postgres) DeleteDoctor(id int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1",
		doctorTable,
	)
	_, err := r.db.Exec(query, id)

	return err
}
