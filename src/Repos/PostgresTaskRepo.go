package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Task struct {
	Id        int `json:"id"`
	IdPatient int `json:"id_patient" binding:"required"`
	IdUser    int `json:"id_user" binding:"required"`
	IdEvent   int `json:"id_event" binding:"required"`
	Weight    int `json:"weight" binding:"required"`
}

func NewTaskPostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateTask(task Task) (Task, error) {
	var result Task

	query := fmt.Sprintf(
		`INSERT INTO %s (id_patient, id_user, id_event, weight)
				values ($1, $2, $3, $4) RETURNING *`,
		tasksTable,
	)
	row := r.db.QueryRow(
		query,
		task.IdPatient,
		task.IdUser,
		task.IdEvent,
		task.Weight,
	)

	if err := row.Scan(&result); err != nil {
		return Task{}, err
	}

	return result, nil
}

func (r *Postgres) GetAllTasks(idDoctor int, idPatient int) ([]Task, error) {
	var result = make([]Task, 0)

	query := fmt.Sprintf(`
		SELECT
		id, id_patient, id_user, id_event, weight FROM %s
		WHERE id_user = $1 AND id_patient = $2`,
		tasksTable,
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
