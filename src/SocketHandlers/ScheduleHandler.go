package SocketHandlers

import (
	"encoding/json"
	RTC "video-chat-app"
)

func (h Hub) SendScheduleEvent(message RTC.BroadcastingMessage) {
	receiverId := message.Payload["id_patient"].(int)
	if h.clients[receiverId] != nil {
		rawMessage, err := json.Marshal(message)
		if err == nil {
			h.clients[receiverId].send <- rawMessage
		}
	}
}
