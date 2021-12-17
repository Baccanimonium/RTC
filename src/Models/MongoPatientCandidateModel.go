package Models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PatientCandidate struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	UserId      int                `bson:"user_id" json:"user_id"`
	Tags        []int              `bson:"tags" json:"tags"`
	Lat         float64            `bson:"lat" json:"lat"`
	Long        float64            `bson:"long" json:"long"`
	Description string             `bson:"description" json:"description"`
}
