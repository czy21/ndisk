package _189

import (
	"encoding/hex"
	"fmt"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider/base"
	"github.com/czy21/ndisk/util"
	"hash"
	"io"
	"io/fs"
	"path"
	"strings"
)

type File struct {
	base.FileBase
}

func (f File) DownloadCreate() (dUrl string, fileSize int64, err error) {
	api := API{File: f.File}
	fileInfo, err := f.FS.GetFileInfo(f.Ctx, f.File.Target.Name, f.File)
	fileInfoVO, err := api.GetFileInfoById(fileInfo.Id)
	return fileInfoVO.FileDownloadUrl, fileInfoVO.Size, err
}
func (f File) DownloadChunk(dUrl string, p []byte, rangeStart int64, rangeEnd int64) (n int, err error) {
	req := http2.GetClient().NewRequest()
	req.SetHeader("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))
	res, _ := req.Get(dUrl)
	return copy(p, res.Body()), err
}

func (f File) Readdir(count int) ([]fs.FileInfo, error) {
	api := API{File: f.File}
	fileInfo, _ := f.FS.GetFileInfo(f.Ctx, f.File.Target.Name, f.File)
	folder, err := api.GetFolderById(fileInfo.Id)
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
	}
	for _, t := range folder.Files {
		fileInfos = append(fileInfos, model.FileInfoDelegate{
			FileInfo: model.FileInfo{
				Name: t.Name,
			},
		})
	}
	return fileInfos, err
}

func (f File) UploadCreate(md5Hash hash.Hash) (fileId string, err error) {
	api := API{File: f.File}
	fileSize := f.UploadFileSize()
	d, fName := path.Split(f.File.Target.Name)
	fileInfo, err := f.FS.GetFileInfo(f.Ctx, d, f.File)
	var fileMd5 string
	if fileSize == 0 {
		fileMd5 = hex.EncodeToString(md5Hash.Sum(nil))
	}
	res, err := api.CreateFile(fileInfo.Id, fName, fileSize, fileMd5)
	return res.UploadFileId, err
}

func (f File) UploadCommit(fileId string, md5Hash hash.Hash, md5s []string, chunkLen int) (err error) {
	api := API{File: f.File}
	fileSize := f.UploadFileSize()
	fileMd5 := hex.EncodeToString(md5Hash.Sum(nil))
	sliceMd5 := fileMd5
	if fileSize > int64(chunkLen) {
		sliceMd5 = util.GetMD5Encode(strings.Join(md5s, "\n"))
	}
	err = api.CommitFile(fileId, fileSize, fileMd5, sliceMd5)
	return err
}

func (f File) UploadChunk(fileId string, b []byte, md5Bytes []byte, index int) (n int, err error) {
	api := API{File: f.File}
	chunkLen := len(b)
	err = api.UploadChunk(fileId, b, md5Bytes, index+1)
	return chunkLen, err
}

//WriteTo CopyTo
func (f File) WriteTo(w io.Writer) (n int64, err error) {
	api := API{File: f.File}
	_, srcFName := path.Split(f.File.Target.Name)
	dstD, _ := path.Split(w.(File).Name())
	srcFileInfo, err := f.FS.GetFileInfo(f.Ctx, f.File.Target.Name, f.File)
	dstFileInfo, err := f.FS.GetFileInfo(f.Ctx, dstD, f.File)
	err = api.Copy(srcFileInfo.Id, srcFName, srcFileInfo.IsDir, dstFileInfo.Id)
	return srcFileInfo.Size, err
}
