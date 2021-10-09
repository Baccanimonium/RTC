package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Event struct {
	Id          int    `json:"id" db:"id"`
	IdPatient   string `json:"id_patient" db:"id_patient" binding:"required"`
	IdUser      int    `json:"id_users" db:"id_user" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeStart   string `json:"time_start"`
	TimeEnd     string `json:"time_end"`
}

func NewEventPostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateEvent(idSchedule int, event Event) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int

	query := fmt.Sprintf(
		`INSERT INTO %s
                (id_patient, id_users, title, description, time_start, time_end) values ($1, $2, $3, $4, $5, $6)
				RETURNING id`,
		eventTable,
	)
	row := tx.QueryRow(query, event.IdPatient, event.IdUser, event.Title, event.Description, event.TimeStart, event.TimeEnd)

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (id_schedule, id_event) VALUES ($1, $2)", scheduleEventTable)

	_, err = tx.Exec(createUsersListQuery, idSchedule, id)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, nil
}

func (r *Postgres) UpdateEvent(event Event, id int) (Event, error) {
	var newSchedule Event

	query := fmt.Sprintf(
		`UPDATE
				%s ev SET title=$1, description=$2, time_start=$3, time_end=$4 
				FROM %s se WHERE ev.id = se.id_event AND se.id_event = $5 RETURNING *`,
		eventTable,
		scheduleEventTable,
	)
	err := r.db.Get(&newSchedule, query, event.Title, event.Description, event.TimeStart, event.TimeEnd, id)

	return newSchedule, err
}

func (r *Postgres) GetEventById(id int) (Event, error) {
	var event Event

	query := fmt.Sprintf(
		`SELECT
				ev.id, ev.id_patient, ev.id_users, ev.title, ev.description, ev.time_start, ev.time_end FROM %s ev
				INNER JOIN %s se on ev.id = se.id_event WHERE se.id_event = $1`,
		eventTable,
		scheduleEventTable,
	)

	err := r.db.Get(&event, query, id)

	return event, err
}

func (r *Postgres) GetAllEvents(idSchedule int) ([]Event, error) {
	var schedules = make([]Event, 0)

	query := fmt.Sprintf(
		`SELECT
				ev.id, ev.id_patient, ev.id_users, ev.title, ev.description, ev.time_start, ev.time_end FROM %s ev
				INNER JOIN %s se on ev.id = se.id_event WHERE se.id_schedule = $1`,
		eventTable,
		scheduleEventTable,
	)

	err := r.db.Select(&schedules, query, idSchedule)

	return schedules, err
}

func (r *Postgres) DeleteEvent(id int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s ev USING %s se WHERE ev.id = se.id_event AND se.id_event = $1",
		eventTable,
		scheduleEventTable,
	)
	_, err := r.db.Exec(query, id)

	return err
}
