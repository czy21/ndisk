package util

import (
	"fmt"
	"github.com/czy21/ndisk/constant"
	log "github.com/sirupsen/logrus"
	"io"
)

func Copy(dst io.Writer, src io.Reader) (n int64, err error) {
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	return n, err
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
	log.Debugf("%s %s %s %s", fnName, fileName, msg, extension)
}
