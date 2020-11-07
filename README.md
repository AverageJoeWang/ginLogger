

### 日志组件功能

+ 1.允许自定义日志等级,小于该等级的日志都不打印
+ 2.允许自定义日志字段
+ 3.封装对外暴露的接口，越简单越好
+ 4.允许管理日志文件路径、名字、大小、保存时间、打印格式、请求行
+ 5.允许自定义带的参数，类型，固定日志参数，生成独一无二的请求id，获取网关id等内容


### 例子

```go

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
			name = "oliver"
			age  = 27
		)
		// 记录日志并使用zap.Xxx(key, val)记录相关字段
		zap.L().Debug("this is hello", zap.String("user", name), zap.Int("age", age))

		c.String(http.StatusOK, "hello averagejoe.wang\n")
	})

	addr := fmt.Sprintf(":%v", ginLogger.Conf.Port)
	r.Run(addr)
}
```

+ 测试

```shell script
curl 127.0.0.1:8180/hello
```

+ 日志


```
{"level":"INFO","time":"2020-11-07 17:10:58.967426","caller":"ginLogger/logger.go:74","msg":"/hello","status":200,"method":"GET","path":"/hello","query":"name=wlf","ip":"127.0.0.1","user-agent":"curl/7.68.0","resId":""}
{"level":"DEBUG","time":"2020-11-07 17:10:58.967465","caller":"example/main.go:31","msg":"this is hello func","user":"q1mi","age":18}
{"level":"INFO","time":"2020-11-07 17:10:58.967482","caller":"ginLogger/logger.go:85","msg":"/hello","status":200,"method":"GET","path":"/hello","query":"name=wlf","ip":"127.0.0.1","user-agent":"curl/7.68.0","errors":"","costTime(s)":0.000070732,"resId":""}
```