package main

import (
	"encoding/json"
	"fmt"
	"github.com/AverageJoeWang/ginLogger"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

var Conf = new(Config)

// Config
type Config struct {
	Mode                 string `json:"mode"`
	Port                 int    `json:"port"`
	*ginLogger.LogConfig `json:"log"`
}

// Init
func Init(filePath string) error {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, Conf)
}

func main() {
	if err := Init("./config.json"); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	if err := ginLogger.InitLogger(Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	gin.SetMode(Conf.Mode)
	r := gin.Default()
	r.Use(ginLogger.GinLogger(), ginLogger.GinRecovery(true))
	r.GET("/hello", func(c *gin.Context) {
		ginLogger.Debug(c, "this is hello")
		c.String(http.StatusOK, "hello AverageJoeWang\n")
	})
	r.Run(":8180")
}
