package Models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id        primitive.ObjectID  `bson:"_id" json:"id"`
	Creator   int32               `bson:"creator" json:"creator"`
	ChannelId primitive.ObjectID  `bson:"channelId" json:"channelId"`
	Created   primitive.Timestamp `bson:"created" json:"created"`
	Updated   primitive.Timestamp `bson:"updated" json:"updated"`
	Text      string              `bson:"text" json:"text"`
	Files     bson.A              `bson:"files" json:"files"`
}
