package _189

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/constant"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider/base"
	"github.com/czy21/ndisk/util"
	"io"
	"io/fs"
	"math"
	"path"
	"strconv"
	"strings"
	"time"
)

type File struct {
	base.FileBase
}

func (f File) Readdir(count int) ([]fs.FileInfo, error) {
	api := API{File: f.File}
	fileInfo, _ := f.FS.GetFileInfo(f.Ctx, f.File.Target.Name, f.File)
	folder, err := api.GetObjectsById(fileInfo.Id, "")
	if err != nil {
		return nil, err
	}
	var fileInfos []fs.FileInfo
	for _, t := range folder.Folders {
		fileInfos = append(fileInfos, model.FileInfoDelegate{
			FileInfo: model.FileInfo{
				Name:  t.Name,
				IsDir: true,
			},
		})
		fi := model.FileInfo{
			Name:    path.Join(f.File.Target.Name, t.Name),
			ModTime: time.Time(t.UpdateDate).Add(-8 * time.Hour),
			Id:      strconv.FormatInt(t.Id, 10),
			IsDir:   true,
		}
		cache.Client.SetObj(f.Ctx, cache.GetFileInfoCacheKey(fi.Name), &fi)
	}
	for _, t := range folder.Files {
		fileInfos = append(fileInfos, model.FileInfoDelegate{
			FileInfo: model.FileInfo{
				Name: t.Name,
			},
		})
		fi := model.FileInfo{
			Name:    path.Join(f.File.Target.Name, t.Name),
			ModTime: time.Time(t.UpdateDate).Add(-8 * time.Hour),
			Size:    t.Size,
			Id:      strconv.FormatInt(t.Id, 10),
		}
		cache.Client.SetObj(f.Ctx, cache.GetFileInfoCacheKey(fi.Name), &fi)
	}
	return fileInfos, err
}

// ReadFrom Put
func (f File) ReadFrom(r io.Reader) (written int64, err error) {
	api := API{File: f.File}
	fileInfo, err := f.FS.GetFileInfo(f.Ctx, f.File.Target.Dir, f.File)
	buf := make([]byte, 1024*1024*10)
	md5s := make([]string, 0)
	md5Hash := md5.New()
	fileMd5 := ""
	if f.UploadFileSize() == 0 {
		fileMd5 = hex.EncodeToString(md5Hash.Sum(nil))
	}
	res, err := api.CreateFile(fileInfo.Id, f.File.Target.BaseName, f.UploadFileSize(), fileMd5)
	if err != nil {
		return written, err
	}
	if res.UploadFileId == "" {
		err = errors.New("create file fail")
		return written, err
	}
	chunkL := len(buf)
	chunks := int(math.Max(1, math.Ceil(float64(f.UploadFileSize())/float64(cap(buf)))))
	rangeS := int64(0)
	rangeE := int64(0)
	for i := 0; i < chunks; i++ {
		nr, _ := io.ReadFull(r, buf)
		rangeS = rangeE
		rangeE += int64(nr)
		util.LogChunk("Put", f.Name(), f.UploadFileSize(), chunks, chunkL, i, rangeS, rangeE, fmt.Sprintf("nr: %d", nr))
		if nr > 0 {
			bufBytes := buf[0:nr]
			md5Hash.Write(bufBytes)
			md5Bytes := util.GetMd5Bytes(bufBytes)
			md5s = append(md5s, strings.ToUpper(hex.EncodeToString(md5Bytes)))
			_ = api.UploadChunk(res.UploadFileId, bufBytes, md5Bytes, i+1)
			written += int64(nr)
		}
	}
	fileMd5 = hex.EncodeToString(md5Hash.Sum(nil))
	sliceMd5 := fileMd5
	if f.UploadFileSize() > int64(chunkL) {
		sliceMd5 = util.GetMD5Encode(strings.Join(md5s, "\n"))
	}
	err = api.CommitFile(res.UploadFileId, f.UploadFileSize(), fileMd5, sliceMd5)
	return written, err
}

// WriteTo Get
func (f File) WriteTo(w io.Writer) (written int64, err error) {
	api := API{File: f.File}
	extra := f.Ctx.Value(constant.HttpExtra).(map[string]interface{})
	if extra[constant.HttpExtraMethod] == "COPY" {
		_, srcFName := path.Split(f.File.Target.Name)
		dstD, _ := path.Split(w.(File).Name())
		srcFileInfo, err := f.FS.GetFileInfo(f.Ctx, f.File.Target.Name, f.File)
		dstFileInfo, err := f.FS.GetFileInfo(f.Ctx, dstD, f.File)
		err = api.Copy(srcFileInfo.Id, srcFName, srcFileInfo.IsDir, dstFileInfo.Id)
		return srcFileInfo.Size, err
	}
	fileInfo, err := f.FS.GetFileInfo(f.Ctx, f.File.Target.Name, f.File)
	fileInfoVO, err := api.GetFileById(fileInfo.Id)
	req := http2.GetClient().NewRequest()
	buf := make([]byte, 1024*1024*f.File.ProviderFolder.Account.GetBuf)
	chunkL := len(buf)
	chunks := int(math.Max(1, math.Ceil(float64(fileInfoVO.Size)/float64(cap(buf)))))
	rangeS := int64(0)
	rangeE := int64(0)
	for i := 0; i < chunks; i++ {
		remain := fileInfoVO.Size - written
		rangeS = rangeE
		if remain > int64(chunkL) {
			rangeE += int64(chunkL)
		} else {
			rangeE += remain
		}
		req.SetHeader("Range", fmt.Sprintf("bytes=%d-%d", rangeS, rangeE))
		res, _ := req.Get(fileInfoVO.FileDownloadUrl)
		util.LogChunk("Get", f.Name(), fileInfoVO.Size, chunks, chunkL, i, rangeS, rangeE, "")
		nw, _ := w.Write(buf[0:copy(buf, res.Body())])
		written += int64(nw)
	}
	return written, err
}
