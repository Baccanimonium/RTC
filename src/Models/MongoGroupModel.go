package Models

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
