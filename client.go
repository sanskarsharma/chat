package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Send pings to peer with this period
	pingPeriod = 10 * time.Second
	
	// Time allowed to read the next ping message from the peer. Must be a greater than pingPeriod
	pongWait = (pingPeriod * 11) / 10

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
	connection *websocket.Conn // The websocket connection.
	send chan Message // Buffered channel of messages to be written to connection (i.e websocket)
	roomId string
	id string
}

// readPump pumps messages from the websocket connection to the hub. i.e from client to server/hub
func (subscriber *Subscriber) websocketReader() {
	defer func() {
		leavingMessage := Message{
			LeavingChat: true,
			RoomId: subscriber.roomId,
			SenderId: subscriber.id,
			SentAt: time.Now(),
		}
		h.broadcast <- leavingMessage
		h.unregister <- subscriber
		subscriber.connection.Close()
	}()
	subscriber.connection.SetReadLimit(maxMessageSize)
	subscriber.connection.SetReadDeadline(time.Now().Add(pongWait))

	// setting pong handler function to respond to ping messages 
	subscriber.connection.SetPongHandler(func(string) error { 
		// extending read deadline of websocket connection on receiving the ping message
		subscriber.connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil 
	})

	joiningMessage := Message{
		JoiningChat: true,
		RoomId: subscriber.roomId,
		SenderId: subscriber.id,
		SentAt: time.Now(),
	}
	h.broadcast <- joiningMessage

	for {
		var message Message
		err := subscriber.connection.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		message.RoomId = subscriber.roomId
		message.SentAt = time.Now()
		h.broadcast <- message
	}
}


// websocketWriter pumps messages from the hub to the websocket connection. i.e from server to client
func (subscriber *Subscriber) websocketWriter() {

	// starting a ticker to periodically send ping messages to peers
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		subscriber.connection.Close()
	}()

	for {
		select {
		case message, ok := <-subscriber.send:
			if !ok {
				subscriber.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := subscriber.connection.WriteJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			// extending write deadline and sending ping message
			subscriber.connection.SetWriteDeadline(time.Now().Add(pingPeriod)) 
			if err := subscriber.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
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
	subscriber := &Subscriber{send: make(chan Message, 256), connection: ws, roomId: roomId, id: r.URL.Query().Get("sender_id")}
	h.register <- subscriber

	// spiing up 2 go-routines per client
	go subscriber.websocketReader() // for reading from websocket and broadcasting to hub
	go subscriber.websocketWriter() // for reading from send channel and writing to websocket
}
