package utils

import (
	tcping "ChiaStart/Tcping"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	IP       string `yaml:"IP"`
	ScanTime int    `yaml:"ScanTime"`
	ChiaPath string `yaml:"ChiaPath"`
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
func CheckConfig(OS, ConfigFile, LinkPathStr, ChiaExec string) (conf *Config, ChiaRun string, err error) {

	confYaml := new(Config)
	yamlFile, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return confYaml, "", errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		return confYaml, "", errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	if !IsDir(confYaml.ChiaPath) {
		return confYaml, "", errors.New("Chia运行目录设置\n10秒后程序自动关闭")
	}
	chiaRun := strings.Join([]string{confYaml.ChiaPath, ChiaExec}, LinkPathStr)
	if !FileExist(chiaRun) {
		return confYaml, "", errors.New("Chia运行目录设置\n10秒后程序自动关闭")
	}
	if len(confYaml.Host) <= 0 {
		if err != nil {
			return confYaml, "", errors.New("获取主机名失败\n10秒后程序自动关闭")
		}
	}
	tp := tcping.Tcping(10, confYaml.Port, confYaml.Host)
	if !tp {
		return confYaml, "", errors.New("端口未开放\n10秒后程序自动关闭")
	}
	if len(confYaml.Port) <= 0 {
		confYaml.Port = "8447"
		data, _ := yaml.Marshal(&confYaml)
		ioutil.WriteFile(ConfigFile, data, 0644)
	}
	if confYaml.ScanTime <= 0 {
		confYaml.ScanTime = 60
		data, _ := yaml.Marshal(&confYaml)
		ioutil.WriteFile(ConfigFile, data, 0644)
	}
	if len(confYaml.IP) <= 0 {
		ip, err := GetDomainIp(confYaml.Host)
		if err != nil {
			return confYaml, "", errors.New("获取IP失败，请检查网络\n10秒后程序自动关闭")
		}
		confYaml.IP = ip
		data, _ := yaml.Marshal(&confYaml)
		ioutil.WriteFile(ConfigFile, data, 0644)
	}
	return confYaml, chiaRun, nil
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

func GetDomainIp(host string) (ip string, err error) {
	addr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return "", errors.New("获取IP失败")
	}
	return addr.String(), nil
}
