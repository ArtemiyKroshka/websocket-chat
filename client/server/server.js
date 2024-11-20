import {WebSocketServer} from "ws";
const PORT = 8080;

const server = new WebSocket.Server({port: PORT});

const clients = new Set();

server.on("connection", (socket) => {
  clients.add(socket);
  console.log("New client connected");

  // Broadcast incoming messages to all clients
  socket.on("message", (message) => {
    console.log("Received:", message);
    for (const client of clients) {
      if (client !== socket && client.readyState === WebSocket.OPEN) {
        client.send(message);
      }
    }
  });

  // Remove the client on disconnect
  socket.on("close", () => {
    clients.delete(socket);
    console.log("Client disconnected");
  });
});

console.log(`WebSocket server is running on ws://localhost:${PORT}`);
