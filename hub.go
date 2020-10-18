package main

type Message struct {
	Data string `json:"data"`
	RoomId string `json:"room_id"`
	SenderId string `json:"sender_id"`
}

type RoomSubscription struct {
	subscriber *Subscriber
	roomId string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*Subscriber]bool

	// Inbound messages from the connections.
	broadcast chan Message

	// Register requests from the connections.
	register chan RoomSubscription

	// Unregister requests from connections.
	unregister chan RoomSubscription
}

var h = hub{
	broadcast:  make(chan Message),
	register:   make(chan RoomSubscription),
	unregister: make(chan RoomSubscription),
	rooms:      make(map[string]map[*Subscriber]bool),
}

func (h *hub) run() {
	for {
		select {
		case roomSubscription := <-h.register:
			subscribers := h.rooms[roomSubscription.roomId]
			if subscribers == nil {
				subscribers = make(map[*Subscriber]bool)
				h.rooms[roomSubscription.roomId] = subscribers
			}
			h.rooms[roomSubscription.roomId][roomSubscription.subscriber] = true
		case roomSubscription := <-h.unregister:
			subscribers := h.rooms[roomSubscription.roomId]
			if subscribers != nil {
				if _, ok := subscribers[roomSubscription.subscriber]; ok {
					delete(subscribers, roomSubscription.subscriber)
					close(roomSubscription.subscriber.send)
					if len(subscribers) == 0 {
						delete(h.rooms, roomSubscription.roomId)
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
