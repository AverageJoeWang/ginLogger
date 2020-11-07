

## Features

+ [X] 允许自定义日志等级,小于该等级的日志都不打印
+ [X] 允许自定义日志字段，全局打印,todo
+ [X] 封装对外暴露的接口，越简单越好,done
+ [X] 允许管理日志文件路径、名字、大小、保存时间、打印格式、请求行
+ [ ] 允许自定义带的参数，类型，固定日志参数，获取网关id等内容
+ [X] 支持json和console两种格式打印


## How to Use

+ download

```
go get -u github.com/AverageJoeWang/ginLogger
```

+ use

    - config.json
    
```
{
  "mode": "debug",
  "port": 8180,
  "log": {
    "level": "INFO",
    "filename": "app.log",
    "log_format": "json",
    "maxsize": 200,
    "max_age": 7,
    "max_backups": 10,
    "time_location": "Asia/Chongqing",
    "time_format": "2006-01-02 15:04:05.000000"
  }
}
```

    + main.go

```go
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
```

+ run

```
go run main.go
```



+ test

```shell script
curl 127.0.0.1:8180/hello?resId=1000
```

+ log file

```
{"level":"INFO","time":"2020-11-07 17:10:58.967426","caller":"ginLogger/logger.go:74","msg":"/hello","status":200,"method":"GET","path":"/hello","query":"name=wlf","ip":"127.0.0.1","user-agent":"curl/7.68.0","resId":""}
{"level":"DEBUG","time":"2020-11-07 17:10:58.967465","caller":"example/main.go:31","msg":"this is hello func","user":"q1mi","age":18}
{"level":"INFO","time":"2020-11-07 17:10:58.967482","caller":"ginLogger/logger.go:85","msg":"/hello","status":200,"method":"GET","path":"/hello","query":"name=wlf","ip":"127.0.0.1","user-agent":"curl/7.68.0","errors":"","costTime(s)":0.000070732,"resId":""}
```

## Intro

+ "level": default debug
+ "filename": default "app.log"
+ "maxsize": the maximum size in megabytes of the log file before it gets rotated. It defaults to 200 megabytes.
+ "max_age": the maximum number of days to retain old log files,default 7 days
+ "max_backups": is the maximum number of old log files to retain.
+ "time_zone": time zone, default "Asia/Chongqing
+ "time_format": default "2006-01-02 15:04:05.000000"
+ "log_format": there two mode such as json/console，default json

