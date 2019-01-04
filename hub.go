package h2oio

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
}

// NewHub buuild new hub
func NewHub(registerHandler func(*Client), unRegisterHandler func(*Client), messageHandler func(*Client, []byte)) *Hub {
	return &Hub{
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		clients:           make(map[string]*Client),
		registerHandler:   registerHandler,
		messageHandler:    messageHandler,
		unRegisterHandler: unRegisterHandler,
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
