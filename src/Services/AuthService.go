package Services

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"video-chat-app/src/Repos"
)

const (
	salt        = "asdasdasd123fdsgdfg"
	signingSalt = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL    = 12 * time.Hour
)

type TokenClaims struct {
	jwt.StandardClaims
	IdUser int `json:"id_user"`
}

type AuthService struct {
	repo Repos.Authorization
}

func NewAuthService(repo Repos.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user Repos.UserCreate) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(user Repos.UserLogin) (string, error) {
	userObj, err := s.repo.GetUser(user.Login, generatePasswordHash(user.Password))

	if err != nil {
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userObj.Id,
	})
	return token.SignedString([]byte(signingSalt))
}

func (s *AuthService) ParseToken(rawToken string) (int, error) {
	token, err := jwt.ParseWithClaims(rawToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingSalt), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*TokenClaims)

	if !ok {
		return 0, errors.New("token claims are not of type *TokenClaims")
	}

	return claims.IdUser, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
