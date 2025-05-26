package chat

import (
	"encoding/json"
	"log"
	"time"
)

// read a websocket message, and dispatch to appropriate channel on hub
// only does the parsing and routing, hub will save to postgres
func dispatch(c *Client, data []byte) {
	var wsMessage WebSocketMessage
	if err := json.Unmarshal(data, &wsMessage); err != nil {
		log.Println("Error unmarshaling websocket message data: ", err)
		return
	}
	switch wsMessage.Type {
	case Chat:
		var chatMessageData ChatMessageData
		chatMessageData.SenderID = c.id
		chatMessageData.SenderUsername = c.username
		chatMessageData.Time = time.Now()
		if err := json.Unmarshal(wsMessage.Payload, &chatMessageData); err != nil {
			log.Println("Error unmarshaling chat message payload: ", err)
			return
		}
		payload, err := json.Marshal(chatMessageData)
		if err != nil {
			log.Println("Error marshaling chat message payload: ", err)
			return
		}
		outboundWsMessage := WebSocketMessage{
			Type:    Chat,
			Payload: payload,
		}
		data, err := json.Marshal(outboundWsMessage)
		if err != nil {
			log.Println("Error marshaling chat message data: ", err)
			return
		}
		c.hub.broadcast <- data

	case UsernameUpdate:
		var usernameUpdateData UsernameUpdateData
		usernameUpdateData.Client = c
		if err := json.Unmarshal(wsMessage.Payload, &usernameUpdateData); err != nil {
			log.Println("Error unmarshaling to UsernameUpdateData")
			return
		}
		//log.Printf("hub addr in Dispatcher: %p", c.hub)
		c.hub.usernameUpdate <- usernameUpdateData
	default:
		log.Println("Unsupported WebSocket message type")
	}
}
