package chat

import (
	"encoding/json"
	"time"
)

// -------------------------- WEB SOCKET MODELS -------------------------------------------

// Instead of interpreting HTTP Methods and URL paths, we create our own custom protocol
// by defining different types of websocket messages and payloads with JSON

type MessageType string

const (
	Chat           MessageType = "chat"
	UsernameUpdate MessageType = "username_update"
	UserList       MessageType = "userlist"
)

// if sending a WebSocket message to peer, will be encoded into a JSON byte slice at the transport layer
// if reading a WebSocket message from peer, will be decoded into one of the message structs for the hub
type WebSocketMessage struct {
	Type    MessageType     `json:"type"`    // chat (bidirectional) | userlist (outbound) | updatename
	Payload json.RawMessage `json:"payload"` // defer decoding of nested json until after interpret what ws message type
}

// Dispatched to Hub
type ChatMessageData struct { // omitempty will leave field out of json, otherwise zero value if not defined
	MessageID        string    `json:"message_id,omitempty"`
	SenderID         string    `json:"sender_id,omitempty"`
	SenderUsername   string    `json:"sender_username,omitempty"`
	ReceiverID       string    `json:"receiver_id,omitempty"`
	ReceiverUsername string    `json:"receiver_username,omitempty"`
	RoomID           string    `json:"room_id,omitempty"`
	Text             string    `json:"text"`
	Time             time.Time `json:"time,omitempty"`
}

// Dispatched to Hub
type UsernameUpdateData struct {
	Client   *Client
	Username string `json:"username"`
}

// Outgoing websocket payload to broadcast active users
type UserListMessage struct {
	Users []UserItem `json:"users"`
}
type UserItem struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
