package web

import (
	"fmt"
	"github.com/czy21/ndisk/constant"
	"github.com/czy21/ndisk/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func LogFormatter() func(param gin.LogFormatterParams) string {
	return func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		return fmt.Sprintf("%v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			param.TimeStamp.Format(model.StandardFormat),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}
}

func LogDav(fnName string, msg string) {
	log.Debugf("[DAV] %10s %s", fnName, msg)
}

func LogChunk(
	fnName string,
	fileName string,
	fileSize int64,
	chunks int,
	chunkL int,
	chunkI int,
	rangeS int64,
	rangeE int64,
	extension string) {
	chunkArr := []interface{}{
		constant.HttpExtraFileSize, fileSize,
		constant.HttpExtraChunks, chunks,
		constant.HttpExtraChunkL, chunkL,
		constant.HttpExtraChunkI, chunkI,
		constant.HttpExtraRangeS, rangeS,
		constant.HttpExtraRangeE, rangeE,
	}
	var msg string
	for i := 0; i < len(chunkArr)/2; i++ {
		msg += fmt.Sprintf(" %s: %d", chunkArr[i*2], chunkArr[i*2+1])
	}
	LogDav(fnName, fmt.Sprintf("%s %s %s", fileName, msg, extension))
}
