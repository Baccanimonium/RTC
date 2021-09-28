package Models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Channel struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Creator     int32              `bson:"creator" json:"creator"`
	Participant int32              `bson:"participant" json:"participant" binding:"required"`
}
