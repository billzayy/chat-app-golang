package pkg

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	conn   *websocket.Conn
	userID primitive.ObjectID
}

type MessageType struct {
	UserID  primitive.ObjectID `json:"userId"`
	Message string             `json:"message"`
}

var clients = make(map[primitive.ObjectID]*Client)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // upgrades the HTTP server connection to the WebSocket protocol.
	if err != nil {
		log.Println(err)
		return
	}

	inputId := r.URL.Query().Get("userId")

	userID, err := primitive.ObjectIDFromHex(inputId)

	if err != nil {
		panic(err)
	}

	client := &Client{conn: conn, userID: userID}
	clients[userID] = client

	connected, err := json.Marshal(clients)

	if err != nil {
		panic(err)
	}

	SendMessageToClient(client.userID, connected)

	go handleMessages(client)
}

func handleMessages(client *Client) {
	defer func() {
		// delete(clients, client.userID)
		client.conn.Close()
	}()

	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket disconnected")
			break
		}
		// Handle incoming messages here
	}
}

func SendMessageToClient(userID primitive.ObjectID, message []byte) {
	client, ok := clients[userID]
	if !ok {
		log.Printf("No client found with userID %s", userID)
		return
	}

	err := client.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Error sending message to client %s: %v", userID, err)
	}

	// fmt.Printf("Messages is sent to %s \n", userID)
}
