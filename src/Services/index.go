package Services

import (
	"video-chat-app/src/Repos"
)

type Authorization interface {
	CreateUser(user Repos.UserCreate) (int, error)
	GenerateToken(user Repos.UserLogin) (string, error)
	ParseToken(rawToken string) (int, error)
}

type DoctorService interface {
	CreateDoctor(doctor Repos.Doctor) (int, error)
	UpdateDoctor(doctor Repos.Doctor, id int) (Repos.Doctor, error)
	GetAllDoctor() ([]Repos.Participant, error)
	GetDoctorById(id int) (Repos.Participant, error)
	DeleteDoctor(id int) error
}
type PatientService interface {
	CreatePatient(patient Repos.Patient) (int, error)
	UpdatePatient(doctor Repos.Patient, id int) (Repos.Patient, error)
	GetPatientById(id int) (Repos.Participant, error)
	GetAllPatient() ([]Repos.Participant, error)
	DeletePatient(id int) error
}
type ScheduleService interface {
	CreateSchedule(schedule Repos.Schedule) (int, error)
	UpdateSchedule(schedule Repos.Schedule, id int) (Repos.Schedule, error)
	GetScheduleById(id int) (Repos.Schedule, error)
	GetAllSchedule() ([]Repos.Schedule, error)
	DeleteSchedule(id int) error
}
type ConsultationService interface {
	CreateConsultation(consultation Repos.Consultation) (int, error)
	UpdateConsultation(consultation Repos.Consultation, id int) (Repos.Consultation, error)
	GetConsultationById(id int) (Repos.Consultation, error)
	GetAllConsultation() ([]Repos.Consultation, error)
	DeleteConsultation(id int) error
}

type Services struct {
	Authorization
	DoctorService
	PatientService
	ScheduleService
	ConsultationService
}

func NewService(repo *Repos.Repo) *Services {
	return &Services{
		Authorization:       NewAuthService(repo.Authorization),
		DoctorService:       NewDoctorService(repo.DoctorRepo),
		PatientService:      NewPatientService(repo.PatientRepo),
		ScheduleService:     NewScheduleService(repo.ScheduleRepo),
		ConsultationService: NewConsultationService(repo.ConsultationRepo),
	}
}
