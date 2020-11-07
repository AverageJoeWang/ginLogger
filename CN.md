

gin的日志中间件，基于zap和lumberjack

## 功能

+ [X] 允许自定义日志等级,小于该等级的日志都不打印
+ [X] 允许自定义日志字段，全局打印,todo
+ [X] 封装对外暴露的接口，越简单越好,done
+ [X] 允许管理日志文件路径、名字、大小、保存时间、打印格式、请求行
+ [ ] 允许自定义带的参数，类型，固定日志参数，获取网关id等内容
+ [X] 支持json和console两种格式打印


## 如何使用

+ 下载

```
go get -u github.com/AverageJoeWang/ginLogger
```

+ 用法

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

+ 运行

```
go run main.go
```



+ 测试

```shell script
curl 127.0.0.1:8180/hello?resId=1000
```

+ 日志文件

```
{"level":"INFO","time":"2020-11-07 17:10:58.967426","caller":"ginLogger/logger.go:74","msg":"/hello","status":200,"method":"GET","path":"/hello","query":"name=wlf","ip":"127.0.0.1","user-agent":"curl/7.68.0","resId":""}
{"level":"DEBUG","time":"2020-11-07 17:10:58.967465","caller":"example/main.go:31","msg":"this is hello func","user":"q1mi","age":18}
{"level":"INFO","time":"2020-11-07 17:10:58.967482","caller":"ginLogger/logger.go:85","msg":"/hello","status":200,"method":"GET","path":"/hello","query":"name=wlf","ip":"127.0.0.1","user-agent":"curl/7.68.0","errors":"","costTime(s)":0.000070732,"resId":""}
```

## 其他说明

+ "level": 日志等级，默认为debug
+ "filename": 日志文件名，需要自定义后缀，默认为"app.log"
+ "maxsize": 单个日志最大容量，默认为200MB
+ "max_age": 归档以后的日志保存天数，默认为7天
+ "max_backups": 最大归档日志数量，默认为30个
+ "time_zone": 时区，默认为"Asia/Chongqing
+ "time_format": 日志中的时间格式，默认为"2006-01-02 15:04:05.000000"
+ "log_format": 日志打印格式，分为json和console两种类型，默认为json格式

