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

func (r *Postgres) CreateConsultation(idSchedule int, consultation Consultation) (Consultation, error) {
	var result Consultation

	query := fmt.Sprintf(
		"INSERT INTO %s (id_user, id_patient, id_course, title, time_start, time_end) values ($1, $2, $3, $4, $5, $6) RETURNING *",
		consultationTable,
	)
	row := r.db.QueryRow(query, consultation.IdPatient, consultation.IdUsers, consultation.IdCourse, consultation.Title, consultation.TimeStart, consultation.TimeEnd)

	if err := row.Scan(&result); err != nil {
		return Consultation{}, err
	}

	return result, nil
}

func (r *Postgres) GetAllConsultation(idSchedule int) ([]Consultation, error) {
	var consultation []Consultation

	query := fmt.Sprintf(`SELECT id, id_user, id_patient, id_course, title, time_start, time_end FROM %s`,
		consultationTable,
	)

	err := r.db.Select(&consultation, query)

	return consultation, err
}

func (r *Postgres) GetConsultationById(idSchedule, idConsultation int) (Consultation, error) {
	var consultation Consultation

	query := fmt.Sprintf(`SELECT id, id_patient, id_users, id_course, title, time_start, time_end FROM %s WHERE id = $1`,
		consultationTable,
	)

	err := r.db.Get(&consultation, query, idConsultation)

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

func (r *Postgres) DeleteConsultation(idSchedule, idConsultation int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1",
		consultationTable,
	)
	_, err := r.db.Exec(query, idConsultation)

	return err

}
