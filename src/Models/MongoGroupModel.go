package Models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetGroupFilterParams struct {
	Tags       []string
	Pagination MongoPagination
}

type Group struct {
	Id                string  `bson:"_id,omitempty" json:"id"`
	Participants      []int   `bson:"participants,omitempty" json:"participants"`
	ParticipantsCount int     `bson:"participants_count,omitempty" json:"participants_count"`
	OwnerId           int     `bson:"owner_id" json:"owner_id"`
	SubscriptionCost  float32 `bson:"subscription_cost" json:"subscription_cost"`
	Tags              []int   `bson:"tags" json:"tags"`
	Description       string  `bson:"description" json:"description"`
	Name              string  `bson:"name" json:"name"`
	PinnedMessageId   string  `bson:"pinned_message_id,omitempty" json:"pinned_message_id"`
	Deleted           bool    `bson:"deleted,omitempty" json:"deleted"`
}

type GroupFiles struct {
	Id        string `bson:"_id,omitempty" json:"id"`
	MessageId string `bson:"message_id" json:"message_id"`
	GroupId   string `bson:"group_id" json:"group_id"`
	File      string `bson:"file" json:"file"`
}

type GroupSubscription struct {
	GroupId       string `bson:"_id" json:"id"`
	ParticipantId int    `bson:"participant_id" json:"participant_id"`
}

type GroupPinMessage struct {
	GroupId         string `bson:"_id" json:"id"`
	PinnedMessageId int    `bson:"pinned_message_id" json:"pinned_message_id"`
}

type GroupMessage struct {
	Id      primitive.ObjectID `bson:"_id" json:"_id"`
	Creator int32              `bson:"creator" json:"creator"`
	GroupId primitive.ObjectID `bson:"groupId" json:"groupId"`
	Created int64              `bson:"created" json:"created"`
	Updated int64              `bson:"updated" json:"updated"`
	Text    string             `bson:"text" json:"text"`
	Files   bson.A             `bson:"files,omitempty" json:"files"`
}

type GetGroupMessages struct {
	GroupId    string `bson:"groupId" json:"groupId"`
	Pagination MongoPagination
}

type GroupMessageComment struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id"`
	RelationId     primitive.ObjectID `bson:"relation_id" json:"relation_id"`
	Creator        int32              `bson:"creator" json:"creator"`
	GroupMessageID primitive.ObjectID `bson:"groupMessageID" json:"groupMessageID"`
	Created        int64              `bson:"created" json:"created"`
	Updated        int64              `bson:"updated" json:"updated"`
	Text           string             `bson:"text" json:"text"`
	Files          bson.A             `bson:"files,omitempty" json:"files" `
}

type GetGroupMessagesComments struct {
	MessageId  string `bson:"MessageId" json:"MessageId"`
	Pagination MongoPagination
}
