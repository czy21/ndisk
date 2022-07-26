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

type WriterProxy interface {
	Create(md5Hash hash.Hash) (string, error)
	Chunk(fileId string, p []byte, md5Bytes []byte, index int) (n int, err error)
	Commit(fileId string, md5Hash hash.Hash, md5s []string, chunkLen int) error
	FileName() string
	FileSize() int64
}

func WriteFull(dst io.Writer, src io.Reader, n int) (written int64, err error) {
	wt, ok := dst.(WriterProxy)
	if !ok {
		return 0, errors.New("no implement WriterProxy interface")
	}
	buf := make([]byte, n)
	md5s := make([]string, 0)
	md5Hash := md5.New()
	fileId, err := wt.Create(md5Hash)
	fileName := wt.FileName()
	fileSize := wt.FileSize()
	if err != nil {
		return 0, err
	}
	chunkI := 0
	chunkL := len(buf)
	chunks := int(math.Max(1, math.Ceil(float64(fileSize)/float64(cap(buf)))))
	rangeS := int64(0)
	rangeE := int64(0)
	for {
		nr, er := io.ReadFull(src, buf)
		rangeS = rangeE
		rangeE += int64(nr)
		logChunk(fileName, fileSize, chunks, chunkL, chunkI, rangeS, rangeE)
		if nr > 0 || fileSize == 0 {
			bufBytes := buf[0:nr]
			md5Hash.Write(bufBytes)
			md5Bytes := GetMd5Bytes(bufBytes)
			md5s = append(md5s, strings.ToUpper(hex.EncodeToString(md5Bytes)))
			nw, ew := wt.Chunk(fileId, bufBytes, md5Bytes, chunkI)
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write result")
				}
			}
			written += int64(nw)
			if ew == io.EOF || fileSize == written {
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
		chunkI++
	}
	err = wt.Commit(fileId, md5Hash, md5s, chunkL)
	return written, err
}

func ReadFull(dst io.Writer, src io.Reader, n int) (written int64, err error) {
	buf := make([]byte, n)
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
			written += int64(nw)
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
	return written, err
}

func logChunk(fileName string, fileSize int64, chunks int, chunkL int, chunkI int, rangeS int64, rangeE int64) {
	chunkArr := []interface{}{
		constant.HttpExtraFileSize, fileSize,
		constant.HttpExtraChunks, chunks,
		constant.HttpExtraChunkL, chunkL,
		constant.HttpExtraChunkI, chunkI,
		constant.HttpExtraRangeS, rangeS,
		constant.HttpExtraRangeE, rangeE,
	}
	var chunkLog string
	for i := 0; i < len(chunkArr)/2; i++ {
		chunkLog += fmt.Sprintf(" %s: %d", chunkArr[i*2], chunkArr[i*2+1])
	}
	log.Infof("%s %s", fileName, chunkLog)
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
