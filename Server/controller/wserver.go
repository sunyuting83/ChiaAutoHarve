package controller

import (
	"ChiaStart/Server/ws"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	Group := strings.Split(ClientName, "_")[0]
	Group = strings.Join([]string{Group, "_"}, "")
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if error != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	client := &ws.Client{ID: ClientName, Socket: conn, Send: make(chan []byte)}
	ws.Manager.GID = Group
	ws.Manager.Register <- client

	go client.Read()
	go client.Write()
}
