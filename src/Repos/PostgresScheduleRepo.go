package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Schedule struct {
	Id          int    `json:"id" db:"id"`
	IdPatient   int    `json:"id_patient" binding:"required"`
	IdUser      int    `json:"id_user" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	TimeStart   string `json:"time_start" binding:"required"`
	TimeEnd     string `json:"time_end" binding:"required"`
}

func NewSchedulePostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateSchedule(schedule Schedule) (Schedule, error) {
	var result Schedule

	query := fmt.Sprintf(
		`INSERT INTO %s (id_user, id_patient, title, description, time_start, time_end)
				values ($1, $2, $3, $4, $5, $6) RETURNING *`,
		scheduleTable,
	)
	row := r.db.QueryRow(
		query,
		schedule.IdUser,
		schedule.IdPatient,
		schedule.Title,
		schedule.Description,
		schedule.TimeStart,
		schedule.TimeEnd,
	)

	if err := row.Scan(&result); err != nil {
		return Schedule{}, err
	}

	return result, nil
}

func (r *Postgres) UpdateSchedule(schedule Schedule, id int) (Schedule, error) {
	var newSchedule Schedule

	query := fmt.Sprintf(
		"UPDATE %s SET title=$1, description=$2, time_end=$3 WHERE id=$4 RETURNING *",
		scheduleTable,
	)
	err := r.db.Get(&newSchedule, query, schedule.Title, schedule.Description, schedule.TimeEnd, id)

	return newSchedule, err
}

func (r *Postgres) GetScheduleById(id int) (Schedule, error) {
	var schedule Schedule

	query := fmt.Sprintf(
		`SELECT id, id_patient, id_user, title, description, time_start, time_end
				FROM %s WHERE id = $1`,
		scheduleTable,
	)

	err := r.db.Get(&schedule, query, id)

	return schedule, err
}

func (r *Postgres) GetAllSchedule(idPatient int) ([]Schedule, error) {
	var schedules = make([]Schedule, 0)

	query := fmt.Sprintf(
		`SELECT id, id_patient, id_user, title, description, time_start, time_end
				FROM %s WHERE id_patient = $1`,
		scheduleTable,
	)

	err := r.db.Select(&schedules, query, idPatient)

	return schedules, err
}

func (r *Postgres) DeleteSchedule(id int) (Schedule, error) {
	var deletedSchedule Schedule
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1 RETURNING *",
		scheduleTable,
	)

	err := r.db.Get(&deletedSchedule, query, id)
	//_, err := r.db.Exec(query, id)

	return deletedSchedule, err
}
