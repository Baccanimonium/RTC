package Services

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	RTC "video-chat-app"
	"video-chat-app/src/Models"
	"video-chat-app/src/Repos"
)

type GroupMessagesRepo struct {
	repo Repos.GroupsMessagesRepo
	b    chan RTC.BroadcastingMessage
}

func NewGroupMessagesService(repo Repos.GroupsMessagesRepo, broadcast chan RTC.BroadcastingMessage) *GroupMessagesRepo {
	return &GroupMessagesRepo{repo: repo, b: broadcast}
}

func (s *GroupMessagesRepo) CreateGroupMessage(newMessage Models.GroupMessage) (bson.M, error) {
	newGroupMessage, err := s.repo.CreateGroupMessage(newMessage)

	rawGroupMessage, convertError := RTC.ConvertToJson(newGroupMessage)

	if err == nil && convertError == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastCreateGroupMessage,
			Payload:     rawGroupMessage,
		}
	}

	return newGroupMessage, err
}

func (s *GroupMessagesRepo) GetGroupMessage(messageId primitive.ObjectID) (bson.M, error) {
	return s.repo.GetGroupMessage(messageId)
}

func (s *GroupMessagesRepo) GetGroupMessages(params Models.GetGroupMessages) ([]Models.GroupMessage, error) {
	return s.repo.GetGroupMessages(params)
}

func (s *GroupMessagesRepo) UpdateGroupMessage(updatedGroup Models.GroupMessage) (bson.M, error) {
	newGroupMessage, err := s.repo.UpdateGroupMessage(updatedGroup)

	rawGroupMessage, convertError := RTC.ConvertToJson(newGroupMessage)

	if err == nil && convertError == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastUpdateGroupMessage,
			Payload:     rawGroupMessage,
		}
	}

	return newGroupMessage, err
}

func (s *GroupMessagesRepo) DeleteGroupMessage(groupId string) error {
	deletedSchedule, err := s.repo.DeleteGroupMessage(groupId)

	rawSchedule, convertError := RTC.ConvertToJson(deletedSchedule)

	if err == nil && convertError == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastDeleteGroupMessage,
			Payload:     rawSchedule,
		}
	}

	return err
}

func (s *GroupMessagesRepo) CreateGroupMessageComment(newMessageComment Models.GroupMessageComment) (string, error) {
	return s.repo.CreateGroupMessageComment(newMessageComment)
}

func (s *GroupMessagesRepo) GetGroupMessageComment(messageCommentId string) (bson.M, error) {
	return s.repo.GetGroupMessageComment(messageCommentId)
}

func (s *GroupMessagesRepo) GetGroupMessagesComment(params Models.GetGroupMessagesComments) ([]Models.GroupMessageComment, error) {
	return s.repo.GetGroupMessagesComment(params)
}

func (s *GroupMessagesRepo) UpdateGroupMessageComment(updatedGroup Models.GroupMessageComment) (bson.M, error) {
	return s.repo.UpdateGroupMessageComment(updatedGroup)
}

func (s *GroupMessagesRepo) DeleteGroupMessageComment(groupId string) (bson.M, error) {
	return s.repo.DeleteGroupMessageComment(groupId)
}
