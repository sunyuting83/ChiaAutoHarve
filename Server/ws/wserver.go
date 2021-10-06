package ws

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

// Start is to start a ws server
func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.Register:
			manager.Clients[conn] = true
		case conn := <-manager.Unregister:
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
			}
		case message := <-manager.Broadcast:
			for conn := range manager.Clients {
				var p *Message
				if err := json.Unmarshal(message, &p); err != nil {
					continue
				}
				if conn.ID == p.Sender {
					continue
				}
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}

// Send is to send ws message to ws client
func (manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		m := string(message)
		jsonMessage, _ := json.Marshal(&Message{Sender: c.ID, Content: m})
		Manager.Broadcast <- jsonMessage
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Mux.Lock()
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				c.Mux.Unlock()
				return
			}
			c.Mux.Lock()
			c.Socket.WriteMessage(websocket.TextMessage, message)
			c.Mux.Unlock()
		}
	}
}
