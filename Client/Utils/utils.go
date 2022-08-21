package utils

import (
	tcping "ChiaStart/Client/Tcping"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WsServer *WsServer `yaml:"WsServer"`
	ChiaPath string    `yaml:"ChiaPath"`
}
type WsServer struct {
	Host       string `yaml:"Host"`
	Port       string `yaml:"Port"`
	Path       string `yaml:"Path"`
	WSType     string `yaml:"WSType"`
	SECRET_KEY string `yaml:"SECRET_KEY"`
}
type Message struct {
	Sender  string `json:"sender,omitempty"`
	Content string `json:"content,omitempty"`
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
	return confYaml, nil
}

// RunCommand run command
func RunCommand(OS, command string) (k string, err error) {
	var cmd *exec.Cmd
	if OS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("/bin/sh", "-c", command)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return "", err
	}

	if len(bytesErr) != 0 {
		return "", errors.New("0")

	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return string(bytes), nil
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
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

func GetUserInfo() (homedir string, err error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}

func MakeRun(ChiaRun, LinkPathStr, CurrentPath, OS, ip string) (command string) {
	var (
		ext string = "sh"
	)
	if OS == "windows" {
		ext = "bat"
	}
	shellFile := strings.Join([]string{"runit", ext}, ".")
	runPath := strings.Join([]string{CurrentPath, "script", shellFile}, LinkPathStr)
	command = strings.Join([]string{runPath, ChiaRun, ip}, " ")
	return
}

func RestartIt(shell, LinkPathStr, OS, ip string) {
	homedir, _ := GetUserInfo()
	configPath := strings.Join([]string{homedir, ".chia", "mainnet", "config", "config.yaml"}, LinkPathStr)
	ChiaConfig, _ := ioutil.ReadFile(configPath)

	if !strings.Contains(string(ChiaConfig), ip) {
		RunCommand(OS, shell)
	}
}
