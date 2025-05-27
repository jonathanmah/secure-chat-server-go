package chat

import (
	"encoding/json"
	"fmt"
	"log"
)

// maintains active peer connections as clients and broadcasts messages
type Hub struct {
	// hashmap of Key:RoomID, Value: hashset of pointers to clients
	rooms map[string]map[*Client]struct{}

	// inbound messages from peers, unbuffered for backpressure
	// don't want a single client taking up a buffered channel with spam messages
	broadcast chan ChatMessage

	usernameUpdate chan UsernameUpdateData // used to update client name for broadcasting new user list

	// clients to register to Hub
	register chan *Client

	// clients to unregister from Hub
	unregister chan *Client
}

type ChatMessage struct {
	RoomID         string // room ID to broadcast message
	Data           []byte // encoded WebSocket data including payload
	SenderUsername string // using for logging
	MessageText    string // using for logging
}

// create and return pointer to new Hub
func NewHub() *Hub {
	return &Hub{
		rooms:          make(map[string]map[*Client]struct{}),
		broadcast:      make(chan ChatMessage),
		usernameUpdate: make(chan UsernameUpdateData),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
	}
}

func (h *Hub) RegisterClient(c *Client) {
	h.register <- c
}

func (h *Hub) UnregisterClient(c *Client) {
	h.unregister <- c
}

// manage clients and broadcasting messages
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegisterClient(client)
		case client := <-h.unregister:
			h.handleUnregisterClient(client)
		case chatMessage := <-h.broadcast:
			h.handleBroadcastChatMessage(chatMessage)
		case usernameUpdate := <-h.usernameUpdate:
			h.handleUsernameUpdate(usernameUpdate)
		}
	}
}

// creates a client for a newly connected peer and registers to a room
func (h *Hub) handleRegisterClient(c *Client) {
	if h.rooms[c.RoomID] == nil { // if room doesn't exist yet create a new one
		h.rooms[c.RoomID] = make(map[*Client]struct{})
	}
	h.rooms[c.RoomID][c] = struct{}{}
	h.broadcastActiveUserList(c.RoomID)

	msg := fmt.Sprintf("%s has joined Room %s ", c.Username, c.RoomID)
	go dispatchNotification(h, c.RoomID, msg)

	log.Printf("%s has joined Room %s ", c.Username, c.RoomID) // #TODO dispatch joined room message custom
}

// removes a client from their current room
func (h *Hub) handleUnregisterClient(c *Client) {
	close(c.Send)
	delete(h.rooms[c.RoomID], c) // remove client from room

	msg := fmt.Sprintf("%s has left Room %s ", c.Username, c.RoomID)
	go dispatchNotification(h, c.RoomID, msg)

	if len(h.rooms[c.RoomID]) == 0 { // delete room if it's empty
		delete(h.rooms, c.RoomID)
		log.Printf("Deleted empty Room %s.", c.RoomID)
	} else {
		h.broadcastActiveUserList(c.RoomID) // broadcast to the room current active users
	}
}

// handler for broadcasting chat messages
func (h *Hub) handleBroadcastChatMessage(message ChatMessage) {
	log.Printf("(Room %s) %s: %s", message.RoomID, message.SenderUsername, message.MessageText)
	h.broadcastData(message.RoomID, message.Data)
}

func (h *Hub) handleUsernameUpdate(update UsernameUpdateData) {
	update.Client.Username = update.Username
	h.broadcastActiveUserList(update.Client.RoomID)
}

// sends updated list of active users in a room
// executes whenever a new client connects, disconnects, or changes name
func (h *Hub) broadcastActiveUserList(RoomID string) {
	room := h.rooms[RoomID]
	if room == nil {
		log.Println("Tried to broadcast to empty room")
		return
	}
	var users []UserItem
	for client := range room {
		users = append(users, UserItem{ID: client.ID, Username: client.Username})
	}
	payload, err := json.Marshal(UserListMessage{Users: users})
	if err != nil {
		log.Println(err)
		return
	}
	data, err := json.Marshal(WebSocketMessage{Type: "userlist", Payload: payload})
	if err != nil {
		log.Println(err)
		return
	}
	// broadcast active user list to all clients
	h.broadcastData(RoomID, data)
}

// broadcasts an encoded WebSocketMessage to all clients in the room
func (h *Hub) broadcastData(RoomID string, data []byte) {
	room := h.rooms[RoomID]
	if room == nil {
		log.Println("Tried to broadcast to empty room")
		return
	}
	for client := range room { // push data to all clients send buffered channels
		select {
		case client.Send <- data:
		default: // default disconnect if client send buffered channel full and being slow
			close(client.Send)
			delete(room, client)
		}
	}
}
