package Repos

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"video-chat-app/src/Models"
)

const ChannelsCollection = "channels"

// переделать участников на список, добавить проверки  в чат при получении сообщений на вхождение в список

func NewMongoChannelsRepo(db *mongo.Database) *Mongo {
	return &Mongo{db: db}
}

func (m *Mongo) CreateChannel(payload Models.Channel) (Models.Channel, error) {

	channelsCollection := m.db.Collection(ChannelsCollection)

	r, err := channelsCollection.InsertOne(context.TODO(), bson.D{
		{Key: "participants", Value: payload.Participants},
	})

	if err != nil {
		logrus.Print("failed to create channel: %s", err.Error())
	}

	return m.GetChannelByID(r.InsertedID)
}

func (m *Mongo) GetChannelByID(documentId interface{}) (Models.Channel, error) {
	var result Models.Channel
	channelsCollection := m.db.Collection(ChannelsCollection)

	cursor := channelsCollection.FindOne(context.TODO(), bson.M{"_id": documentId})

	if cursor.Err() != nil {
		logrus.Print("failed to get channel: %s", cursor.Err())

		return result, cursor.Err()
	}

	err := cursor.Decode(&result)

	if err != nil {
		logrus.Print("failed to get channel: %s", err.Error())
		return result, err
	}

	return result, err
}

func (m *Mongo) GetChannelByParticipants(userId int, payload map[string]interface{}) (Models.Channel, error) {
	var result Models.Channel
	channelsCollection := m.db.Collection(ChannelsCollection)

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

		return Models.Channel{}, nil
	}

	err := cursor.Decode(&result)

	if err != nil {
		logrus.Print("failed to decode gotten channel: ", err.Error())
	}

	return result, err
}

func (m *Mongo) GetAllChannelsBelongsToUser(userId int) ([]Models.Channel, error) {
	var result = make([]Models.Channel, 0)
	channelsCollection := m.db.Collection(ChannelsCollection)

	cursor, err := channelsCollection.Find(context.TODO(), bson.M{"participants": userId})

	if cursor.Err() != nil {
		logrus.Print("failed to get channels: ", cursor.Err())
		return result, err
	}

	err = cursor.All(context.TODO(), &result)

	if err != nil {
		logrus.Print("failed to decode gotten channels: ", err.Error())
	}

	return result, err
}

func (m *Mongo) SubscribeToChannel(userId int, channel Models.Channel) (bson.M, error) {
	channelsCollection := m.db.Collection(ChannelsCollection)

	result := channelsCollection.FindOneAndUpdate(context.TODO(), bson.M{
		"_id": channel.Id,
		"$or": bson.A{
			bson.M{"creator": userId},
			bson.M{"participant": userId},
		},
	}, bson.M{"$addToSet": bson.M{"participants": userId}})

	if result.Err() != nil {
		logrus.Print("failed to create message: %s", result.Err())
		return nil, result.Err()
	}

	var nextMessage bson.M
	err := result.Decode(&nextMessage)

	if err != nil {
		logrus.Print("fail to decode updated message")
		return nil, err
	}
	return bson.M{"_id": channel.Id, "creator": channel.Creator, "participant": channel.Participant}, nil
}

func (m *Mongo) DeleteChannel(userId int, channel Models.Channel) (bson.M, error) {
	channelsCollection := m.db.Collection(ChannelsCollection)

	_, err := channelsCollection.DeleteOne(context.TODO(), bson.M{
		"_id": channel.Id,
		"$or": bson.A{
			bson.M{"creator": userId},
			bson.M{"participant": userId},
		},
	})

	if err != nil {
		logrus.Print("error occur during channel deleting")
		return nil, err
	}

	return bson.M{"_id": channel.Id, "creator": channel.Creator, "participant": channel.Participant}, nil
}
