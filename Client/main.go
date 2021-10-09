package main

import (
	utils "ChiaStart/Client/Utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// ConnentWs Connent Ws
func ConnentWs(confYaml *utils.Config, OS, ChiaRun, LinkPathStr, CurrentPath string) {
	host := strings.Join([]string{confYaml.WsServer.Host, confYaml.WsServer.Port}, ":")
	wsurl := url.URL{Scheme: confYaml.WsServer.WSType, Host: host, Path: confYaml.WsServer.Path}
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(wsurl.String(), http.Header{"X-Api-Key": []string{confYaml.WsServer.SECRET_KEY}})

	if err != nil {
		fmt.Println("连接失败，10秒后重连")
		time.Sleep(time.Duration(10) * time.Second)
		ConnentWs(confYaml, OS, ChiaRun, LinkPathStr, CurrentPath)
	}
	conn.WriteMessage(websocket.TextMessage, []byte("getip"))

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("连接失败，10秒后重连")
			time.Sleep(time.Duration(10) * time.Second)
			ConnentWs(confYaml, OS, ChiaRun, LinkPathStr, CurrentPath)
		} else {
			var p *utils.Message
			if err := json.Unmarshal(message, &p); err != nil {
				fmt.Println("error")
			}
			IsIP := utils.CheckIP(p.Content)
			if IsIP {
				shell := utils.MakeRun(ChiaRun, LinkPathStr, CurrentPath, OS, p.Content)
				go utils.RestartIt(shell, LinkPathStr, OS, p.Content)
			}
		}
	}
}

func main() {
	OS := runtime.GOOS
	CurrentPath, _ := utils.GetCurrentPath()

	LinkPathStr := "/"
	ChiaExec := "chia"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	ConfigFile := strings.Join([]string{CurrentPath, "config.yaml"}, LinkPathStr)

	confYaml, ChiaRun, err := utils.CheckConfig(OS, ConfigFile, LinkPathStr, ChiaExec)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	ConnentWs(confYaml, OS, ChiaRun, LinkPathStr, CurrentPath)
}
