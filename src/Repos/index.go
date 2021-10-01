package Repos

import (
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"video-chat-app/src/Models"
)

type Postgres struct {
	db *sqlx.DB
}

type Mongo struct {
	db *mongo.Database
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

type UserRepo interface {
	UpdateUser(user UserCreate, id int) (UserCreate, error)
	GetAllUser() ([]UserCreate, error)
	GetUserById(id int) (UserCreate, error)
	DeleteUser(id int) error
}

type PatientRepo interface {
	CreatePatient(patient Patient) (int, error)
	UpdatePatient(patient Patient, id int) (Patient, error)
	GetPatientById(id int) (Participant, error)
	GetAllPatient(userId int) ([]Participant, error)
	DeletePatient(id int) error
}

type ScheduleRepo interface {
	CreateSchedule(schedule Schedule) (Schedule, error)
	UpdateSchedule(schedule Schedule, id int) (Schedule, error)
	GetScheduleById(id int) (Schedule, error)
	GetAllSchedule() ([]Schedule, error)
	DeleteSchedule(id int) error
}

type ConsultationRepo interface {
	CreateConsultation(consultation Consultation) (int, error)
	GetConsultationById(id int) (Consultation, error)
	GetAllConsultation() ([]Consultation, error)
	UpdateConsultation(consultation Consultation, id int) (Consultation, error)
	DeleteConsultation(id int) error
}

type EventRepo interface {
	CreateEvent(idSchedule int, event Event) (int, error)
	UpdateEvent(event Event, id int) (Event, error)
	GetEventById(id int) (Event, error)
	GetAllEvents(idSchedule int) ([]Event, error)
	DeleteEvent(id int) error
}

type MessagesRepo interface {
	CreateMessage(newMessage Models.CreateMessage) (bson.M, error)
	GetMessage(messageId interface{}) (bson.M, error)
	GetMessages(channelId string) ([]Models.Message, error)
	UpdateMessage(updatedMessage Models.Message) (bson.M, error)
	DeleteMessage(message Models.DeleteMessage) (bson.M, error)
}

type ChannelsRepo interface {
	CreateChannel(userId int, payload Models.Channel) (bson.M, error)
	DeleteChannel(userId int, channel Models.Channel) (bson.M, error)
	GetChannelByID(documentId interface{}) (bson.M, error)
	GetChannelByParticipants(userId int, payload map[string]interface{}) (Models.Channel, error)
	GetAllChannelsBelongsToUser(creatorId int) ([]Models.Channel, error)
}

type Repo struct {
	Authorization
	DoctorRepo
	PatientRepo
	ScheduleRepo
	ConsultationRepo
	EventRepo
	UserRepo
	MessagesRepo
	ChannelsRepo
}

func NewRepo(db *sqlx.DB, mongoDB *mongo.Database) *Repo {
	return &Repo{
		Authorization: NewAuthPostgresRepo(db),
		DoctorRepo:    NewDoctorPostgresRepo(db),
		PatientRepo:   NewPatientPostgresRepo(db),
		ScheduleRepo:  NewSchedulePostgresRepo(db),
		ConsultationRepo: NewConsultationPostgresRepo(db),
		EventRepo:     NewEventPostgresRepo(db),
		UserRepo:      NewUserPostgresRepo(db),
		MessagesRepo:  NewMongoMessagesRepo(mongoDB),
		ChannelsRepo:  NewMongoChannelsRepo(mongoDB),
	}
}
