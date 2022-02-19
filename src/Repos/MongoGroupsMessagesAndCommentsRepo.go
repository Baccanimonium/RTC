package Repos

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"video-chat-app/src/Models"
)

/*
	Сообщения в группы пишет только автор кнала они всем падают по сокету
	Комменты в группы может писать кто угодно
	комменты будут рефрешиться с фронта раз в 30 секунд, отправка сообщения в комменты будет откладывать таймер рефреша
	у комментов должна быть система лайков как в ютубе
*/

const GroupMessagesCollection = "groupMessages"
const GroupMessagesCommentsCollection = "groupMessagesComments"

func NewMongoGroupMessagesRepo(db *mongo.Database) *Mongo {
	return &Mongo{db: db}
}

func (m *Mongo) CreateGroupMessage(newMessage Models.GroupMessage) (bson.M, error) {
	groupCollection := m.db.Collection(GroupMessagesCollection)

	newMessage.Created = time.Now().Unix()
	newMessage.Updated = time.Now().Unix()

	r, err := groupCollection.InsertOne(context.TODO(), newMessage)

	if err != nil {
		logrus.Print("failed to create group: %s", err.Error())
		return nil, err
	}

	oid, ok := r.InsertedID.(primitive.ObjectID)

	if ok {
		return m.GetGroupMessage(oid)
	}
	return nil, err
}

func (m *Mongo) GetGroupMessage(messageId primitive.ObjectID) (bson.M, error) {
	var result bson.M
	groupMessagesCollection := m.db.Collection(GroupMessagesCollection)
	ctx := context.TODO()
	cursor := groupMessagesCollection.FindOne(ctx, bson.M{"_id": messageId})

	if cursor.Err() != nil {
		logrus.Print("failed to get groupMessage: ", cursor.Err())
		return result, cursor.Err()
	}

	err := cursor.Decode(&result)

	if err != nil {
		logrus.Print("failed to decode groupMessage: ", err.Error())
	}

	return result, err
}

func (m *Mongo) GetGroupMessages(params Models.GetGroupMessages) ([]Models.GroupMessage, error) {
	var result = make([]Models.GroupMessage, 0)
	groupMessagesCollection := m.db.Collection(GroupMessagesCollection)
	ctx := context.TODO()
	objID, _ := primitive.ObjectIDFromHex(params.GroupId)

	opt := options.FindOptions{
		Limit: params.Pagination.Limit,
		Skip:  params.Pagination.Skip,
	}

	cursor, err := groupMessagesCollection.Find(ctx, bson.M{"_id": objID}, &opt)

	if cursor.Err() != nil {
		logrus.Print("failed to get groups: ", cursor.Err())
		return result, err
	}

	for cursor.Next(ctx) {
		var groupMessage Models.GroupMessage
		err := cursor.Decode(&groupMessage)
		if err != nil {
			logrus.Print("failed to decode gotten groupMessage: ", err.Error())
		} else {
			result = append(result, groupMessage)
		}
	}

	if err != nil {
		logrus.Print("failed to decode gotten groupMessages: ", err.Error())
	}

	return result, err
}

func (m *Mongo) UpdateGroupMessage(updatedGroup Models.GroupMessage) (bson.M, error) {
	result := m.db.Collection(GroupMessagesCollection).FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": updatedGroup.Id},
		bson.M{
			"$set": bson.M{
				"text":    updatedGroup.Text,
				"files":   updatedGroup.Files,
				"updated": time.Now().Unix(),
			},
		},
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

func (m *Mongo) DeleteGroupMessage(groupId string) (bson.M, error) {
	var result bson.M
	docID, err := primitive.ObjectIDFromHex(groupId)
	groupMessagesCollection := m.db.Collection(GroupMessagesCollection)
	commentsCollection := m.db.Collection(GroupMessagesCommentsCollection)
	err = m.db.Client().UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		cursor := groupMessagesCollection.FindOneAndDelete(sessionContext, bson.M{"_id": docID})

		if cursor.Err() != nil {
			logrus.Print("failed to get group: ", cursor.Err())
			return cursor.Err()
		}

		err = cursor.Decode(&result)

		_, err = commentsCollection.DeleteMany(sessionContext, bson.M{"groupMessageID": docID})

		if err != nil {
			logrus.Print("failed to delete groupMessagesComments: ", err.Error())
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		return nil
	})

	return result, err
}

