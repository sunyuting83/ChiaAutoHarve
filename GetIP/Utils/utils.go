package utils

import (
	tcping "ChiaStart/Client/Tcping"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	IP         string    `yaml:"IP"`
	ScanTime   int       `yaml:"ScanTime"`
	SendStatus bool      `yaml:"SendStatus"`
	WsServer   *WsServer `yaml:"WsServer"`
}
type WsServer struct {
	Host       string `yaml:"Host"`
	Port       string `yaml:"Port"`
	Path       string `yaml:"Path"`
	SECRET_KEY string `yaml:"SECRET_KEY"`
}

var confYaml *Config

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
func CheckConfig(OS, ConfigFile, LinkPathStr string) (conf *Config, err error) {
	yamlFile, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	if len(confYaml.WsServer.Host) <= 0 {
		if err != nil {
			return confYaml, errors.New("获取主机名失败\n10秒后程序自动关闭")
		}
	}
	tp := tcping.Tcping(10, confYaml.WsServer.Port, confYaml.WsServer.Host)
	if !tp {
		return confYaml, errors.New("WS服务器链接失败\n10秒后程序自动关闭")
	}
	if confYaml.ScanTime <= 0 {
		confYaml.ScanTime = 60
		data, _ := yaml.Marshal(&confYaml)
		ioutil.WriteFile(ConfigFile, data, 0644)
	}
	if len(confYaml.IP) <= 0 {
		ip, err := GetIpAddr()
		if err != nil {
			return confYaml, errors.New("获取IP失败，请检查网络\n10秒后程序自动关闭")
		}
		confYaml.IP = ip
		data, _ := yaml.Marshal(&confYaml)
		ioutil.WriteFile(ConfigFile, data, 0644)
	}
	return confYaml, nil
}

func GetConfigIP(OS, ConfigFile, LinkPathStr string) (conf *Config, err error) {
	yamlFile, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	return confYaml, nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func GetIpAddr() (i string, err error) {
	ip, err := getData("https://www.taobao.com/help/getip.php", "GET", []byte(""), "")
	if err != nil {
		return "", errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	ips := string(ip)
	length := len(ips)
	start := strings.Index(ips, `ip:"`)
	a := ips[start+4 : length]
	end := strings.Index(a, `"}`)
	i = a[0:end]
	return
}

// getData get data
func getData(url string, types string, data []byte, password string) (s []byte, err error) {
	client := &http.Client{}
	reqest, err := http.NewRequest(types, url, bytes.NewBuffer(data))

	if len(password) > 0 {
		reqest.Header.Add("Authorization", strings.Join([]string{"Bearer", password}, " "))
	}

	if err != nil {
		return []byte(""), err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return []byte(""), err
	}
	defer response.Body.Close()
	d, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte(""), err
	}
	return d, nil
}

func SaveConfigFile(confYaml *Config, ConfigFile string) {
	data, _ := yaml.Marshal(&confYaml)
	ioutil.WriteFile(ConfigFile, data, 0644)
}
