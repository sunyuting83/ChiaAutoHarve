package main

import (
	utils "ChiaStart/Client/Utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/websocket"
)

type Config struct {
	OS                 string
	LinkPathStr        string
	CurrentPath        string
	WSHost             string
	WSPort             string
	WSPath             string
	WSType             string
	SECRET_KEY         string
	RustDeskConfigFile string
}

type RustDeskConfig struct {
	RendezvousServer string `toml:"rendezvous_server"`
	NatType          int    `toml:"nat_type"`
	Serial           int    `toml:"serial"`

	Options *RustDeskOptions `toml:"options"`
}
type RustDeskOptions struct {
	Key                    string `toml:"key"`
	RendezvousServers      string `toml:"rendezvous-servers"`
	CustomRendezvousServer string `toml:"custom-rendezvous-server"`
}

var config *Config

// ConnentWs Connent Ws
func ConnentWs() {
	host := strings.Join([]string{config.WSHost, config.WSPort}, ":")
	if len(config.WSPort) == 0 {
		host = config.WSHost
	}
	wsurl := url.URL{Scheme: config.WSType, Host: host, Path: config.WSPath}
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(wsurl.String(), http.Header{"X-Api-Key": []string{config.SECRET_KEY}})

	if err != nil {
		fmt.Println("连接失败，10秒后重连")
		time.Sleep(time.Duration(10) * time.Second)
		ConnentWs()
	}
	conn.WriteMessage(websocket.TextMessage, []byte("getip"))

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("连接失败，10秒后重连")
			time.Sleep(time.Duration(10) * time.Second)
			ConnentWs()
		} else {
			var p *utils.Message
			if err := json.Unmarshal(message, &p); err != nil {
				fmt.Println("error")
			}
			IsIP := utils.CheckIP(p.Content)
			if IsIP {
				var conf *RustDeskConfig
				toml.DecodeFile(config.RustDeskConfigFile, &conf)
				if &p.Content != &conf.RendezvousServer {
					buf := new(bytes.Buffer)
					conf.RendezvousServer = p.Content
					conf.Options.CustomRendezvousServer = p.Content
					toml.NewEncoder(buf).Encode(&conf)
					ioutil.WriteFile(config.RustDeskConfigFile, buf.Bytes(), 0644)
				}
			}
		}
	}
}

func main() {
	OS := runtime.GOOS
	CurrentPath, _ := utils.GetCurrentPath()
	homeDir, _ := utils.GetUserInfo()
	LinkPathStr := "/"
	rustDeskConfigFile := strings.Join([]string{homeDir, ".config", "rustdesk", "RustDesk2.toml"}, LinkPathStr)
	if OS == "windows" {
		LinkPathStr = "\\"
		rustDeskConfigFile = strings.Join([]string{homeDir, "AppData", "Roaming", "RustDesk", "config", "RustDesk2.toml"}, LinkPathStr)
	}
	ConfigFile := strings.Join([]string{CurrentPath, "config.yaml"}, LinkPathStr)

	confYaml, err := utils.CheckConfig(OS, ConfigFile, LinkPathStr)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}

	config = &Config{
		OS:                 OS,
		LinkPathStr:        LinkPathStr,
		CurrentPath:        CurrentPath,
		WSHost:             confYaml.WsServer.Host,
		WSPort:             confYaml.WsServer.Port,
		WSPath:             confYaml.WsServer.Path,
		WSType:             confYaml.WsServer.WSType,
		SECRET_KEY:         confYaml.WsServer.SECRET_KEY,
		RustDeskConfigFile: rustDeskConfigFile,
	}
	ConnentWs()
}
