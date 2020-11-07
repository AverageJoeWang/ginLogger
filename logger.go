package ginLogger

import (
	"github.com/gin-gonic/gin"
	"github.com/snowair/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var ginLogger *zap.Logger

const (
	resIdKey     = "resId"
	statusKey    = "status"
	methodKey    = "method"
	pathKey      = "path"
	queryKey     = "query"
	ipKey        = "ip"
	userAgentKey = "user-agent"
	errorsKey    = "errors"
	costTimeKey  = "costTime(s)"
)

var timeLocation string
var logLevel string
var mutex sync.Mutex
var logLevelMap = map[string]int{
	"TRACE": 1,
	"DEBUG": 2,
	"INFO":  3,
	"WARN":  4,
	"ERROR": 5,
	"FATAL": 6,
}

//init logger
func InitLogger(config *LogConfig) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	//check config
	checkDefaultSetting(config)
	//init
	writeSyncer := getLogWriter(config.Filename, config.MaxSize, config.MaxBackups, config.MaxAge)
	encoder := getEncoder(config)
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(config.Level))
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	ginLogger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(ginLogger)
	return nil
}

func checkDefaultSetting(config *LogConfig) {
	//check time location
	if config.TimeLocation == "" {
		config.TimeLocation = "Asia/Chongqing"
	}
	timeLocation = config.TimeLocation
	//check time format
	if config.TimeFormat == "" {
		config.TimeFormat = "2006-01-02 15:04:05.000000"
	}
	//check maxsize
	if config.MaxSize == 0 {
		config.MaxSize = 200
	}
	//check maxAge
	if config.MaxAge == 0 {
		config.MaxAge = 7
	}
	//check max_backups
	config.MaxBackups = 30
	if config.MaxBackups == 0 {
	}
	//check level
	if config.Level == "" {
		config.Level = "DEBUG"
	}
	logLevel = config.Level
	//check log file name
	if config.Filename == "" {
		config.Filename = "app.log"
	}
	//check log format default json
	if config.LogFormat == "" {
		config.LogFormat = "json"
	}
}

//self config
func getEncoder(config *LogConfig) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(config.TimeFormat)
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	if config.LogFormat == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}else {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		loc, _ := time.LoadLocation(timeLocation)
		start := time.Now().In(loc)
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		gateResId, _ := c.GetQuery(resIdKey)
		ginLogger.Info(path,
			zap.Int(statusKey, c.Writer.Status()),
			zap.String(methodKey, c.Request.Method),
			zap.String(pathKey, path),
			zap.String(queryKey, query),
			zap.String(ipKey, c.ClientIP()),
			zap.String(userAgentKey, c.Request.UserAgent()),
			zap.String(resIdKey, gateResId),
		)
		c.Next()
		cost := time.Since(start)
		ginLogger.Info(path,
			zap.Int(statusKey, c.Writer.Status()),
			zap.String(methodKey, c.Request.Method),
			zap.String(pathKey, path),
			zap.String(queryKey, query),
			zap.String(ipKey, c.ClientIP()),
			zap.String(userAgentKey, c.Request.UserAgent()),
			zap.String(errorsKey, c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration(costTimeKey, cost),
			zap.String(resIdKey, gateResId),
		)
	}
}

// GinRecovery
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					ginLogger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					ginLogger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					ginLogger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
