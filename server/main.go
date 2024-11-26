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
	ID       string
	conn     *websocket.Conn
	messages chan Message
}

type Message struct {
	SenderID string `json:"senderId"`
	Content  string `json:"content"`
}

var (
	clients  = make(map[string]*Client)
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

func (c *Client) send() {
	for msg := range c.messages {
		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("Error marshalling message: %v\n", err)
			continue
		}
		err = c.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Printf("Error writing message to client %s: %v\n", c.ID, err)
			return
		}
	}
}

func removeClient(clientID string) {
	clientMu.Lock()
	defer clientMu.Unlock()

	if client, exists := clients[clientID]; exists {
		close(client.messages)
		client.conn.Close()
		delete(clients, clientID)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading to WebSocket: %v\n", err)
		return
	}

	clientID := uuid.NewString()
	client := &Client{
		ID:       clientID,
		conn:     conn,
		messages: make(chan Message, 100),
	}

	clientMu.Lock()
	clients[clientID] = client
	clientMu.Unlock()

	go client.send()

	client.messages <- Message{
		SenderID: "server",
		Content:  clientID,
	}

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

func broadcast(clientID, messageData string) {
	message := Message{
		SenderID: clientID,
		Content:  messageData,
	}

	clientMu.Lock()
	for _, client := range clients {
		select {
		case client.messages <- message:
		default:
			fmt.Printf("Client %s is slow, dropping message\n", client.ID)
		}
	}
	clientMu.Unlock()
}

func main() {
	http.HandleFunc("/ws", handler)

	fmt.Printf("WebSocket server is running on ws://%s/ws\n", address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
