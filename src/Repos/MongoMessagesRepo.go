package Repos

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"video-chat-app/src/Models"
)

const MessagesCollection = "channels"

func NewMongoMessagesRepo(db *mongo.Database) *Mongo {
	return &Mongo{db: db}
}

func (m *Mongo) CreateMessage(newMessage Models.Message) (bson.D, error) {
	messagesCollection := m.db.Collection(MessagesCollection)

	r, err := messagesCollection.InsertOne(context.TODO(), newMessage)

	if err != nil {
		logrus.Print("failed to create message: %s", err.Error())
		return nil, err
	}

	return m.GetMessage(r.InsertedID)
}

func (m *Mongo) GetMessage(messageId interface{}) (bson.D, error) {
	var result bson.D
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
