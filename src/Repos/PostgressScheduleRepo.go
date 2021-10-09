package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Schedule struct {
	Id          int    `json:"id" db:"id"`
	IdUser      int    `json:"id_user" db:"id_user" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewSchedulePostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateSchedule(schedule Schedule) (int, error) {
	var id int

	query := fmt.Sprintf(
		"INSERT INTO %s (id_user, title, description) values ($1, $2, $3) RETURNING id",
		scheduleTable,
	)
	row := r.db.QueryRow(query, schedule.IdUser, schedule.Title, schedule.Description)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Postgres) GetScheduleById(id int) (Schedule, error) {
	var schedule Schedule

	query := fmt.Sprintf(`SELECT id, id_user, title, description FROM %s WHERE id = $1`,
		scheduleTable,
	)

	err := r.db.Get(&schedule, query, id)

	return schedule, err
}

func (r *Postgres) GetAllSchedule() ([]Schedule, error) {
	var schedules []Schedule

	query := fmt.Sprintf(`SELECT id, id_user, title description FROM %s`,
		scheduleTable,
	)

	err := r.db.Select(&schedules, query)

	return schedules, err
}

func (r *Postgres) UpdateSchedule(schedule Schedule, id int) (Schedule, error) {
	var newSchedule Schedule
	query := fmt.Sprintf(
		"UPDATE %s SET title=$1, description=$2 WHERE id=$3 RETURNING *",
		scheduleTable,
	)
	err := r.db.Get(&newSchedule, query, schedule.Title, schedule.Description, id)
	return newSchedule, err
}

func (r *Postgres) getAllSchedule() ([]Schedule, error) {
	var schedules []Schedule

	query := fmt.Sprintf(`SELECT id, id_user, title, description FROM %s WHERE id = $1`,
		scheduleTable,
	)

	err := r.db.Select(&schedules, query)

	return schedules, err
}

func (r *Postgres) DeleteSchedule(id int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1",
		scheduleTable,
	)
	_, err := r.db.Exec(query, id)

	return err

}
