package Services

import (
	"crypto/sha1"
	"fmt"
	"video-chat-app/src/Repos"
)

const (
	salt = "asdasdasd123fdsgdfg"
)

type AuthService struct {
	repo Authorization
}

func NewAuthService(repo Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user Repos.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
