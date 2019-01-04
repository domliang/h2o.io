package h2oio

import "time"

// Hub 客户端hub
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	//
	registerHandler func(*Client)

	//
	unRegisterHandler func(*Client)

	//
	messageHandler func(*Client, []byte)

	// Time allowed to write a message to the peer.
	writeWait time.Duration

	// Time allowed to read the next pong message from the peer.
	pongWait time.Duration

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod time.Duration

	// Maximum message size allowed from peer.
	maxMessageSize int64
}

// NewHub buuild new hub
func NewHub(
	registerHandler func(*Client),
	unRegisterHandler func(*Client),
	messageHandler func(*Client, []byte),
	writeWait time.Duration,
	pongWait time.Duration,
	pingPeriod time.Duration,
	maxMessageSize int64) *Hub {
	return &Hub{
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		clients:           make(map[string]*Client),
		registerHandler:   registerHandler,
		messageHandler:    messageHandler,
		unRegisterHandler: unRegisterHandler,
		writeWait:         writeWait,
		pongWait:          pongWait,
		pingPeriod:        pingPeriod,
		maxMessageSize:    maxMessageSize,
	}
}

// Run run hubs
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerHandler(client)
			h.clients[client.ID] = client
		case client := <-h.unregister:
			h.unRegisterHandler(client)
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
			}
		}
	}
}
