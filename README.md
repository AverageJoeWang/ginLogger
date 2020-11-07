

### 日志组件功能

+ 1.允许自定义日志等级,小于该等级的日志都不打印,done
+ 2.允许自定义日志字段，全局打印,todo
+ 3.封装对外暴露的接口，越简单越好,done
+ 4.允许管理日志文件路径、名字、大小、保存时间、打印格式、请求行,done
+ 5.允许自定义带的参数，类型，固定日志参数，生成独一无二的请求id，获取网关id等内容,todo
+ 6.支持json和console两种格式打印，done


## 使用

+ 下载

```
go get github.com/AverageJoeWang/ginLogger
```

+ 使用


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

## 说明

```go
type LogConfig struct {
	Level        string `json:"level"`       //default debug
	Filename     string `json:"filename"`    //default "app.log"
	MaxSize      int    `json:"maxsize"`     //the maximum size in megabytes of the log file before it gets rotated. It defaults to 200 megabytes.
	MaxAge       int    `json:"max_age"`     //the maximum number of days to retain old log files,default 7 days
	MaxBackups   int    `json:"max_backups"` //is the maximum number of old log files to retain.
	TimeLocation string `json:"time_zone"`   //time zone, default "Asia/Chongqing
	TimeFormat   string `json:"time_format"` //default 2006-01-02 15:04:05.000000
	LogFormat 	string	`json:"log_format"`//there two mode such as json/console，default json
}
```