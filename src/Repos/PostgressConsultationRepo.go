package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
	"video-chat-app/src/Models"
)

/*
	Консультация создается каждый раз вручную
	О необходимости создания консультации необходимо уведомлять и доктора и пациента
	Пушить доктору и пациенту сообщение которое будет показывать кнопку входа в консультациюю как в скайинге
	Иметь отдельную джобу которая будет рассылать уведомления о предстоящей консультации или о начале консультации

	На фронте при создани консультации должен быть особый UI который позволять создавать заблаговременное уведомление по
	упрощенной форме
*/

func NewConsultationPostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateConsultation(consultation Models.Consultation) (Models.Consultation, error) {
	var result Models.Consultation

	query := fmt.Sprintf(
		`INSERT INTO %s (id_patient, id_doctor, start, end, offline)
				values ($1, $2, $3, $4, $5, $6) RETURNING *`,
		consultationTable,
	)
	row := r.db.QueryRow(
		query,
		consultation.IdPatient,
		consultation.IdDoctor,
		consultation.Start,
		consultation.End,
		consultation.Offline,
	)

	if err := row.Scan(&result); err != nil {
		return Models.Consultation{}, err
	}

	return result, nil
}

func (r *Postgres) GetAllConsultation(params Models.GetConsultationList) ([]Models.Consultation, error) {
	var result = make([]Models.Consultation, 0)
	var conditionQuery []string
	argsArray := make([]interface{}, 0)
	argsCount := 0

	if params.IdPatient != 0 {
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			"id_patient = $%d",
			argsCount+1,
		))
		argsCount += 1
		argsArray = append(argsArray, params.IdPatient)
	}

	if params.IdDoctor != 0 {
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			"id_doctor = $%d",
			argsCount+1,
		))
		argsCount += 1
		argsArray = append(argsArray, params.IdDoctor)

	}

	if params.Start != 0 {
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			"start >= $%d",
			argsCount+1,
		))
		argsCount += 1
		argsArray = append(argsArray, params.Start)

	}

	if params.End != 0 {
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			"end <= $%d",
			argsCount+1,
		))
		argsCount += 1
		argsArray = append(argsArray, params.End)

	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s",
		consultationTable,
		strings.Join(conditionQuery, " AND "),
	)

	err := r.db.Select(&result, query)

	return result, err
}

func (r *Postgres) GetConsultationById(idConsultation int) (Models.Consultation, error) {
	var consultation Models.Consultation

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`,
		consultationTable,
	)

	err := r.db.Get(&consultation, query, idConsultation)

	return consultation, err
}

func (r *Postgres) UpdateConsultation(consultation Models.Consultation, id int) (Models.Consultation, error) {
	var newConsultation Models.Consultation
	query := fmt.Sprintf(
		"UPDATE %s SET start=$1, end=$2, offline=$3 WHERE id=$4 RETURNING *",
		consultationTable,
	)
	err := r.db.Get(&newConsultation, query, consultation.Start, consultation.End, consultation.Offline, id)
	return newConsultation, err
}

func (r *Postgres) SetDoctorJoinTime(id int) error {
	query := fmt.Sprintf(
		"UPDATE %s SET doctor_joined=$1 WHERE id=$2 RETURNING *",
		consultationTable,
	)
	_, err := r.db.Exec(query, time.Now().Format(time.RFC3339), id)

	return err
}

func (r *Postgres) DeleteConsultation(idConsultation int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1",
		consultationTable,
	)
	_, err := r.db.Exec(query, idConsultation)

	return err
}

func (r *Postgres) CreateConsultationNotes(notes Models.Notes) (Models.Notes, error) {
	var result Models.Notes

	query := fmt.Sprintf(
		`INSERT INTO %s (notes, id_consultation, for_doctor, files)
				values ($1, $2, $3, $4) RETURNING *`,
		consultationFilesTable,
	)
	row := r.db.QueryRow(
		query,
		notes.Notes,
		notes.IdConsultation,
		notes.ForDoctor,
		notes.Files,
	)

	if err := row.Scan(&result); err != nil {
		return Models.Notes{}, err
	}

	return result, nil
}

func (r *Postgres) UpdateConsultationNotes(notes Models.Notes) (Models.Notes, error) {
	var result Models.Notes

	query := fmt.Sprintf(
		`UPDATE %s SET notes=$1, files=$2 WHERE id=$3 RETURNING *`,
		consultationFilesTable,
	)
	row := r.db.QueryRow(
		query,
		notes.Notes,
		notes.Files,
		notes.Id,
	)

	if err := row.Scan(&result); err != nil {
		return Models.Notes{}, err
	}

	return result, nil
}

func (r *Postgres) DeleteConsultationNotes(idNotes int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1",
		consultationFilesTable,
	)
	_, err := r.db.Exec(query, idNotes)

	return err
}
