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
	Chat           MessageType = "chat"            // (bidirectional) - receives messages from clients and broadcasts them
	UsernameUpdate MessageType = "username_update" // (inbound) - updates the clients username and triggers a new userlist broadcast
	UserList       MessageType = "userlist"        // (outbound) - updates active user lists with current connected clients
)

const (
	NotificationSenderID string = "notification" // SenderID used when Hub is sending notifications to a room
)

// if sending an outbound WebSocket message to peer, will be encoded into a JSON byte slice at the transport layer
// if reading an inbound WebSocket message from peer, will be decoded into one of the structs below
type WebSocketMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Message Type: Chat
// Direction: Bidirectional
// Purpose: Inbound chat messages are added to Hub broadcast channel, then are sent outbound. Also used for chat notifications.
type ChatMessageData struct {
	MessageID        string    `json:"message_id,omitempty"`
	SenderID         string    `json:"sender_id,omitempty"`
	SenderUsername   string    `json:"sender_username,omitempty"`
	ReceiverID       string    `json:"receiver_id,omitempty"`
	ReceiverUsername string    `json:"receiver_username,omitempty"`
	RoomID           string    `json:"room_id,omitempty"`
	Text             string    `json:"text"`
	Time             time.Time `json:"time,omitempty"`
}

// Message Type: UsernameUpdate
// Direction: Inbound
// Purpose:
type UsernameUpdateData struct {
	Client   *Client
	Username string `json:"username"`
}

// Message Type: UserList
// Direction: Outbound
// Purpose: The payload inside a WebSocketMessage to update currently active users
type UserListMessage struct {
	Users []UserItem `json:"users"`
}

// A single entry to represent a connected client
type UserItem struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
