package Services

import (
	"go.mongodb.org/mongo-driver/bson"
	RTC "video-chat-app"
	"video-chat-app/src/Models"
	"video-chat-app/src/Repos"
)

type ChannelsRepo struct {
	repo Repos.ChannelsRepo
	b    chan RTC.BroadcastingMessage
}

func NewChannelsService(repo Repos.ChannelsRepo, broadcast chan RTC.BroadcastingMessage) *ChannelsRepo {
	return &ChannelsRepo{repo: repo, b: broadcast}
}

func (s *ChannelsRepo) CreateChannel(userId int, payload Models.Channel) (bson.M, error) {
	channel, err := s.repo.CreateChannel(userId, payload)

	if err != nil {
		return nil, err
	}

	s.b <- RTC.BroadcastingMessage{
		MessageType: RTC.BroadcastCreateChannel,
		Payload:     channel,
	}

	return channel, err
}

func (s *ChannelsRepo) DeleteChannel(userId int, payload Models.Channel) (bson.M, error) {
	channel, err := s.repo.DeleteChannel(userId, payload)

	if err != nil {
		return nil, err
	}

	s.b <- RTC.BroadcastingMessage{
		MessageType: RTC.BroadcastDeleteChannel,
		Payload:     channel,
	}

	return channel, err
}

func (s *ChannelsRepo) GetChannelByID(documentId interface{}) (bson.M, error) {
	return s.repo.GetChannelByID(documentId)
}
func (s *ChannelsRepo) GetChannelByParticipants(userId int, payload map[string]interface{}) (Models.Channel, error) {
	return s.repo.GetChannelByParticipants(userId, payload)
}
func (s *ChannelsRepo) GetAllChannelsBelongsToUser(userId int) ([]Models.Channel, error) {
	return s.repo.GetAllChannelsBelongsToUser(userId)
}
