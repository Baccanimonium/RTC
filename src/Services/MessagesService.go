package Services

import (
	"go.mongodb.org/mongo-driver/bson"
	"video-chat-app"
	"video-chat-app/src/Models"
	"video-chat-app/src/Repos"
)

type MessagesRepo struct {
	repo Repos.MessagesRepo
	b    chan RTC.BroadcastingMessage
}

func NewMessagesService(repo Repos.MessagesRepo, broadcast chan RTC.BroadcastingMessage) *MessagesRepo {
	return &MessagesRepo{repo: repo, b: broadcast}
}

func (s *MessagesRepo) CreateMessage(newMessage Models.CreateMessage) (bson.M, error) {
	createdMessage, err := s.repo.CreateMessage(newMessage)
	if err == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastCreateChatMessage,
			Payload:     createdMessage,
		}
	}

	return createdMessage, err
}

func (s *MessagesRepo) GetMessage(messageId interface{}) (bson.M, error) {
	return s.repo.GetMessage(messageId)
}

func (s *MessagesRepo) GetMessages(channelId string, userId interface{}) ([]Models.Message, error) {
	return s.repo.GetMessages(channelId)
}

func (s *MessagesRepo) UpdateMessage(updatedMessage Models.Message, userId int) (bson.M, error) {

	nextMessage, err := s.repo.UpdateMessage(updatedMessage, userId)

	if err == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastUpdateChatMessage,
			Payload:     nextMessage,
		}
	}

	return nextMessage, err
}

func (s *MessagesRepo) DeleteMessage(dMessage Models.DeleteMessage, userId int) (bson.M, error) {

	deletedMessage, err := s.repo.DeleteMessage(dMessage, userId)

	if err == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastDeleteChatMessage,
			Payload:     deletedMessage,
		}
	}

	return deletedMessage, err
}
