package chat

import (
	"log"
	"time"
)

// decodes an inbound message, parses message type, and sends it to correct Hub channel
func dispatch(c *Client, data []byte) {
	wsMessage, err := Decode[WebSocketMessage](data)
	if err != nil {
		log.Println(err)
		return
	}
	switch wsMessage.Type {
	case Chat:
		// messages from clients should only contain Text in payload
		chatMessageData, err := Decode[ChatMessageData](wsMessage.Payload)
		if err != nil {
			log.Println(err)
			return
		}
		// after reading in only the text from message, update the rest of message with client details
		updateChatMessageData(chatMessageData, c)
		// call dispatch to send to hub broadcast channel
		dispatchChatMessage(c.Hub, *chatMessageData)

	case UsernameUpdate:
		usernameUpdateData, err := Decode[UsernameUpdateData](wsMessage.Payload)
		if err != nil {
			log.Println(err)
			return
		}
		usernameUpdateData.Client = c
		c.Hub.usernameUpdate <- *usernameUpdateData
	default:
		log.Println("Unsupported WebSocket message type")
	}
}

// updates a chat message with the details of the sender client who is broadcasting it
func updateChatMessageData(chatMessageData *ChatMessageData, c *Client) {
	chatMessageData.SenderID = c.ID
	chatMessageData.SenderUsername = c.Username
	chatMessageData.Time = time.Now()
	chatMessageData.RoomID = c.RoomID
}

// notifications to a room when a new client joins or leaves the room
func dispatchNotification(hub *Hub, roomID string, text string) {
	chatMessageData := ChatMessageData{
		SenderID: NotificationSenderID,
		RoomID:   roomID,
		Text:     text,
		Time:     time.Now(),
	}
	dispatchChatMessage(hub, chatMessageData)
}

// enqueues a message a to hub broadcast channel to get sent to the room
func dispatchChatMessage(hub *Hub, chatMessageData ChatMessageData) {
	data, err := Encode(chatMessageData)
	if err != nil {
		log.Println(err)
		return
	}
	outboundWsMessage := WebSocketMessage{
		Type:    Chat,
		Payload: data,
	}
	data, err = Encode(outboundWsMessage)
	if err != nil {
		log.Println(err)
		return
	}
	hub.broadcast <- ChatMessage{chatMessageData.RoomID, data, chatMessageData.SenderUsername, chatMessageData.Text}
}
