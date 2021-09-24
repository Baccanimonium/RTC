package Services

import (
	"go.mongodb.org/mongo-driver/bson"
	"video-chat-app/src/Models"
	"video-chat-app/src/Repos"
)

type MessagesRepo struct {
	repo Repos.MessagesRepo
}

func NewMessagesService(repo Repos.MessagesRepo) *MessagesRepo {
	return &MessagesRepo{repo: repo}
}

func (s *MessagesRepo) CreateMessage(newMessage Models.Message) (bson.D, error) {
	return s.repo.CreateMessage(newMessage)
}

func (s *MessagesRepo) GetMessage(messageId interface{}) (bson.D, error) {
	return s.repo.GetMessage(messageId)
}
