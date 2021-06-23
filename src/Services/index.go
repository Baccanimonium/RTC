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
	GetAllDoctor() ([]Repos.Participant, error)
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
