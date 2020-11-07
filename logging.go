package ginLogger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Trace(c *gin.Context, msg string) {
	if logLevelMap["TRACE"] >= logLevelMap[logLevel] {
		zap.L().DPanic(msg, zap.String(resIdKey, getCommonParam(c)))
	}
}

func Debug(c *gin.Context, msg string) {
	if logLevelMap["DEBUG"] >= logLevelMap[logLevel] {
		zap.L().Debug(msg, zap.String(resIdKey, getCommonParam(c)))
	}
}

func Info(c *gin.Context, msg string) {
	if logLevelMap["INFO"] >= logLevelMap[logLevel] {
		zap.L().Info(msg, zap.String(resIdKey, getCommonParam(c)))
	}
}

func Warn(c *gin.Context, msg string) {
	if logLevelMap["WARN"] >= logLevelMap[logLevel] {
		zap.L().Warn(msg, zap.String(resIdKey, getCommonParam(c)))
	}
}

func Error(c *gin.Context, msg string) {
	if logLevelMap["ERROR"] >= logLevelMap[logLevel] {
		zap.L().Error(msg, zap.String(resIdKey, getCommonParam(c)))
	}
}

func Fatal(c *gin.Context, msg string) {
	if logLevelMap["FATAL"] >= logLevelMap[logLevel] {
		zap.L().Fatal(msg, zap.String(resIdKey, getCommonParam(c)))
	}
}

func getCommonParam(c *gin.Context) string {
	var cp commonParam
	cp.ResId, _ = c.GetQuery(resIdKey)
	return cp.ResId
}
