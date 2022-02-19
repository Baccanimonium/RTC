package Repos

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"video-chat-app/src/Models"
)

type Postgres struct {
	db *sqlx.DB
}

type Mongo struct {
	db *mongo.Database
}

type Redis struct {
	rdb *redis.Client
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
type ScheduleRepo interface {
	CreateSchedule(schedule Models.DoctorSchedule) (Models.DoctorSchedule, error)
	UpdateSchedule(schedule Models.DoctorSchedule) (Models.DoctorSchedule, error)
	GetScheduleByDoctorId(id int) (Models.DoctorSchedule, error)
	GetAllSchedule(params Models.PostgresPagination) ([]Models.DoctorSchedule, error)
	DeleteSchedule(idDoctor int) (Models.DoctorSchedule, error)
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

type ConsultationRepo interface {
	CreateConsultation(consultation Models.Consultation) (Models.Consultation, error)
	GetAllConsultation(Models.GetConsultationList) ([]Models.Consultation, error)
	GetConsultationById(idConsultation int) (Models.Consultation, error)
	UpdateConsultation(consultation Models.Consultation, id int) (Models.Consultation, error)
	SetDoctorJoinTime(id int) error
	DeleteConsultation(idConsultation int) error
	CreateConsultationNotes(notes Models.Notes) (Models.Notes, error)
	UpdateConsultationNotes(notes Models.Notes) (Models.Notes, error)
	DeleteConsultationNotes(idNotes int) error
}

type EventRepo interface {
	CreateEvent(event Event) (Event, error)
	UpdateEvent(event Event) (Event, error)
	GetEventById(id int) (Event, error)
	GetAllEvents(request GetAllEventsParams) ([]Event, error)
	DeleteEvent(id int) (Event, error)
}

type MessagesRepo interface {
	CreateMessage(newMessage Models.CreateMessage) (bson.M, error)
	GetMessage(messageId interface{}) (bson.M, error)
	GetMessages(channelId string) ([]Models.Message, error)
	UpdateMessage(updatedMessage Models.Message, userId int) (bson.M, error)
	DeleteMessage(message Models.DeleteMessage, userId int) (bson.M, error)
}

type ChannelsRepo interface {
	CreateChannel(userId int, payload Models.Channel) (bson.M, error)
	DeleteChannel(userId int, channel Models.Channel) (bson.M, error)
	GetChannelByID(documentId interface{}) (bson.M, error)
	GetChannelByParticipants(userId int, payload map[string]interface{}) (Models.Channel, error)
	GetAllChannelsBelongsToUser(creatorId int) ([]Models.Channel, error)
}

type TaskCandidatesRepo interface {
	CreateTaskCandidates(unConfirmedEvents []Event, currentDate, currentTime string) error
	DeleteTaskCandidate(candidateId int) error
	ExtractTaskCandidates(taskTime string) ([]TaskCandidate, error)
	GetTaskCandidatesByPatient(patientId int) ([]TaskCandidate, error)
}

type TaskRepo interface {
	CreateTask(task TaskCandidate) (Task, error)
	GetAllTasks(idDoctor int, idPatient int) ([]Task, error)
	DeleteTask(idTask int) (Task, error)
}

type PatientCandidatesRepo interface {
	CreatePatientCandidate(patientCandidate Models.PatientCandidate) (interface{}, error)
	GetAllPatientCandidates() ([]Models.PatientCandidate, error)
}

type GroupsRepo interface {
	CreateGroup(newGroup Models.Group) (string, error)
	GetGroup(groupID string) (bson.M, error)
	GetGroups(params Models.GetGroupFilterParams) ([]Models.Group, error)
	UpdateGroup(updatedGroup Models.Group) (bson.M, error)
	DeleteGroup(groupId string) (bson.M, error)
	SubscribeToGroup(subscription Models.GroupSubscription) (bson.M, error)
	UnSubscribeToGroup(subscription Models.GroupSubscription) (bson.M, error)
	PinGroupMessage(pinMessage Models.GroupPinMessage) (bson.M, error)
}

type GroupsMessagesRepo interface {
	CreateGroupMessage(newMessage Models.GroupMessage) (bson.M, error)
	GetGroupMessage(messageId primitive.ObjectID) (bson.M, error)
	GetGroupMessages(params Models.GetGroupMessages) ([]Models.GroupMessage, error)
	UpdateGroupMessage(updatedGroup Models.GroupMessage) (bson.M, error)
	DeleteGroupMessage(groupId string) (bson.M, error)
	CreateGroupMessageComment(newMessageComment Models.GroupMessageComment) (string, error)
	GetGroupMessageComment(messageCommentId string) (bson.M, error)
	GetGroupMessagesComment(params Models.GetGroupMessagesComments) ([]Models.GroupMessageComment, error)
	UpdateGroupMessageComment(updatedGroup Models.GroupMessageComment) (bson.M, error)
	DeleteGroupMessageComment(groupId string) (bson.M, error)
}

type Repo struct {
	Authorization
	DoctorRepo
	ScheduleRepo
	PatientRepo
	ConsultationRepo
	EventRepo
	UserRepo
	MessagesRepo
	ChannelsRepo
	TaskCandidatesRepo
	TaskRepo
	PatientCandidatesRepo
	GroupsRepo
	GroupsMessagesRepo
}

func NewRepo(db *sqlx.DB, mongoDB *mongo.Database, rdb *redis.Client) *Repo {
	return &Repo{
		Authorization:         NewAuthPostgresRepo(db),
		DoctorRepo:            NewDoctorPostgresRepo(db),
		ScheduleRepo:          NewSchedulePostgresRepo(db),
		PatientRepo:           NewPatientPostgresRepo(db),
		ConsultationRepo:      NewConsultationPostgresRepo(db),
		EventRepo:             NewEventPostgresRepo(db),
		UserRepo:              NewUserPostgresRepo(db),
		MessagesRepo:          NewMongoMessagesRepo(mongoDB),
		ChannelsRepo:          NewMongoChannelsRepo(mongoDB),
		TaskCandidatesRepo:    NewTaskCandidatesRedisRepo(rdb),
		TaskRepo:              NewTaskPostgresRepo(db),
		PatientCandidatesRepo: NewPatientCandidatesRepo(mongoDB),
		GroupsRepo:            NewMongoGroupRepo(mongoDB),
		GroupsMessagesRepo:    NewMongoGroupMessagesRepo(mongoDB),
	}
}
