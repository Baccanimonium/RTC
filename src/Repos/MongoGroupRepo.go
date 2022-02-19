package Repos

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"video-chat-app/src/Models"
)

/*
	Должны иметь пулл пользователей которые могут ими управлять
	Должны иметь роли
	Должны иметь стоимость подписки
	тэги
*/

const GroupCollection = "group"

func NewMongoGroupRepo(db *mongo.Database) *Mongo {
	return &Mongo{db: db}
}

func (m *Mongo) CreateGroup(newGroup Models.Group) (string, error) {
	groupCollection := m.db.Collection(GroupCollection)

	r, err := groupCollection.InsertOne(context.TODO(), newGroup)

	if err != nil {
		logrus.Print("failed to create group: %s", err.Error())
		return "", err
	}

	oid, ok := r.InsertedID.(primitive.ObjectID)

	if ok {
		return oid.Hex(), nil
	}
	return "", errors.New("failed to convert objectID to hex")
}

func (m *Mongo) GetGroup(groupID string) (bson.M, error) {
	var result bson.M
	groupCollection := m.db.Collection(GroupCollection)
	ctx := context.TODO()
	objID, _ := primitive.ObjectIDFromHex(groupID)
	cursor := groupCollection.FindOne(ctx, bson.M{"_id": objID})

	if cursor.Err() != nil {
		logrus.Print("failed to get group: ", cursor.Err())
		return result, cursor.Err()
	}

	err := cursor.Decode(&result)

	if err != nil {
		logrus.Print("failed to decode group: ", err.Error())
	}

	return result, err
}

func (m *Mongo) GetGroups(params Models.GetGroupFilterParams) ([]Models.Group, error) {
	var result = make([]Models.Group, 0)
	groupCollection := m.db.Collection(GroupCollection)
	ctx := context.TODO()

	tagsBytes, err := bson.Marshal(params.Tags)

	if err != nil {
		return result, fmt.Errorf("failed to marhsal user. error: %v", err)
	}

	//// TODO Проверить рааботоспособность опций, если что применить агрегацию
	//matchStage := bson.D{{"$match", bson.D{{"tags", tagsBytes}}}}
	//limitStage := bson.D{{"$limit", params.Pagination.Limit}}
	//skipStage := bson.D{{"$skip", params.Pagination.Skip}}
	//projectStage := bson.D{{
	//	"$project",
	//	bson.D{{"participants_count", bson.D{{"$size", "$participants"}}}},
	//}}
	//
	//showLoadedCursor, err := groupCollection.Aggregate(ctx, mongo.Pipeline{matchStage, limitStage, skipStage, projectStage})
	//if err != nil {
	//	panic(err)
	//}
	//var showsLoaded []bson.M
	//if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
	//	panic(err)
	//}

	opt := options.FindOptions{
		Limit: params.Pagination.Limit,
		Skip:  params.Pagination.Skip,
		Projection: bson.D{
			{"participants", 0},
			{"pinned_message_id", 0},
			{"participants_count", bson.D{{"$size", "$participants"}}},
		},
	}

	cursor, err := groupCollection.Find(ctx, bson.M{"tags": tagsBytes}, &opt)

	if cursor.Err() != nil {
		logrus.Print("failed to get groups: ", cursor.Err())
		return result, err
	}

	for cursor.Next(ctx) {
		var group Models.Group
		err := cursor.Decode(&group)
		if err != nil {
			logrus.Print("failed to decode gotten group: ", err.Error())
		} else {
			result = append(result, group)
		}
	}

	if err != nil {
		logrus.Print("failed to decode gotten groups: ", err.Error())
	}

	return result, err
}

