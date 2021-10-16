package Services

import (
	"go.mongodb.org/mongo-driver/bson"
	RTC "video-chat-app"
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
	GetAllPatient(userId int) ([]Repos.Participant, error)
	DeletePatient(id int) error
}

type ConsultationService interface {
	CreateConsultation(idSchedule int, consultation Repos.Consultation) (Repos.Consultation, error)
	UpdateConsultation(consultation Repos.Consultation, id int) (Repos.Consultation, error)
	GetConsultationById(idSchedule, idConsultation int) (Repos.Consultation, error)
	GetAllConsultation(idSchedule int) ([]Repos.Consultation, error)
	DeleteConsultation(idSchedule, idConsultation int) error
}

type ScheduleService interface {
	CreateSchedule(schedule Repos.Schedule) (Repos.Schedule, error)
	UpdateSchedule(schedule Repos.Schedule, id int) (Repos.Schedule, error)
	GetScheduleById(id int) (Repos.Schedule, error)
	GetAllSchedule(idPatient int) ([]Repos.Schedule, error)
	DeleteSchedule(id int) error
}

type EventService interface {
	CreateEvent(idSchedule int, event Repos.Event) (int, error)
	UpdateEvent(event Repos.Event, id int) (Repos.Event, error)
	GetEventById(id int) (Repos.Event, error)
	GetAllEvents(idSchedule int) ([]Repos.Event, error)
	GetEventsByDate(date string) ([]Repos.Event, error)
	DeleteEvent(id int) error
}

type MessagesService interface {
	CreateMessage(newMessage Models.CreateMessage) (bson.M, error)
	GetMessage(messageId interface{}) (bson.M, error)
	GetMessages(channelId string, userId interface{}) ([]Models.Message, error)
}

type ChannelsService interface {
	CreateChannel(userId int, payload Models.Channel) (bson.M, error)
	DeleteChannel(userId int, payload Models.Channel) (bson.M, error)
	GetChannelByID(documentId interface{}) (bson.M, error)
	GetChannelByParticipants(userId int, payload map[string]interface{}) (Models.Channel, error)
	GetAllChannelsBelongsToUser(userId int) ([]Models.Channel, error)
}

type Services struct {
	Authorization
	DoctorService
	PatientService
	ScheduleService
	ConsultationService
	EventService
	UserService
	MessagesService
	ChannelsService
}

func NewService(repo *Repos.Repo, broadcast chan RTC.BroadcastingMessage) *Services {
	return &Services{
		Authorization:       NewAuthService(repo.Authorization),
		DoctorService:       NewDoctorService(repo.DoctorRepo),
		PatientService:      NewPatientService(repo.PatientRepo),
		ScheduleService:     NewScheduleService(repo.ScheduleRepo, broadcast),
		ConsultationService: NewConsultationService(repo.ConsultationRepo),
		EventService:        NewEventService(repo.EventRepo),
		UserService:         NewUserService(repo.UserRepo),
		MessagesService:     NewMessagesService(repo.MessagesRepo, broadcast),
		ChannelsService:     NewChannelsService(repo.ChannelsRepo, broadcast),
	}
}
