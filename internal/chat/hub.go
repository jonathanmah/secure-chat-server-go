package chat

import (
	"encoding/json"
	"log"
)

// maintains active peer connections as clients and broadcasts messages
type Hub struct {
	// hashset of pointers to clients
	clients map[*Client]struct{}

	// inbound messages from peers, unbuffered to apply backpressure
	// don't want a single client taking up a buffered channel with spam messages
	broadcast chan []byte

	usernameUpdate chan UsernameUpdateData // used to update client name for broadcasting new user list

	// clients to register to Hub
	register chan *Client

	// clients to unregister from Hub
	unregister chan *Client
}

// create and return pointer to new Hub
func NewHub() *Hub {
	return &Hub{
		clients:        make(map[*Client]struct{}),
		broadcast:      make(chan []byte),
		usernameUpdate: make(chan UsernameUpdateData),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// manage clients and broadcasting messages
func (h *Hub) Run() {
	//log.Printf("hub addr in Run(): %p", h)
	for {
		select {
		// read in and add clients pending registration
		case client := <-h.register:
			h.clients[client] = struct{}{}
			h.broadcastActiveUserList()
		case client := <-h.unregister: // read in and remove clients pending unregistration
			close(client.send)
			delete(h.clients, client)
			h.broadcastActiveUserList()
		case message := <-h.broadcast: // read in a message from broadcast channel
			//#TODO save in postgres
			h.broadcastData(message)
		case clientUpdate := <-h.usernameUpdate:
			clientUpdate.Client.username = clientUpdate.Username
			h.broadcastActiveUserList()
		}
	}
}

// sends updated list of active usernames to all peers
// executes whenever a new client connects, disconnects, or changes name
func (h *Hub) broadcastActiveUserList() {
	var users []UserItem
	for client := range h.clients {
		users = append(users, UserItem{ID: client.id, Username: client.username})
	}
	userListMessage := UserListMessage{Users: users}
	payload, err := json.Marshal(userListMessage)
	if err != nil {
		log.Println("Error marshaling user list payload:", err)
		return
	}
	message := WebSocketMessage{
		Type:    "userlist",
		Payload: payload, // json.RawMessage type allows using already-marshaled payload
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshaling WebSocketMessage:", err)
		return
	}
	// broadcast active user list to all clients
	h.broadcastData(data)
}

// broadcast data to all active websocket connections
func (h *Hub) broadcastData(data []byte) {
	for client := range h.clients { // push data to all clients send buffered channels
		select {
		case client.send <- data:
		default: // default disconnect if client send buffered channel full and being slow
			close(client.send)
			delete(h.clients, client)
		}
	}
}
