package util

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
)

func CopyN(dst io.Writer, src io.Reader, buf []byte) (n int64, err error) {
	for {
		nr, er := io.ReadFull(src, buf)
		if nr > 0 {
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

func GetChunk(name string, fileSize int64, chunkSize int64, extra map[string]interface{}) (int64, int, int64, int64) {
	var chunks int64
	chunkI := 0
	rangeL := int64(0)
	rangeR := chunkSize
	if extra["chunkI"] != nil {
		chunkI = extra["chunkI"].(int) + 1
	}
	if extra["rangeR"] != nil {
		v := extra["rangeR"].(int64)
		rangeL = v
		rangeR += v
	}
	if extra["chunks"] == nil {
		chunks = int64(math.Ceil(float64(fileSize)) / float64(chunkSize))
	} else {
		chunks = extra["chunks"].(int64)
	}
	extra["chunks"] = chunks
	extra["chunkI"] = chunkI
	extra["rangeR"] = rangeR
	extra["chunks"] = chunks
	log.Debugf("%s chunks: %d chunkSize: %d chunkI: %d rangeL: %d rangeR: %d", name, chunks, chunkSize, chunkI, rangeL, rangeR)
	return chunks, chunkI, rangeL, rangeR
}
