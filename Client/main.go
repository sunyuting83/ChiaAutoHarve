package main

import (
	utils "ChiaStart/Client/Utils"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func main() {
	OS := runtime.GOOS
	CurrentPath, _ := utils.GetCurrentPath()

	LinkPathStr := "/"
	ChiaExec := "chia"
	if OS == "windows" {
		LinkPathStr = "\\"
		ChiaExec = "chia.exe"
	}
	ConfigFile := strings.Join([]string{CurrentPath, "config.yaml"}, LinkPathStr)

	confYaml, ChiaRun, err := utils.CheckConfig(OS, ConfigFile, LinkPathStr, ChiaExec)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}

	var ch chan int
	ticker := time.NewTicker(time.Second * time.Duration(confYaml.ScanTime))
	go func() {
		for range ticker.C {
			checkIP(confYaml.Host, ConfigFile, ChiaRun, OS, LinkPathStr)
		}
		ch <- 1
	}()
	<-ch
}

func checkIP(host, ConfigFile, ChiaRun, OS, LinkPathStr string) {
	ip, err := utils.GetDomainIp(host)
	if err != nil {
		fmt.Println(err)
	}
	confIP, err := utils.GetConfigIP(OS, ConfigFile, LinkPathStr)
	if err != nil {
		fmt.Println(err)
	}
	if ip != confIP.IP {
		confIP.IP = ip
		data, _ := yaml.Marshal(&confIP)
		ioutil.WriteFile(ConfigFile, data, 0644)
		command := strings.Join([]string{ChiaRun, "start", "harvester", "-r"}, " ")
		utils.RunCommand(OS, command)
	}
}
