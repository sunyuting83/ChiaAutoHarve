package main

import (
	"ChiaStart/Server/controller"
	"ChiaStart/Server/utils"
	"ChiaStart/Server/ws"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	OS := runtime.GOOS
	CurrentPath, _ := utils.GetCurrentPath()

	confYaml, err := utils.CheckConfig(OS, CurrentPath)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	port := strconv.Itoa(confYaml.Port)
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(utils.CORSMiddleware())
	{
		app.GET("/server", controller.WsServer)
	}

	go ws.Manager.Start()
	app.Run(strings.Join([]string{":", port}, ""))
}
