package ws

import (
	"encoding/json"
	"strings"

	"github.com/gorilla/websocket"
)

// Start is to start a ws server
func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.Register:
			manager.Clients[conn] = true
			if strings.HasPrefix(conn.ID, "manage_") {
				c := "sdfsdf"
				jsonMessage, _ := json.Marshal(&Message{Sender: conn.ID, Content: string(c)})
				manager.Send(jsonMessage, nil, "2m")
			}
		case conn := <-manager.Unregister:
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
				if strings.HasPrefix(conn.ID, "client_") {
					jsonMessage, _ := json.Marshal(&Message{Sender: conn.ID, Content: "offline"})
					manager.Send(jsonMessage, conn, "2m")
				}
			}
		case message := <-manager.Broadcast:
			for conn := range manager.Clients {
				if !strings.HasPrefix(conn.ID, manager.GID) {
					continue
				}
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.Clients, conn)
					conn.Mux.Lock()
				}
			}
		}
	}
}

// Send is to send ws message to ws client
func (manager *ClientManager) Send(message []byte, ignore *Client, SendTo string) {
	if SendTo == "m2c" {
		var (
			messages *Message
			command  *Command
		)
		json.Unmarshal([]byte(message), &messages)
		json.Unmarshal([]byte(messages.Content), &command)
		if len(command.ClientList) > 0 {
			for _, ig := range command.ClientList {
				for conn := range manager.Clients {
					if conn.ID == ig {
						conn.Send <- message
					}
				}
			}
		}
	} else {
		for conn := range manager.Clients {
			if SendTo == "2m" {
				if strings.HasPrefix(conn.ID, "manage_") {
					conn.Send <- message
				}
			} else {
				if conn != ignore {
					conn.Send <- message
				}
			}
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
		var (
			command *Command
		)
		if err := json.Unmarshal([]byte(message), &command); err != nil {
			jsonMessage, _ := json.Marshal(&Message{Sender: c.ID, Content: `{"status":1,"message":"error"}`})
			Manager.Send(jsonMessage, c, "2m")
		} else {
			jsonMessage, _ := json.Marshal(&Message{Sender: c.ID, Content: m})
			// fmt.Println(command)
			if command.SendTo == "m2c" {
				// fmt.Println("here")
				Manager.Send(jsonMessage, c, "m2c")
			} else {
				SetGid(c.ID, &Manager)
				Manager.Broadcast <- jsonMessage
			}
		}

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
func SetGid(ClientID string, manager *ClientManager) {
	if strings.HasPrefix(ClientID, "client_") {
		manager.GID = "manage_"
	} else {
		manager.GID = "client_"
	}
}
