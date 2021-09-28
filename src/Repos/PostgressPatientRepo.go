package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Patient struct {
	Id              int    `json:"id" db:"id_patient"`
	IdUser          int    `json:"id_user" db:"id_user" binding:"required"`
	IdCurrentDoctor int    `json:"id_current_doctor" db:"id_current_doctor"`
	Description     string `json:"description"`
	Recovered       bool   `json:"recovered"`
}

func NewPatientPostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreatePatient(patient Patient) (int, error) {
	var id int

	query := fmt.Sprintf(
		`INSERT INTO %s (id_user, description, recovered, id_current_doctor)
				values ($1, $2, $3, $4)
				RETURNING id`,
		patientTable,
	)
	row := r.db.QueryRow(query, patient.IdUser, patient.Description, patient.Recovered, patient.IdCurrentDoctor)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Postgres) GetAllPatient(userId int) ([]Participant, error) {
	var patients = make([]Participant, 0)

	query := fmt.Sprintf(`SELECT
		doc.id as id_patient, doc.id_user, doc.description, doc.recovered, doc.id_current_doctor,
		us.address, us.about, us.name, us.phone, us.id, us.avatar 
		FROM %s doc INNER JOIN %s us ON doc.id_user = us.id WHERE doc.id_current_doctor = $1`,
		patientTable, usersTable,
	)

	err := r.db.Select(&patients, query, userId)

	return patients, err
}

func (r *Postgres) GetPatientById(id int) (Participant, error) {
	var patients Participant

	query := fmt.Sprintf(`SELECT
		doc.id as id_patient, doc.id_user, doc.description, doc.recovered, us.address, us.about, us.name, us.phone, us.id 
		FROM %s doc INNER JOIN %s us ON doc.id_user = us.id WHERE doc.id = $1`,
		patientTable, usersTable,
	)

	err := r.db.Get(&patients, query, id)

	return patients, err
}

func (r *Postgres) UpdatePatient(patient Patient, id int) (Patient, error) {
	var newPatient Patient

	query := fmt.Sprintf(
		`UPDATE %s SET description=$1, recovered=$2, WHERE id=$3
				RETURNING doc.id as id_patient, doc.id_user, doc.description, doc.recovered`,
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
