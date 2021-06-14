package Services

import (
	"video-chat-app/src/Repos"
)

type AuthService struct {
	repo Authorization
}

func NewAuthService(repo Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user Repos.User) (int, error) {
	return s.repo.CreateUser(user)
}
