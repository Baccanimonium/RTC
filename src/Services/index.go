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

type Services struct {
	Authorization
	DoctorService
}

func NewService(repo *Repos.Repo) *Services {
	return &Services{
		Authorization: NewAuthService(repo.Authorization),
		DoctorService: NewDoctorService(repo.DoctorRepo),
	}
}
