package controller

import (
	"ChiaStart/Server/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func WsServer(c *gin.Context) {
	var ClientName string
	ClientName = c.GetHeader("X-Api-Key")
	if len(ClientName) <= 0 {
		Token, ok := c.GetQuery("token")
		if ok {
			ClientName = Token
		} else {
			http.NotFound(c.Writer, c.Request)
			return
		}
	}
	SECRET_KEY, _ := c.Get("token")
	if ClientName != SECRET_KEY.(string) {
		http.NotFound(c.Writer, c.Request)
	}
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if error != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	client := &ws.Client{ID: uuid.NewV4().String(), Socket: conn, Send: make(chan []byte)}
	ws.Manager.Register <- client
	go client.Read()
	go client.Write()
}
