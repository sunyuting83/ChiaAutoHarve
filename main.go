package main

import (
	utils "ChiaStart/Utils"
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
			checkIP(confYaml, ConfigFile, ChiaRun, OS)
		}
		ch <- 1
	}()
	<-ch
}

func checkIP(confYaml *utils.Config, ConfigFile, ChiaRun, OS string) {
	ip, err := utils.GetDomainIp(confYaml.Host)
	if err != nil {
		fmt.Println(err)
	}
	if ip != confYaml.IP {
		confYaml.IP = ip
		data, _ := yaml.Marshal(&confYaml)
		ioutil.WriteFile(ConfigFile, data, 0644)
		command := strings.Join([]string{ChiaRun, "start", "harvester", "-r"}, " ")
		go utils.RunCommand(OS, command)
	}
}
