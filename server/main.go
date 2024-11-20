package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type IClient interface {
	send(message Message)
}

type Client struct {
	ID   string
	conn *websocket.Conn
	mu   sync.Mutex
}

type Message struct {
	SenderID string `json:"senderId"`
	Content  string `json:"content"`
}

var (
	clients  = make([]*Client, 0)
	clientMu sync.Mutex
	address  = "localhost:8080"
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (c *Client) send(message Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error marshalling message: %v\n", err)
		return
	}

	err = c.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Printf("Error writing message to client %s: %v\n", c.ID, err)
	}
}

func removeClient(clientID string) {
	clientMu.Lock()
	defer clientMu.Unlock()

	for idx, el := range clients {
		if el.ID == clientID {
			clients = append(clients[:idx], clients[:idx+1]...)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading to WebSocket: %v\n", err)
		return
	}
	defer conn.Close()

	clientID := uuid.NewString()
	client := &Client{
		ID:   clientID,
		conn: conn,
	}

	clientMu.Lock()
	clients = append(clients, client)
	clientMu.Unlock()

	client.send(Message{
		SenderID: "server",
		Content:  clientID,
	})

	fmt.Printf("Client connected: %s\n", clientID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			removeClient(clientID)
			break
		}

		broadcast(clientID, string(message))
	}
}

func broadcast(clientId, messageData string) {
	clientMu.Lock()
	defer clientMu.Unlock()
	message := Message{
		SenderID: clientId,
		Content:  messageData,
	}

	for _, client := range clients {
		client.send(message)
	}
}

func main() {
	http.HandleFunc("/ws", handler)

	fmt.Printf("WebSocket server is running on ws://%s/ws\n", address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
