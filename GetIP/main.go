package main

import (
	utils "ChiaStart/GetIP/Utils"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	OS := runtime.GOOS
	CurrentPath, _ := utils.GetCurrentPath()

	LinkPathStr := "/"
	ConfigFile := strings.Join([]string{CurrentPath, "config.yaml"}, LinkPathStr)

	confYaml, err := utils.CheckConfig(OS, ConfigFile, LinkPathStr)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}

	var ch chan int
	ticker := time.NewTicker(time.Second * time.Duration(confYaml.ScanTime))
	go func() {
		for range ticker.C {
			config, err := utils.GetConfigIP(OS, ConfigFile, LinkPathStr)
			if err != nil {
				fmt.Println(err)
			}

			IP, err := utils.GetIpAddr()
			if err != nil {
				fmt.Println(err)
			}
			if !config.SendStatus {
				wsSended := ConnentWs(confYaml, ConfigFile, IP)
				if wsSended {
					confYaml.SendStatus = true
				} else {
					confYaml.SendStatus = false
				}
				utils.SaveConfigFile(confYaml, ConfigFile)
			} else {
				if config.IP != IP {
					wsSended := ConnentWs(confYaml, ConfigFile, IP)
					if wsSended {
						confYaml.SendStatus = true
						confYaml.IP = IP
					} else {
						confYaml.SendStatus = false
						confYaml.IP = IP
					}
					utils.SaveConfigFile(confYaml, ConfigFile)
				}
			}

		}
		ch <- 1
	}()
	<-ch
}

// ConnentWs Connent Ws
func ConnentWs(confYaml *utils.Config, ConfigFile, IP string) bool {
	host := strings.Join([]string{confYaml.WsServer.Host, confYaml.WsServer.Port}, ":")
	wsurl := url.URL{Scheme: confYaml.WsServer.WSType, Host: host, Path: confYaml.WsServer.Path}
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(wsurl.String(), http.Header{"X-Api-Key": []string{confYaml.WsServer.SECRET_KEY}})

	if err != nil {
		fmt.Println("连接失败，10秒后重连")
		return false
	} else {
		fmt.Println("连接WServer成功")
		conn.WriteMessage(websocket.TextMessage, []byte(IP))
		return true
	}
}
