package SocketHandlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	RTC "video-chat-app"
	"video-chat-app/src/Repos"
)

type Hub struct {
	// Registered clients.
	clients map[int]*Client

	// Inbound messages from the clients.
	broadcast chan RTC.BroadcastingMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	repo *Repos.Repo
}

func NewHub(repo *Repos.Repo, broadcast chan RTC.BroadcastingMessage) *Hub {

	//mongo.Collection("channels").Drop(context.TODO())
	//asd, err := mongo.Collection("channels").EstimatedDocumentCount(context.TODO())
	//
	//if err != nil {
	//	logrus.Fatalf("failed to initialize db: %s", err.Error())
	//}
	//
	//logrus.Print("44doc num ", asd)

	return &Hub{
		broadcast:  broadcast,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int]*Client),
		repo:       repo,
	}
}

func (h *Hub) Run() {

	handlers := map[string]func(message RTC.BroadcastingMessage){
		RTC.BroadcastCreateSchedule: h.SendScheduleEvent,
		RTC.BroadcastUpdateSchedule: h.SendScheduleEvent,
		RTC.BroadcastDeleteSchedule: h.SendScheduleEvent,
	}

	for {
		select {
		case client := <-h.register:
			h.clients[client.userId] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.userId]; ok {
				delete(h.clients, client.userId)
				close(client.send)
			}
		case message := <-h.broadcast:
			logrus.Print("broadcasting", message.MessageType, message.Payload)
			if message.MessageType == RTC.BroadcastCreateChatMessage ||
				message.MessageType == RTC.BroadcastDeleteChatMessage ||
				message.MessageType == RTC.BroadcastUpdateChatMessage {
				channel, err := h.repo.GetChannelByID(message.Payload["channelId"])
				if err != nil {
					logrus.Print("error during messages distribution")
					return
				}

				var receiverId int

				if channel["creator"] == message.Payload["creator"] {
					receiverId = channel["creator"].(int)
				} else {
					receiverId = channel["participant"].(int)
				}

				if h.clients[receiverId] != nil {
					rawMessage, err := json.Marshal(message)
					if err == nil {
						h.clients[receiverId].send <- rawMessage
					}
				}
			}

			handlers[message.MessageType](message)
			//rawJson, _ := json.Marshal(message)
			//if message["channelId"] != nil {
			//	channel, err := h.repo.ChannelsRepo.GetChannelByID(message["channelId"])
			//	if err == nil {
			//		h.clients[channel["participant"].(int)].send <- rawJson
			//	 	h.clients[channel["participant"].(int)].send <- rawJson
			//	}
			//}
		}
	}
}
