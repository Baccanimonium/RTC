package Repos

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"video-chat-app/src/Models"
)

const patientCandidates = "patientCandidates"

func NewPatientCandidatesRepo(db *mongo.Database) *Mongo {
	return &Mongo{db: db}
}

func (m *Mongo) CreatePatientCandidate(patientCandidate Models.PatientCandidate) (interface{}, error) {
	patientCandidateCollection := m.db.Collection(patientCandidates)

	r, err := patientCandidateCollection.InsertOne(context.TODO(), bson.D{
		{Key: "user_id", Value: patientCandidate.UserId},
		{Key: "tags", Value: patientCandidate.Tags},
		{Key: "lat", Value: patientCandidate.Lat},
		{Key: "long", Value: patientCandidate.Long},
		{Key: "description", Value: patientCandidate.Description},
	})

	if err != nil {
		logrus.Print("failed to create channel: %s", err.Error())
	}

	return r.InsertedID, err
}

func (m *Mongo) GetAllPatientCandidates() ([]Models.PatientCandidate, error) {
	var result = make([]Models.PatientCandidate, 0)
	patientCandidateCollection := m.db.Collection(patientCandidates)

	cursor, err := patientCandidateCollection.Find(context.TODO(), bson.M{})

	if cursor.Err() != nil {
		logrus.Print("failed to get PatientCandidates: ", cursor.Err())
		return result, err
	}

	err = cursor.All(context.TODO(), &result)

	if err != nil {
		logrus.Print("failed to decode gotten PatientCandidates: ", err.Error())
	}

	return result, err
}
