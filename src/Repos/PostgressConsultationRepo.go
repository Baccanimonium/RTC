package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Consultation struct {
	Id        int    `json:"id" db:"id"`
	IdPatient int    `json:"id_patient"`
	IdUsers   int    `json:"id_users"`
	IdCourse  int    `json:"id_course"`
	Title     string `json:"title"`
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
}

func NewConsultationPostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateConsultation(consultation Consultation) (int, error) {
	var id int

	query := fmt.Sprintf(
		"INSERT INTO %s (id_user, id_patient, id_course, title, time_start, time_end) values ($1, $2, $3, $4, $5, $6) RETURNING id",
		consultationTable,
	)
	row := r.db.QueryRow(query, consultation.IdPatient, consultation.IdUsers, consultation.IdCourse, consultation.Title, consultation.TimeStart, consultation.TimeEnd)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Postgres) GetAllConsultation() ([]Consultation, error) {
	var consultation []Consultation

	query := fmt.Sprintf(`SELECT id, id_user, id_patient, id_course, title, time_start, time_end FROM %s`,
		consultationTable,
	)

	err := r.db.Select(&consultation, query)

	return consultation, err
}

func (r *Postgres) GetConsultationById(id int) (Consultation, error) {
	var consultation Consultation

	query := fmt.Sprintf(`SELECT id, id_patient, id_users, id_course, title, time_start, time_end FROM %s WHERE id = $1`,
		consultationTable,
	)

	err := r.db.Get(&consultation, query, id)

	return consultation, err
}

func (r *Postgres) UpdateConsultation(consultation Consultation, id int) (Consultation, error) {
	var newConsultation Consultation
	query := fmt.Sprintf(
		"UPDATE %s SET title=$1, time_start=$2, time_end=$3 WHERE id=$4 RETURNING *",
		consultationTable,
	)
	err := r.db.Get(&newConsultation, query, consultation.Title, consultation.TimeStart, consultation.TimeEnd, id)
	return newConsultation, err
}

func (r *Postgres) getAllConsultation() ([]Consultation, error) {
	var consultations []Consultation

	query := fmt.Sprintf(`SELECT id, id_patient, id_users, id_course, title, time_start, time_end FROM %s WHERE id = $1`,
		consultationTable,
	)

	err := r.db.Select(&consultations, query)

	return consultations, err
}

func (r *Postgres) DeleteConsultation(id int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1",
		consultationTable,
	)
	_, err := r.db.Exec(query, id)

	return err

}
