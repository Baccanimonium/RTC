package SocketController

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"video-chat-app/src/Services"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	mongo *mongo.Database
}

// OnNewSocketClient handles websocket requests from the peer.
func (h Hub) OnNewSocketClient(context *gin.Context) {
	//logrus.Print(context.Get(RTC.UserContext))
	userId, err := Services.ParseToken(context.Query("authorization"))
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		logrus.Print(err)
		return
	}
	client := &Client{hub: &h, conn: conn, send: make(chan []byte, 256), userId: userId}
	client.hub.register <- client
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func NewHub(mongo *mongo.Database) *Hub {

	//mongo.Collection("channels").Drop(context.TODO())
	//asd, err := mongo.Collection("channels").EstimatedDocumentCount(context.TODO())
	//
	//if err != nil {
	//	logrus.Fatalf("failed to initialize db: %s", err.Error())
	//}
	//
	//logrus.Print("44doc num ", asd)

	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		mongo:      mongo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
