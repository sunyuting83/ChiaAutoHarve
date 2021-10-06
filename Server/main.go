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

// SetConfigMiddleWare set config
func SetConfigMiddleWare(config *utils.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("token", config.SECRET_KEY)
		c.Writer.Status()
	}
}

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
	app.Use(SetConfigMiddleWare(confYaml))
	{
		app.GET("/server", controller.WsServer)
	}

	go ws.Manager.Start()
	app.Run(strings.Join([]string{":", port}, ""))
}
