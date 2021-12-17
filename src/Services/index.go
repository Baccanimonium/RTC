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
	CreateConsultation(consultation Repos.Consultation) (Repos.Consultation, error)
	GetAllConsultation(idDoctor int, idPatient int) ([]Repos.Consultation, error)
	GetConsultationById(idConsultation int) (Repos.Consultation, error)
	UpdateConsultation(consultation Repos.Consultation, id int) (Repos.Consultation, error)
	SetDoctorJoinTime(id int) error
	DeleteConsultation(idConsultation int) error
	CreateConsultationNotes(notes Repos.Notes) (Repos.Notes, error)
	UpdateConsultationNotes(notes Repos.Notes) (Repos.Notes, error)
	DeleteConsultationNotes(idNotes int) error
}

type EventService interface {
	CreateEvent(event Repos.Event) (Repos.Event, error)
	UpdateEvent(event Repos.Event) (Repos.Event, error)
	GetEventById(id int) (Repos.Event, error)
	GetAllEvents(request Repos.GetAllEventsParams) ([]Repos.Event, error)
	DeleteEvent(id int) (Repos.Event, error)
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

type TaskCandidatesService interface {
	CreateTaskCandidates(unConfirmedEvents []Repos.Event, currentDate, currentTime string) error
	DeleteTaskCandidate(taskId int) error
	ExtractTaskCandidates(taskTime string) ([]Repos.TaskCandidate, error)
	GetTaskCandidatesByPatient(patientId int) ([]Repos.TaskCandidate, error)
}

type TaskService interface {
	CreateTask(task Repos.TaskCandidate) (Repos.Task, error)
	GetAllTasks(idDoctor int, idPatient int) ([]Repos.Task, error)
	DeleteTask(idTask int) (Repos.Task, error)
}

type PatientCandidatesService interface {
	CreatePatientCandidate(patientCandidate Models.PatientCandidate) (interface{}, error)
	GetAllPatientCandidates() ([]Models.PatientCandidate, error)
}

type Services struct {
	Authorization
	DoctorService
	PatientService
	ConsultationService
	EventService
	UserService
	MessagesService
	ChannelsService
	TaskCandidatesService
	TaskService
	PatientCandidatesService
}

func NewService(repo *Repos.Repo, broadcast chan RTC.BroadcastingMessage) *Services {
	return &Services{
		Authorization:            NewAuthService(repo.Authorization),
		DoctorService:            NewDoctorService(repo.DoctorRepo),
		PatientService:           NewPatientService(repo.PatientRepo),
		ConsultationService:      NewConsultationService(repo.ConsultationRepo),
		EventService:             NewEventService(repo.EventRepo),
		UserService:              NewUserService(repo.UserRepo),
		MessagesService:          NewMessagesService(repo.MessagesRepo, broadcast),
		ChannelsService:          NewChannelsService(repo.ChannelsRepo, broadcast),
		TaskCandidatesService:    NewTaskCandidatesService(repo.TaskCandidatesRepo),
		TaskService:              NewTaskService(repo.TaskRepo),
		PatientCandidatesService: NewPatientCandidatesService(repo.PatientCandidatesRepo),
	}
}
