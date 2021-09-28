package SocketHandlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"time"
	"video-chat-app/src/Services"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	//
	hub *Hub

	// services
	services *Services.Services

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	userId int
}

type SocketClientFactory struct {
	services *Services.Services
	hub      *Hub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewSocketClientFactory(services *Services.Services, hub *Hub) *SocketClientFactory {
	return &SocketClientFactory{
		hub:      hub,
		services: services,
	}
}

func (s SocketClientFactory) OnNewSocketClient(context *gin.Context) {
	//logrus.Print(context.Get(RTC.UserContext))
	userId, err := Services.ParseToken(context.Query("authorization"))
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		logrus.Print(err)
		return
	}
	client := &Client{hub: s.hub, conn: conn, send: make(chan []byte, 256), userId: userId, services: s.services}
	client.hub.register <- client
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		var m map[string]interface{}
		json.Unmarshal(message, &m)
		var response = make(map[string]interface{})
		response["messageId"] = m["messageId"]
		rawJson, _ := json.Marshal(m["payload"])
		logrus.Print("socket", response)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch m["type"] {
		case createChatChannel:
			payload, err := c.CreateChannel(c.userId, rawJson)
			if err != nil {
				response["error"] = err.Error()
			}
			response["payload"] = payload
			break

		case createChatMessage:
			payload, err := c.createMessage(rawJson)

			if err != nil {
				response["error"] = err.Error()
			}
			response["payload"] = payload
			break

		default:
			//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			//c.hub.broadcast <- message
		}
		rawMessage, err := json.Marshal(response)
		c.send <- rawMessage
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
