

[中文说明](./CN.md)

gin-logger middleware based on zap,lumberjack

## Features

+ [X] Allow self config log level, it will not print if less than config level, TRACE<DEBUG<INFO<WARN<ERROR<FATAL
+ [ ] Allow custom field and print in every log
+ [X] Easy to use
+ [X] Allow manage log file name,size,backups,print format
+ [ ] Fixed parameters, also get resId from gateway
+ [X] Support json/console print format


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

