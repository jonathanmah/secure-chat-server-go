package chat

import (
	"log"
	"time"
)

// read a websocket message, and dispatch to appropriate channel on hub
// only does the parsing and routing, hub will save to postgres
func dispatch(c *Client, data []byte) {
	wsMessage, err := Decode[WebSocketMessage](data)
	if err != nil {
		log.Println(err)
		return
	}
	switch wsMessage.Type {
	case Chat:
		chatMessageData, err := Decode[ChatMessageData](wsMessage.Payload)
		if err != nil {
			log.Println(err)
			return
		}
		updateChatMessageData(chatMessageData, c)
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

func updateChatMessageData(chatMessageData *ChatMessageData, c *Client) {
	chatMessageData.SenderID = c.ID
	chatMessageData.SenderUsername = c.Username
	chatMessageData.Time = time.Now()
	chatMessageData.RoomID = c.RoomID
}

// enqueues a message received from a listening websocket connection to hub broadcast channel
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

// notifications to a room if new client joins or leaves the room
func dispatchNotification(hub *Hub, roomID string, text string) {
	chatMessageData := ChatMessageData{
		SenderUsername: "Hub",
		SenderID:       "notification",
		RoomID:         roomID,
		Text:           text,
		Time:           time.Now(),
	}
	dispatchChatMessage(hub, chatMessageData)
}
