package util

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/czy21/ndisk/constant"
	log "github.com/sirupsen/logrus"
	"hash"
	"io"
	"math"
	"strings"
)

type WriterMd5 interface {
	WriteMd5(p []byte, md5Bytes []byte, md5Hash hash.Hash, md5s []string) (n int, err error)
}

func CopyN(dst io.Writer, src io.Reader, buf []byte) (n int64, err error) {
	chunkMd5s := make([]string, 0)
	fileMd5Hash := md5.New()
	for {
		nr, er := io.ReadFull(src, buf)
		if nr > 0 || er == io.EOF {
			bufBytes := buf[0:nr]
			var nw int
			var ew error
			if wt, ok := dst.(WriterMd5); ok {
				fileMd5Hash.Write(bufBytes)
				chunkMd5Bytes := GetMd5Bytes(bufBytes)
				chunkMd5s = append(chunkMd5s, strings.ToUpper(hex.EncodeToString(chunkMd5Bytes)))
				nw, er = wt.WriteMd5(bufBytes, chunkMd5Bytes, fileMd5Hash, chunkMd5s)
			} else {
				nw, ew = dst.Write(bufBytes)
			}
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
	chunkArr := []interface{}{
		constant.HttpExtraFileSize, fileSize,
		constant.HttpExtraChunks, chunks,
		constant.HttpExtraChunkL, chunkLen,
		constant.HttpExtraChunkI, chunkI,
		constant.HttpExtraRangeS, rangeS,
		constant.HttpExtraRangeE, rangeE,
	}
	var chunkLog string
	for i := 0; i < len(chunkArr)/2; i++ {
		chunkLog += fmt.Sprintf(" %s: %d", chunkArr[i*2], chunkArr[i*2+1])
	}
	log.Infof("%s %s", name, chunkLog)
	return chunks, chunkI, rangeS, rangeE
}
