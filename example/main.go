package main

import (
	"fmt"
	"ginLogger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)


func main()  {
	if err := ginLogger.Init("./config.json"); err  != nil{
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	if err := ginLogger.InitLogger(ginLogger.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	gin.SetMode(ginLogger.Conf.Mode)
	r := gin.Default()
	r.Use(ginLogger.GinLogger(), ginLogger.GinRecovery(true))
	r.GET("/hello", func(c *gin.Context) {
		// 假设你有一些数据需要记录到日志中
		var (
			name = "q1mi"
			age  = 18
		)
		// 记录日志并使用zap.Xxx(key, val)记录相关字段
		zap.L().Debug("this is hello func", zap.String("user", name), zap.Int("age", age))

		c.String(http.StatusOK, "hello averagejoe.wang\n")
	})

	addr := fmt.Sprintf(":%v", ginLogger.Conf.Port)
	r.Run(addr)
}
