package s3

import (
	"github.com/czy21/ndisk/provider/base"
	"hash"
	"io"
	"io/fs"
)

type File struct {
	base.FileBase
}

func (f File) DownloadCreate() (dUrl string, fileSize int64, err error) {
	panic("")
}
func (f File) DownloadChunk(dUrl string, p []byte, rangeStart int64, rangeEnd int64) (n int, err error) {
	panic("")
}

func (f File) Readdir(count int) ([]fs.FileInfo, error) {
	panic("")
}

func (f File) UploadCreate(md5Hash hash.Hash) (fileId string, err error) {
	panic("")
}

func (f File) UploadCommit(fileId string, md5Hash hash.Hash, md5s []string, chunkLen int) (err error) {
	panic("")
}

func (f File) UploadChunk(fileId string, b []byte, md5Bytes []byte, index int) (n int, err error) {
	panic("")
}

//WriteTo CopyTo
func (f File) WriteTo(w io.Writer) (n int64, err error) {
	panic("")
}
