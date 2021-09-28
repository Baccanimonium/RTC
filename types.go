package RTC

type BroadcastingMessage struct {
	MessageType string `bson:"type" json:"type"`
	Payload     map[string]interface{}
}
