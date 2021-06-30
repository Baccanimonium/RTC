package Repos

import (
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

type Authorization interface {
	CreateUser(user UserCreate) (int, error)
	GetUser(login, password string) (UserCreate, error)
}

type DoctorRepo interface {
	CreateDoctor(doctor Doctor) (int, error)
	UpdateDoctor(doctor Doctor, id int) (Doctor, error)
	GetAllDoctor() ([]Participant, error)
	GetDoctorById(id int) (Participant, error)
	DeleteDoctor(id int) error
}

type PatientRepo interface {
	CreatePatient(patient Patient) (int, error)
	UpdatePatient(patient Patient, id int) (Patient, error)
	GetPatientById(id int) (Participant, error)
	GetAllPatient() ([]Participant, error)
	DeletePatient(id int) error
}

type ScheduleRepo interface {
	CreateSchedule(schedule Schedule) (int, error)
	UpdateSchedule(schedule Schedule, id int) (Schedule, error)
	GetScheduleById(id int) (Schedule, error)
	GetAllSchedule() ([]Schedule, error)
	DeleteSchedule(id int) error
}

type EventRepo interface {
	CreateEvent(idSchedule int, event Event) (int, error)
	UpdateEvent(event Event, id int) (Event, error)
	GetEventById(id int) (Event, error)
	GetAllEvents(idSchedule int) ([]Event, error)
	DeleteEvent(id int) error
}

type Repo struct {
	Authorization
	DoctorRepo
	PatientRepo
	ScheduleRepo
	EventRepo
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		Authorization: NewAuthPostgresRepo(db),
		DoctorRepo:    NewDoctorPostgresRepo(db),
		PatientRepo:   NewPatientPostgresRepo(db),
		ScheduleRepo:  NewSchedulePostgresRepo(db),
		EventRepo:     NewEventPostgresRepo(db),
	}
}
