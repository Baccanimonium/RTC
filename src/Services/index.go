package Services

import (
	"go.mongodb.org/mongo-driver/bson"
	"video-chat-app/src/Models"
	"video-chat-app/src/Repos"
)

type Authorization interface {
	CreateUser(user Repos.UserCreate) (int, error)
	GenerateToken(user Repos.UserLogin) (string, error)
}

type DoctorService interface {
	CreateDoctor(doctor Repos.Doctor) (int, error)
	UpdateDoctor(doctor Repos.Doctor, id int) (Repos.Doctor, error)
	GetAllDoctor() ([]Repos.Participant, error)
	GetDoctorById(id int) (Repos.Participant, error)
	DeleteDoctor(id int) error
}
type UserService interface {
	UpdateUser(user Repos.UserCreate, id int) (Repos.UserCreate, error)
	GetAllUser() ([]Repos.UserCreate, error)
	GetUserById(id int) (Repos.UserCreate, error)
	DeleteUser(id int) error
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

type EventService interface {
	CreateEvent(idSchedule int, event Repos.Event) (int, error)
	UpdateEvent(event Repos.Event, id int) (Repos.Event, error)
	GetEventById(id int) (Repos.Event, error)
	GetAllEvents(idSchedule int) ([]Repos.Event, error)
	DeleteEvent(id int) error
}

type MessagesService interface {
	CreateMessage(newMessage Models.Message) (bson.D, error)
	GetMessage(messageId interface{}) (bson.D, error)
}

type Services struct {
	Authorization
	DoctorService
	PatientService
	ScheduleService
	EventService
	UserService
	MessagesService
}

func NewService(repo *Repos.Repo) *Services {
	return &Services{
		Authorization:   NewAuthService(repo.Authorization),
		DoctorService:   NewDoctorService(repo.DoctorRepo),
		PatientService:  NewPatientService(repo.PatientRepo),
		ScheduleService: NewScheduleService(repo.ScheduleRepo),
		EventService:    NewEventService(repo.EventRepo),
		UserService:     NewUserService(repo.UserRepo),
		MessagesService: NewMessagesService(repo.MessagesRepo),
	}
}
