package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// client is a middleman between websocket connection and hub
type Client struct {
	ID       string // uuid of the peer
	Username string // username of peer
	RoomID   string // room id peer is subscribed to

	Hub *Hub // the hub managing this client
	// the websocket connection.
	Conn *websocket.Conn
	// buffered channel of outbound messages
	Send chan []byte
}

const (
	// time allowed to write a message to peer
	writeWait = 10 * time.Second
	// ping and pong used to detect broken/dead connections to close
	// server waits up to 60 seconds to receive a pong after sending a ping
	pongWait = 60 * time.Second
	// send pings to peer with this period, make it 90% of pong wait to give second ping a chance incase
	// first ping was lost or dead
	pingPeriod = (pongWait * 9) / 10
	// maximum total message size allowed to receive from peer
	maxMessageSize = 512
)

func NewClient(id string, username string, roomID string, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		ID:       id,
		Username: username,
		RoomID:   roomID,
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}
}

// transfers messages from websocket connection receive buffer to the hub broadcast channel
func (c *Client) ReceiveWsMessage() {
	defer func() { // unregister from hub and close connection when client no longer reading
		c.Hub.UnregisterClient(c)
		c.Conn.Close()
		log.Printf("Closed connection with %s in Room %s", c.Username, c.RoomID)
	}()
	c.Conn.SetReadLimit(maxMessageSize)              // max message size for all frames combined
	c.Conn.SetReadDeadline(time.Now().Add(pongWait)) // ReadMessage() will error if called after deadline
	// only a pong message can reset the pong timeout
	c.Conn.SetPongHandler(func(string) error { // pong handler is a callback function that gets called when pong frame is received
		c.Conn.SetReadDeadline(time.Now().Add(pongWait)) // resets read deadline for next pong message
		return nil
	})
	for {
		// call SetReadDeadline explicitly, which is done by the Pong Handler to manage heartbeat
		// do not enforce readmessage with a deadline to read normal messages, only for pong, because some peers might only listen
		_, message, err := c.Conn.ReadMessage() // read all frames of a message into []byte
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}
		// push message into hub broadcast channel buffer
		dispatch(c, message)
	}
}

// sends message []byte from client send buffer to websocket connection
func (c *Client) SendWsMessage() {
	ticker := time.NewTicker(pingPeriod) // create ticker to send pings as a heartbeat for checking websocket connection
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			// write deadline gets reset before each write. it may be passed deadline if inactive but always resets here
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait)) // must set write deadline, otherwise none responsive client may hang
			if !ok {
				// hub closed this clients send channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{}) // this triggers the frontend JS socket.OnClose()
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage) // create a websocket writer
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C: // signal sent to ticker.C every ping period
			// reset write deadline every ping
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
