package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Patient struct {
	Id          int    `json:"id" db:"id"`
	IdUser      int    `json:"id_user" db:"id_user" binding:"required"`
	Description string `json:"description"`
	Recovered   bool   `json:"recovered"`
}

func NewPatientPostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreatePatient(patient Patient) (int, error) {
	var id int

	query := fmt.Sprintf(
		"INSERT INTO %s (id_user, description, recovered) values ($1, $2, $3) RETURNING id",
		patientTable,
	)
	row := r.db.QueryRow(query, patient.IdUser, patient.Description, patient.Recovered)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Postgres) GetAllPatient() ([]Participant, error) {
	var patients []Participant

	query := fmt.Sprintf(`SELECT
		doc.id, doc.id_user, doc.description, doc.recovered, us.address, us.about, us.name, us.phone
		FROM %s doc INNER JOIN %s us ON doc.id_user = us.id`,
		patientTable, usersTable,
	)

	err := r.db.Select(&patients, query)

	return patients, err
}

func (r *Postgres) GetPatientById(id int) (Participant, error) {
	var patients Participant

	query := fmt.Sprintf(`SELECT
		doc.id, doc.id_user, doc.description, doc.recovered, us.address, us.about, us.name, us.phone
		FROM %s doc INNER JOIN %s us ON doc.id_user = us.id WHERE doc.id = $1`,
		patientTable, usersTable,
	)

	err := r.db.Get(&patients, query, id)

	return patients, err
}

func (r *Postgres) UpdatePatient(patient Patient, id int) (Patient, error) {
	var newPatient Patient

	query := fmt.Sprintf(
		"UPDATE %s SET description=$1, recovered=$2, WHERE id=$3 RETURNING *",
		patientTable,
	)
	err := r.db.Get(&newPatient, query, patient.Description, patient.Recovered, id)

	return newPatient, err
}

func (r *Postgres) DeletePatient(id int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1",
		patientTable,
	)
	_, err := r.db.Exec(query, id)

	return err
}
