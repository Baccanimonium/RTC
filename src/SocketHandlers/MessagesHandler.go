package SocketHandlers

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"video-chat-app/src/Models"
)

func (c *Client) createMessage(rawJson []byte) (bson.M, error) {
	var message Models.CreateMessage
	if err := json.Unmarshal(rawJson, &message); err != nil {
		return bson.M{}, err
	}

	return c.services.MessagesService.CreateMessage(message)

}
