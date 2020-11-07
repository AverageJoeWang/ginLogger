package ginLogger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

//func Debug(c *gin.Context, msg string)  {
//	zap.L().Debug(msg, zap.String(""),)
//}