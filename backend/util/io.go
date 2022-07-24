package util

import (
	"errors"
	"github.com/czy21/ndisk/constant"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
)

func CopyN(dst io.Writer, src io.Reader, buf []byte) (n int64, err error) {
	for {
		nr, er := io.ReadFull(src, buf)
		if nr > 0 || er == io.EOF {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write result")
				}
			}
			n += int64(nw)
			if ew == io.EOF {
				break
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return n, err
}

func GetChunk(name string, fileSize int64, chunkLen int64, extra map[string]interface{}) (int, int, int64, int64) {
	var chunks int
	chunkI := 0
	rangeS := int64(0)
	rangeE := chunkLen
	if extra[constant.HttpExtraChunkI] != nil {
		chunkI = extra[constant.HttpExtraChunkI].(int) + 1
	}
	if extra[constant.HttpExtraRangeE] != nil {
		v := extra[constant.HttpExtraRangeE].(int64)
		rangeS = v
		rangeE += v
	}
	if extra[constant.HttpExtraChunks] == nil {
		if chunkLen == 0 {
			chunks = 1
		} else {
			chunks = int(math.Max(1, math.Ceil(float64(fileSize))/float64(chunkLen)))
		}
	} else {
		chunks = extra[constant.HttpExtraChunks].(int)
	}
	extra[constant.HttpExtraChunks] = chunks
	extra[constant.HttpExtraChunkI] = chunkI
	extra[constant.HttpExtraRangeE] = rangeE
	extra[constant.HttpExtraChunks] = chunks
	log.Infof("%s chunks: %d chunkL: %d chunkI: %d rangeS: %d rangeE: %d", name, chunks, chunkLen, chunkI, rangeS, rangeE)
	return chunks, chunkI, rangeS, rangeE
}
