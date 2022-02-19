package SocketHandlers

import (
	"encoding/json"
	"video-chat-app/src/Models"
)

func (c *Client) CreateGroupMessage(rawJson []byte) error {
	var message Models.GroupMessage
	if err := json.Unmarshal(rawJson, &message); err != nil {
		return err
	}

	_, err := c.services.GroupMessagesService.CreateGroupMessage(message)
	return err
}

func (c *Client) UpdateGroupMessage(rawJson []byte) error {
	var message Models.GroupMessage
	if err := json.Unmarshal(rawJson, &message); err != nil {
		return err
	}

	_, err := c.services.GroupMessagesService.UpdateGroupMessage(message)

	return err
}

func (c *Client) DeleteGroupMessage(rawJson []byte) error {
	var groupId string
	if err := json.Unmarshal(rawJson, &groupId); err != nil {
		return err
	}
	return c.services.GroupMessagesService.DeleteGroupMessage(groupId)
}

func (c *Client) createGroupMessageComment(rawJson []byte) error {
	var message Models.GroupMessageComment
	if err := json.Unmarshal(rawJson, &message); err != nil {
		return err
	}

	_, err := c.services.GroupMessagesService.CreateGroupMessageComment(message)
	return err
}

func (c *Client) updateGroupMessageComment(rawJson []byte) error {
	var message Models.GroupMessageComment
	if err := json.Unmarshal(rawJson, &message); err != nil {
		return err
	}

	_, err := c.services.GroupMessagesService.UpdateGroupMessageComment(message)

	return err
}

func (c *Client) deleteGroupMessageComment(groupId string) error {
	_, err := c.services.GroupMessagesService.DeleteGroupMessageComment(groupId)

	return err
}
