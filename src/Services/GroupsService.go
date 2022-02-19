package Services

import (
	"go.mongodb.org/mongo-driver/bson"
	RTC "video-chat-app"
	"video-chat-app/src/Models"
	"video-chat-app/src/Repos"
)

type GroupRepo struct {
	repo Repos.GroupsRepo
	b    chan RTC.BroadcastingMessage
}

func NewGroupService(repo Repos.GroupsRepo, broadcast chan RTC.BroadcastingMessage) *GroupRepo {
	return &GroupRepo{repo: repo, b: broadcast}
}

func (s *GroupRepo) CreateGroup(newGroup Models.Group) (string, error) {
	return s.repo.CreateGroup(newGroup)
}

func (s *GroupRepo) GetGroup(groupID string) (bson.M, error) {
	return s.repo.GetGroup(groupID)
}

func (s *GroupRepo) GetGroups(params Models.GetGroupFilterParams) ([]Models.Group, error) {
	return s.repo.GetGroups(params)
}

func (s *GroupRepo) UpdateGroup(updatedGroup Models.Group) (bson.M, error) {
	return s.repo.UpdateGroup(updatedGroup)
}

func (s *GroupRepo) DeleteGroup(groupId string) (bson.M, error) {
	return s.repo.DeleteGroup(groupId)
}

func (s *GroupRepo) SubscribeToGroup(subscription Models.GroupSubscription) (bson.M, error) {
	return s.repo.SubscribeToGroup(subscription)
}

func (s *GroupRepo) UnSubscribeToGroup(subscription Models.GroupSubscription) (bson.M, error) {
	return s.repo.UnSubscribeToGroup(subscription)
}

func (s *GroupRepo) PinGroupMessage(pinMessage Models.GroupPinMessage) error {
	pinnedMessage, err := s.repo.PinGroupMessage(pinMessage)

	rawGroupMessage, convertError := RTC.ConvertToJson(pinnedMessage)

	if err == nil && convertError == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastPinGroupMessage,
			Payload:     rawGroupMessage,
		}
	}

	return err
}
