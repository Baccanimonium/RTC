package Services

import (
	"video-chat-app/src/Repos"
)

type Authorization interface {
	CreateUser(user Repos.User) (int, error)
}

type Services struct {
	Authorization
}

func NewService(repo *Repos.Repo) *Services {
	return &Services{
		Authorization: NewAuthService(repo.Authorization),
	}
}