func (m *Mongo) UpdateGroup(updatedGroup Models.Group) (bson.M, error) {
	opt := options.FindOneAndUpdateOptions{
		Projection: bson.D{
			{"participants", 0},
			{"pinned_message_id", 0},
			{"participants_count", bson.D{{"$size", "$participants"}}},
		},
	}
	objID, _ := primitive.ObjectIDFromHex(updatedGroup.Id)
	result := m.db.Collection(GroupCollection).FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{
			"$set": bson.M{
				"subscription_cost": updatedGroup.SubscriptionCost,
				"tags":              updatedGroup.Tags,
				"description":       updatedGroup.Description,
				"name":              updatedGroup.Name,
			},
		},
		&opt,
	)

	if result.Err() != nil {
		logrus.Print("failed to update group: %s", result.Err())
		return nil, result.Err()
	}

	var nextGroup bson.M
	err := result.Decode(&nextGroup)

	if err != nil {
		logrus.Print("fail to decode updated group")
		return nil, err
	}

	return nextGroup, nil
}

func (m *Mongo) DeleteGroup(groupId string) (bson.M, error) {
	var result bson.M
	docID, err := primitive.ObjectIDFromHex(groupId)
	groupCollection := m.db.Collection(GroupCollection)

	cursor := groupCollection.FindOneAndDelete(context.TODO(), bson.M{"_id": docID})

	if cursor.Err() != nil {
		logrus.Print("failed to get group: ", cursor.Err())
		return result, cursor.Err()
	}

	err = cursor.Decode(&result)

	if err != nil {
		logrus.Print("failed to decode group: ", err.Error())
	}

	return result, err
}

func (m *Mongo) SubscribeToGroup(subscription Models.GroupSubscription) (bson.M, error) {
	// TODO проверить параметр updatedRows
	opt := options.FindOneAndUpdateOptions{
		Projection: bson.D{
			{"participants", 0},
			{"pinned_message_id", 0},
			{"participants_count", bson.D{{"$size", "$participants"}}},
		},
	}

	objID, _ := primitive.ObjectIDFromHex(subscription.GroupId)

	result := m.db.Collection(GroupCollection).FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$addToSet": bson.M{"participants": subscription.ParticipantId}},
		&opt,
	)

	if result.Err() != nil {
		logrus.Print("group subscription has failed: %s", result.Err())
		return nil, result.Err()
	}
	var nextGroup bson.M
	err := result.Decode(&nextGroup)

	if err != nil {
		logrus.Print("fail to decode updated group")
		return nil, err
	}

	return nextGroup, nil
}

func (m *Mongo) UnSubscribeToGroup(subscription Models.GroupSubscription) (bson.M, error) {
	// TODO проверить параметр updatedRows
	opt := options.FindOneAndUpdateOptions{
		Projection: bson.D{
			{"participants", 0},
			{"pinned_message_id", 0},
			{"participants_count", bson.D{{"$size", "$participants"}}},
		},
	}

	objID, _ := primitive.ObjectIDFromHex(subscription.GroupId)

	result := m.db.Collection(GroupCollection).FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$pull": bson.M{"participants": subscription.ParticipantId}},
		&opt,
	)

	if result.Err() != nil {
		logrus.Print("group unsubscription has failed: %s", result.Err())
		return nil, result.Err()
	}
	var nextGroup bson.M
	err := result.Decode(&nextGroup)

	if err != nil {
		logrus.Print("fail to decode updated group")
		return nil, err
	}

	return nextGroup, nil
}

func (m *Mongo) PinGroupMessage(pinMessage Models.GroupPinMessage) (bson.M, error) {
	// TODO проверить параметр updatedRows
	opt := options.FindOneAndUpdateOptions{
		Projection: bson.D{{"pinned_message_id", 1}},
	}

	objID, _ := primitive.ObjectIDFromHex(pinMessage.GroupId)

	result := m.db.Collection(GroupCollection).FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{
			"$set": bson.M{
				"subscription_cost": pinMessage.PinnedMessageId,
			},
		},
		&opt,
	)

	if result.Err() != nil {
		logrus.Print("group subscription has failed: %s", result.Err())
		return nil, result.Err()
	}
	var nextGroup bson.M
	err := result.Decode(&nextGroup)

	if err != nil {
		logrus.Print("fail to decode updated group")
		return nil, err
	}

	return nextGroup, nil
}
