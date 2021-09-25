package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port       int    `yaml:"port"`
	SECRET_KEY string `yaml:"SECRET_KEY"`
}

// Verification power;
func Verification(c *gin.Context) {
	Token, ok := c.GetQuery("token")
	if ok {
		if CheckToken(Token) {
			c.Next()
		} else {
			c.AbortWithStatus(403)
		}
	} else {
		c.AbortWithStatus(403)
	}
}

// CheckToken is a check token function
func CheckToken(a string) bool {
	fmt.Println(a)
	return true
}

// checkHas check token has Bearer
func checkHas(a string, b string) bool {
	return strings.ContainsAny(a, b)
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

// SetConfigMiddleWare set config
func SetConfigMiddleWare(SECRET_KEY string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("secret_key", SECRET_KEY)
		c.Writer.Status()
	}
}
