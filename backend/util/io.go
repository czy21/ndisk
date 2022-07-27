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

type UpDownWriter interface {
	FileName() string
	UploadFileSize() int64
	UploadCreate(md5Hash hash.Hash) (string, error)
	UploadChunk(fileId string, p []byte, md5Bytes []byte, index int) (n int, err error)
	UploadCommit(fileId string, md5Hash hash.Hash, md5s []string, chunkLen int) error
	DownloadCreate() (string, int64, error)
	DownloadChunk(dUrl string, p []byte, rangeStart int64, rangeEnd int64) (m int, err error)
}

func WriteFull(dst io.Writer, src io.Reader, n int) (written int64, err error) {
	wt, ok := dst.(UpDownWriter)
	if !ok {
		return 0, errors.New("no implement UpDownWriter interface")
	}
	fileSize := wt.UploadFileSize()
	if fileSize == 0 {
		return 0, nil
	}
	// cache get u:/189/test/t1 ret: fileId,fileSize,writtenSize,last
	buf := make([]byte, n)
	md5s := make([]string, 0)
	md5Hash := md5.New()
	fileId, err := wt.UploadCreate(md5Hash)
	fileName := wt.FileName()
	if err != nil {
		return 0, err
	}
	chunkL := len(buf)
	chunks := int(math.Max(1, math.Ceil(float64(fileSize)/float64(cap(buf)))))
	rangeS := int64(0)
	rangeE := int64(0)
	for i := 0; i < chunks; i++ {
		nr, er := io.ReadFull(src, buf)
		rangeS = rangeE
		rangeE += int64(chunkL)
		logChunk("Put", fileName, fileSize, chunks, chunkL, i, rangeS, rangeE)
		if nr > 0 {
			md5Hash.Write(buf)
			md5Bytes := GetMd5Bytes(buf)
			md5s = append(md5s, strings.ToUpper(hex.EncodeToString(md5Bytes)))
			nw, ew := wt.UploadChunk(fileId, buf, md5Bytes, i)
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
	}
	err = wt.UploadCommit(fileId, md5Hash, md5s, chunkL)
	return written, err
}

func ReadFull(dst io.Writer, src io.Reader, n int) (written int64, err error) {
	rt, ok := src.(UpDownWriter)
	if !ok {
		return 0, errors.New("no implement UpDownWriter interface")
	}
	dUrl, fileSize, err := rt.DownloadCreate()
	fileName := rt.FileName()
	buf := make([]byte, n)
	chunkL := len(buf)
	chunks := int(math.Max(1, math.Ceil(float64(fileSize)/float64(cap(buf)))))
	rangeS := int64(0)
	rangeE := int64(0)
	for i := 0; i < chunks; i++ {
		remain := fileSize - written
		rangeS = rangeE
		if remain > int64(chunkL) {
			rangeE += int64(chunkL)
		} else {
			rangeE += remain
		}
		nr, er := rt.DownloadChunk(dUrl, buf, rangeS, rangeE)
		logChunk("Get", fileName, fileSize, chunks, chunkL, i, rangeS, rangeE)
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

func logChunk(fnName, fileName string, fileSize int64, chunks int, chunkL int, chunkI int, rangeS int64, rangeE int64) {
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
	log.Infof("%s %s %s", fnName, fileName, chunkLog)
}
