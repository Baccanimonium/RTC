package SocketController

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"video-chat-app/src"
)

const ChannelsDocument = "channels"

func (h Hub) createChannel(userId int, payload map[string]interface{}) (interface{}, error) {

	existedChannel, err := h.getChannelByParticipants(userId, payload)

	// when we don't have existed channel and error -> create new channel
	if existedChannel == nil && err == nil {

		channelsCollection := h.mongo.Collection(ChannelsDocument)

		r, err := channelsCollection.InsertOne(context.TODO(), bson.D{
			{Key: "creator", Value: userId},
			{Key: "participant", Value: payload["userId"]},
		})

		if err != nil {
			logrus.Print("failed to create channel: %s", err.Error())
		}

		return h.getChannelByID(r.InsertedID)

	}

	return existedChannel, err

}

func (h Hub) getChannelByParticipants(userId int, payload map[string]interface{}) (bson.M, error) {
	var result bson.M
	channelsCollection := h.mongo.Collection(ChannelsDocument)

	cursor := channelsCollection.FindOne(context.TODO(), bson.M{
		"$and": bson.A{
			bson.M{"$or": bson.A{
				bson.M{"creator": userId},
				bson.M{"participant": userId},
			}},
			bson.M{"$or": bson.A{
				bson.M{"creator": payload["userId"]},
				bson.M{"participant": payload["userId"]},
			}},
		},
	})

	if cursor.Err() != nil {
		logrus.Print("failed to get channel: ", cursor.Err())

		return nil, nil
	}

	err := cursor.Decode(&result)

	if err != nil {
		logrus.Print("failed to decode gotten channel: ", err.Error())
	}

	return result, err
}

func (h Hub) getChannelByID(documentId interface{}) (bson.M, error) {
	var result bson.M
	channelsCollection := h.mongo.Collection(ChannelsDocument)

	cursor := channelsCollection.FindOne(context.TODO(), bson.M{"_id": documentId})

	if cursor.Err() != nil {
		logrus.Print("failed to get channel: %s", cursor.Err())

		return nil, cursor.Err()
	}

	err := cursor.Decode(&result)

	if err != nil {
		logrus.Print("failed to get channel: %s", err.Error())
		return nil, err
	}

	return result, err
}

type Channel struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Creator     int32              `bson:"creator" json:"creator"`
	Participant int32              `bson:"participant" json:"participant"`
}

func (h Hub) GetAllChannelsBelongsToUser(c *gin.Context) {
	userId, userCtxError := c.Get(src.UserContext)

	if !userCtxError {
		c.JSON(http.StatusBadRequest, "user does not exist")
		return
	}

	var result = make([]Channel, 0)
	channelsCollection := h.mongo.Collection(ChannelsDocument)

	cursor, err := channelsCollection.Find(context.TODO(), bson.M{"$or": bson.A{
		bson.M{"creator": userId},
		bson.M{"participant": userId},
	}})

	if cursor.Err() != nil {
		logrus.Print("failed to get channels: ", cursor.Err())
	}

	err = cursor.All(context.TODO(), &result)

	if err != nil {
		logrus.Print("failed to decode gotten channels: ", err.Error())
	}
	logrus.Print(result, reflect.TypeOf(result))
	c.JSON(http.StatusOK, result)

}
