package main

import (
	utils "ChiaStart/Utils"
	"fmt"
	"os"
	"runtime"
	"time"
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

	err, ip := utils.GetDomainIp(confYaml.Host)
	if err != nil {
		fmt.Println(err)
	}
	if ip == confYaml.IP {
		fmt.Println(ip)
	}
}
