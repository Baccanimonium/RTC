package Repos

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"video-chat-app/src/Models"
)

const MessagesCollection = "messages"

func NewMongoMessagesRepo(db *mongo.Database) *Mongo {
	return &Mongo{db: db}
}

func (m *Mongo) CreateMessage(newMessage Models.CreateMessage) (bson.M, error) {
	messagesCollection := m.db.Collection(MessagesCollection)
	newMessage.Created = time.Now().Format("02.01.2006 15:04:05")
	newMessage.Updated = time.Now().Format("02.01.2006 15:04:05")

	r, err := messagesCollection.InsertOne(context.TODO(), newMessage)

	if err != nil {
		logrus.Print("failed to create message: %s", err.Error())
		return nil, err
	}
	logrus.Print(r.InsertedID)
	return m.GetMessage(r.InsertedID)
}

func (m *Mongo) UpdateMessage(updatedMessage Models.Message, userId int) (bson.M, error) {
	messagesCollection := m.db.Collection(MessagesCollection)
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	result := messagesCollection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": updatedMessage.Id, "creator": userId},
		bson.M{
			"$set": bson.M{
				"text":    updatedMessage.Text,
				"files":   updatedMessage.Files,
				"updated": time.Now().Format("02.01.2006 15:04:05"),
			},
		},
		&opt,
	)

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

	return nextMessage, nil
}

func (m *Mongo) DeleteMessage(message Models.DeleteMessage, userId int) (bson.M, error) {
	messagesCollection := m.db.Collection(MessagesCollection)

	_, err := messagesCollection.DeleteOne(context.TODO(), bson.M{"_id": message.Id, "creator": userId})

	if err != nil {
		logrus.Print("error occur during message deleting")
		return nil, err
	}

	return bson.M{"_id": message.Id, "channelId": message.ChannelId, "creator": message.Creator}, nil
}

func (m *Mongo) GetMessage(messageId interface{}) (bson.M, error) {
	var result bson.M
	messagesCollection := m.db.Collection(MessagesCollection)
	cursor := messagesCollection.FindOne(context.TODO(), bson.M{"_id": messageId})

	if cursor.Err() != nil {
		logrus.Print("failed to get a message: %s", cursor.Err())

		return nil, cursor.Err()
	}

	err := cursor.Decode(&result)

	if err != nil {
		logrus.Print("failed to get a message: %s", err.Error())
		return nil, err
	}

	return result, err
}

func (m *Mongo) GetMessages(channelId string) ([]Models.Message, error) {
	var result = make([]Models.Message, 0)
	messagesCollection := m.db.Collection(MessagesCollection)
	ctx := context.TODO()
	objID, _ := primitive.ObjectIDFromHex(channelId)
	cursor, err := messagesCollection.Find(ctx, bson.M{"channelId": objID})

	if cursor.Err() != nil {
		logrus.Print("failed to get messages: ", cursor.Err())
		return result, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var message Models.Message
		cursor.Decode(&message)
		result = append(result, message)

	}

	if err != nil {
		logrus.Print("failed to decode gotten messages: ", err.Error())
	}

	return result, err
}
