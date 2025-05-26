package handlers

import (
	"chatapp/internal/auth"
	"chatapp/internal/chat"
	"chatapp/internal/postgres"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// used to upgrade HTTP to Web sockets, websocket library manages the reassembly of websocket frames
var upgrader = websocket.Upgrader{ // websocket buffers use send/recv internally, just a pointer to userspace buffer
	WriteBufferSize: 1024, // I/O buffer sizes in user space, this is different from TCP buffer in kernel memory
	ReadBufferSize:  1024, // read and write buffers can only process 1 websocket frame at a time
}

// establish the websocket connection with client here
func ServeWsConn(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	// get userID, username from HTTP only cookie to populate name
	id, username, err := getClientInfo(w, r)
	if err != nil {
		log.Println(err)
		return // would've already wrote error to response, just return
	}
	// upgrade connection from HTTP to WebSocket protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	// create new client for the connection
	client := chat.NewClient(id, username, hub, conn)
	hub.RegisterClient(client) // push onto hub register channel
	log.Printf("New peer connecting from address: %s", r.RemoteAddr)
	go client.ReceiveWsMessage() // receive websocket frames on separate thread
	go client.SendWsMessage()    // send websocket frames on separate thread
}

// retrieve the users id and username from the first HTTP1.1 req that
// is starting the websocket handshake
func getClientInfo(w http.ResponseWriter, r *http.Request) (string, string, error) {
	claims, err := auth.VerifySessionToken(r)
	if err != nil {
		http.Error(w, "Failed to verify user", http.StatusBadRequest)
		return "", "", err
	}
	id := claims["id"].(string)
	username, err := postgres.GetUsernameById(id)
	if err != nil {
		http.Error(w, "Failed to fetch username from postgres", http.StatusBadRequest)
		return "", "", err
	}
	return id, username, err
}
