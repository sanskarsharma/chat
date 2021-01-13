package main
import (
	"time"
)

type Message struct {
	Data string `json:"data"`
	RoomId string `json:"room_id"`
	SenderId string `json:"sender_id"`
	SentAt time.Time `json:"sent_at"`
}

// hub maintains the set of active connections and broadcasts messages to the connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*Subscriber]bool

	// Inbound messages from the connections.
	broadcast chan Message

	// Register requests from the connections.
	register chan *Subscriber

	// Unregister requests from connections.
	unregister chan *Subscriber
}

var h = hub{
	broadcast:  make(chan Message),
	register:   make(chan *Subscriber),
	unregister: make(chan *Subscriber),
	rooms:      make(map[string]map[*Subscriber]bool),
}

func (h *hub) run() {
	for {
		select {
		case subscriber := <-h.register:
			roomSubscribers := h.rooms[subscriber.roomId]
			if roomSubscribers == nil {
				roomSubscribers = make(map[*Subscriber]bool)
				h.rooms[subscriber.roomId] = roomSubscribers
			}
			h.rooms[subscriber.roomId][subscriber] = true
		case subscriber := <-h.unregister:
			roomSubscribers := h.rooms[subscriber.roomId]
			if roomSubscribers != nil {
				if _, ok := roomSubscribers[subscriber]; ok {
					delete(roomSubscribers, subscriber)
					close(subscriber.send)
					if len(roomSubscribers) == 0 {
						delete(h.rooms, subscriber.roomId)
					}
				}
			}
		case m := <-h.broadcast:
			subscribers := h.rooms[m.RoomId]
			for c := range subscribers {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(subscribers, c)
					if len(subscribers) == 0 {
						delete(h.rooms, m.RoomId)
					}
				}
			}
		}
	}
}