func (m *Mongo) CreateGroupMessageComment(newMessageComment Models.GroupMessageComment) (string, error) {
	commentsCollection := m.db.Collection(GroupMessagesCommentsCollection)

	newMessageComment.Created = time.Now().Unix()
	newMessageComment.Updated = time.Now().Unix()

	r, err := commentsCollection.InsertOne(context.TODO(), newMessageComment)

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

func (m *Mongo) GetGroupMessageComment(messageCommentId string) (bson.M, error) {
	var result bson.M
	groupMessagesCollection := m.db.Collection(GroupMessagesCommentsCollection)
	ctx := context.TODO()
	objID, _ := primitive.ObjectIDFromHex(messageCommentId)
	cursor := groupMessagesCollection.FindOne(ctx, bson.M{"_id": objID})

	if cursor.Err() != nil {
		logrus.Print("failed to get groupMessageComment: ", cursor.Err())
		return result, cursor.Err()
	}

	err := cursor.Decode(&result)

	if err != nil {
		logrus.Print("failed to decode groupMessageComment: ", err.Error())
	}

	return result, err
}

func (m *Mongo) GetGroupMessagesComment(params Models.GetGroupMessagesComments) ([]Models.GroupMessageComment, error) {
	var result = make([]Models.GroupMessageComment, 0)
	groupMessagesCollection := m.db.Collection(GroupMessagesCollection)
	ctx := context.TODO()
	objID, _ := primitive.ObjectIDFromHex(params.MessageId)

	opt := options.FindOptions{
		Limit: params.Pagination.Limit,
		Skip:  params.Pagination.Skip,
	}

	cursor, err := groupMessagesCollection.Find(ctx, bson.M{"_id": objID}, &opt)

	if cursor.Err() != nil {
		logrus.Print("failed to get groups: ", cursor.Err())
		return result, err
	}

	for cursor.Next(ctx) {
		var groupMessageComments Models.GroupMessageComment
		err := cursor.Decode(&groupMessageComments)
		if err != nil {
			logrus.Print("failed to decode gotten groupMessageComment: ", err.Error())
		} else {
			result = append(result, groupMessageComments)
		}
	}

	if err != nil {
		logrus.Print("failed to decode gotten groupMessageComments: ", err.Error())
	}

	return result, err
}

func (m *Mongo) UpdateGroupMessageComment(updatedGroup Models.GroupMessageComment) (bson.M, error) {
	result := m.db.Collection(GroupMessagesCommentsCollection).FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": updatedGroup.Id},
		bson.M{
			"$set": bson.M{
				"text":    updatedGroup.Text,
				"files":   updatedGroup.Files,
				"updated": time.Now().Unix(),
			},
		},
	)

	if result.Err() != nil {
		logrus.Print("failed to update messageComment: %s", result.Err())
		return nil, result.Err()
	}

	var nextGroup bson.M
	err := result.Decode(&nextGroup)

	if err != nil {
		logrus.Print("fail to decode updated messageComment")
		return nil, err
	}

	return nextGroup, nil
}

func (m *Mongo) DeleteGroupMessageComment(groupId string) (bson.M, error) {
	var result bson.M
	docID, err := primitive.ObjectIDFromHex(groupId)
	groupCollection := m.db.Collection(GroupMessagesCommentsCollection)

	cursor := groupCollection.FindOneAndDelete(context.TODO(), bson.M{"_id": docID})

	if cursor.Err() != nil {
		logrus.Print("failed to get messageComment: ", cursor.Err())
		return result, cursor.Err()
	}

	err = cursor.Decode(&result)

	if err != nil {
		logrus.Print("failed to decode messageComment: ", err.Error())
	}

	return result, err
}
