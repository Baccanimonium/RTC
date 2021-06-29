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

type Repo struct {
	Authorization
	DoctorRepo
	PatientRepo
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		Authorization: NewAuthPostgresRepo(db),
		DoctorRepo:    NewDoctorPostgresRepo(db),
		PatientRepo:   NewPatientPostgresRepo(db),
	}
}
