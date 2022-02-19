package Models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Creator   int                `bson:"creator" json:"creator"`
	ChannelId primitive.ObjectID `bson:"channelId" json:"channelId"`
	Created   string             `bson:"created" json:"created"`
	Updated   string             `bson:"updated" json:"updated"`
	Text      string             `bson:"text" json:"text"`
	Files     bson.A             `bson:"files,omitempty" json:"files"`
}

type CreateMessage struct {
	Creator   int                `bson:"creator" json:"creator"`
	ChannelId primitive.ObjectID `bson:"channelId" json:"channelId"`
	Created   string             `bson:"created" json:"created"`
	Updated   string             `bson:"updated" json:"updated"`
	Text      string             `bson:"text" json:"text"`
	Files     bson.A             `bson:"files,omitempty" json:"files"`
}

type DeleteMessage struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	ChannelId primitive.ObjectID `bson:"channelId" json:"channelId"`
	Creator   int                `bson:"creator" json:"creator"`
}
