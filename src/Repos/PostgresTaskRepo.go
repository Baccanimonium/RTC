package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Task struct {
	Id        int    `json:"id"`
	IdPatient int    `json:"id_patient" binding:"required"`
	IdDoctor  int    `json:"id_doctor" binding:"required"`
	IdEvent   int    `json:"id_event" binding:"required"`
	Title     string `json:"title"`
	Date      string `json:"at_date"`
	Weight    int    `json:"weight"`
}

func NewTaskPostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateTask(task TaskCandidate) (Task, error) {
	var result Task

	query := fmt.Sprintf(
		`INSERT INTO %s (id_patient, id_doctor, id_event, weight, title, at_date)
				values ($1, $2, $3, $4, $5, $6) RETURNING *`,
		tasksTable,
	)
	row := r.db.QueryRow(
		query,
		task.IdPatient,
		task.IdDoctor,
		task.IdEvent,
		task.Weight,
		task.Title,
		task.Date,
	)

	if err := row.Scan(&result); err != nil {
		return Task{}, err
	}

	return result, nil
}

func (r *Postgres) GetAllTasks(idDoctor int, idPatient int) ([]Task, error) {
	var result = make([]Task, 0)
	var conditionQuery []string
	argsArray := make([]interface{}, 0)
	argsCount := 0

	if idPatient != 0 {
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			"id_patient = $%d",
			argsCount+1,
		))
		argsCount += 1
		argsArray = append(argsArray, idPatient)
	}

	conditionQuery = append(conditionQuery, fmt.Sprintf(
		"id_doctor = $%d",
		argsCount+1,
	))
	argsCount += 1
	argsArray = append(argsArray, idDoctor)

	query := fmt.Sprintf(`
		SELECT id, id_patient, id_doctor, id_event, weight, title, array_agg(at_date) FROM %s
		WHERE %s
		GROUP BY id_event
		`,
		tasksTable,
		strings.Join(conditionQuery, " AND "),
	)

	err := r.db.Select(&result, query, idDoctor, idPatient)

	return result, err
}

func (r *Postgres) DeleteTask(idTask int) (Task, error) {
	var result Task
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1 RETURNING *",
		tasksTable,
	)

	err := r.db.Get(&result, query, idTask)

	return result, err
}
