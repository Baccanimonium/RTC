package Models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Channel struct {
	Id           primitive.ObjectID `bson:"_id" json:"id"`
	Participants []int              `bson:"participants" json:"participants"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
}
