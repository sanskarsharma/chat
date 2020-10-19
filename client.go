package main

import (
	"log"
	"net/http"
	"time"

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// connection is an middleman between the websocket connection and the hub.
type Subscriber struct {
	// The websocket connection.
	connection *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Message
}

// write writes a message with the given message type and payload.
func (c *Subscriber) write(mt int, payload []byte) error {
	c.connection.SetWriteDeadline(time.Now().Add(writeWait))
	return c.connection.WriteMessage(mt, payload)
}

// readPump pumps messages from the websocket connection to the hub.
func (roomSubscription RoomSubscription) readPump() {
	subscriber := roomSubscription.subscriber
	defer func() {
		h.unregister <- roomSubscription
		subscriber.connection.Close()
	}()
	subscriber.connection.SetReadLimit(maxMessageSize)
	subscriber.connection.SetReadDeadline(time.Now().Add(pongWait))
	subscriber.connection.SetPongHandler(func(string) error { subscriber.connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message Message
		err := subscriber.connection.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		message.RoomId = roomSubscription.roomId
		h.broadcast <- message
	}
}



// writePump pumps messages from the hub to the websocket connection.
func (roomSubscription *RoomSubscription) writePump() {
	subscriber := roomSubscription.subscriber
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		subscriber.connection.Close()
	}()
	for {
		select {
		case message, ok := <-subscriber.send:
			if !ok {
				subscriber.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := subscriber.connection.WriteJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			if err := subscriber.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, roomId string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	subscriber := &Subscriber{send: make(chan Message, 256), connection: ws}
	roomSubscription := RoomSubscription{subscriber, roomId}
	h.register <- roomSubscription
	go roomSubscription.writePump()
	go roomSubscription.readPump()
}
