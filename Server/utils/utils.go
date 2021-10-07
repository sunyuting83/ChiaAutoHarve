package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port       string `yaml:"port"`
	WsPath     string `yaml:"WsPath"`
	SECRET_KEY string `yaml:"SECRET_KEY"`
}

// GetCurrentPath Get Current Path
func GetCurrentPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	return dir, nil
}

// CheckConfig check config
func CheckConfig(OS, CurrentPath string) (conf *Config, err error) {
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	ConfigFile := strings.Join([]string{CurrentPath, "config.yaml"}, LinkPathStr)

	var confYaml *Config
	yamlFile, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	if len(confYaml.Port) <= 0 {
		confYaml.Port = "13001"
		data, _ := yaml.Marshal(&confYaml)
		ioutil.WriteFile(ConfigFile, data, 0644)
	}
	if len(confYaml.WsPath) <= 0 {
		confYaml.WsPath = "server"
		data, _ := yaml.Marshal(&confYaml)
		ioutil.WriteFile(ConfigFile, data, 0644)
	}
	return confYaml, nil
}

// CORSMiddleware cors middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func CheckIP(ip string) bool {
	if len(ip) < 7 || len(ip) > 15 {
		return false
	}

	ipArray := strings.Split(ip, ".")
	if len(ipArray) != 4 {
		return false
	}
	for _, v := range ipArray {
		number, err := strconv.Atoi(v)
		if err != nil {
			return false
		}
		if number < 0 || number > 255 {
			return false
		}
	}
	return true
}

func ReadIPFile() (string, error) {
	OS := runtime.GOOS
	CurrentPath, _ := GetCurrentPath()
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	IPFile := strings.Join([]string{CurrentPath, "ipdb"}, LinkPathStr)
	IP, err := ioutil.ReadFile(IPFile)
	if err != nil {
		return "0", err
	}
	return string(IP), err
}

func SaveIPFile(IP []byte) {
	OS := runtime.GOOS
	CurrentPath, _ := GetCurrentPath()
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	IPFile := strings.Join([]string{CurrentPath, "ipdb"}, LinkPathStr)
	ioutil.WriteFile(IPFile, IP, 0644)
}
