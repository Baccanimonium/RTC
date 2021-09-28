package Services

import (
	"errors"
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

func (s *MessagesRepo) UpdateMessage(updatedMessage Models.Message, userId interface{}) (bson.M, error) {
	// TODO переделать проверки, на проверки в базе. типо искать и по _id и по creator
	message, err := s.repo.GetMessage(updatedMessage.Id)
	if err != nil || message["creator"] != userId {
		return nil, errors.New("creators ids comparison failed")
	}

	nextMessage, err := s.repo.UpdateMessage(updatedMessage)

	if err == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastUpdateChatMessage,
			Payload:     nextMessage,
		}
	}

	return nextMessage, err
}

func (s *MessagesRepo) DeleteMessage(dMessage Models.DeleteMessage, userId interface{}) (bson.M, error) {
	message, err := s.repo.GetMessage(dMessage.Id)
	if err != nil || message["creator"] != userId {
		return nil, errors.New("creators ids comparison failed")
	}

	deletedMessage, err := s.repo.DeleteMessage(dMessage)

	if err == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastDeleteChatMessage,
			Payload:     deletedMessage,
		}
	}

	return deletedMessage, err
}
