package _189

import (
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/util"
	"io/fs"
	"os"
	"path"
	"time"
)

type FileInfo struct {
	name       string
	size       int64
	mode       os.FileMode
	modTime    time.Time
	isDir      bool
	sys        any
	remoteName string
}

func (c FileInfo) Name() string {
	return c.name
}

func (c FileInfo) Size() int64 {
	return c.size
}

func (c FileInfo) Mode() fs.FileMode {
	return c.mode
}

func (c FileInfo) ModTime() time.Time {
	return c.modTime
}

func (c FileInfo) IsDir() bool {
	return c.isDir
}

func (c FileInfo) Sys() any {
	return c.sys
}

type File struct {
	name string
	pctx model.ProviderContext
	env  map[string]interface{}
}

func (f File) Close() error {
	return nil
}

func (f File) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	panic("implement me")
}

func (f File) Readdir(count int) ([]fs.FileInfo, error) {
	var ret FileListAORes
	var client = util.HttpUtil{}.NewClient()
	client.SetHeader("accept", "application/json;charset=UTF-8")
	client.SetHeader("cookie", "s_fid=1F2141B769232BD6-27945D9DC425F8FF; lvid=a8761a577d0946ea770ac65cdf877c2f; nvid=1; trkId=645D4484-F660-49CE-9983-355F77E5D334; _gscu_1708861450=45501776bz4o6613; svid=40D644AB56B89B6BEED64A023263A993; userId=201%7C20170100000261869905; apm_ua=8B11E0A1C25A29CA8CD6B530E64C5294; apm_ct=20220620141608000; apm_ip=116.247.110.46; apm_uid=35A656C0E78BB334950E945E5DFFC2E1; apm_sid=2FEFD7ABE7002318AE3829E902CACA81; JSESSIONID=866FF0F3B3B373DA48810B9AB109A9F9; COOKIE_LOGIN_USER=81ACFDB17EFBF1BDE2E6339CE631B77F236E3B80BDC0C72440B457B650ED35EA1B5E21388310F531FCE9B745EC8B61F728687EACB9DD00C0BB7E745A83867C0D55BA1331")
	client.Get(queryFolder(f.env[f.name].(FileInfo).remoteName), &ret)
	var fileInfos []fs.FileInfo
	for _, t := range ret.FileListAO.Files {
		fileInfos = append(fileInfos, FileInfo{
			name: t.Name,
			size: t.Size,
		})
	}
	for _, t := range ret.FileListAO.Folders {
		fileInfos = append(fileInfos, FileInfo{
			name:  t.Name,
			isDir: true,
		})
	}
	return fileInfos, nil
}

func (f File) Stat() (fs.FileInfo, error) {
	//ret := FileListAORes{}
	var client = util.HttpUtil{}.NewClient()
	client.SetHeader("accept", "application/json;charset=UTF-8")
	client.SetHeader("cookie", "s_fid=1F2141B769232BD6-27945D9DC425F8FF; lvid=a8761a577d0946ea770ac65cdf877c2f; nvid=1; trkId=645D4484-F660-49CE-9983-355F77E5D334; _gscu_1708861450=45501776bz4o6613; svid=40D644AB56B89B6BEED64A023263A993; userId=201%7C20170100000261869905; apm_ua=8B11E0A1C25A29CA8CD6B530E64C5294; apm_ct=20220620141608000; apm_ip=116.247.110.46; apm_uid=35A656C0E78BB334950E945E5DFFC2E1; apm_sid=2FEFD7ABE7002318AE3829E902CACA81; JSESSIONID=866FF0F3B3B373DA48810B9AB109A9F9; COOKIE_LOGIN_USER=81ACFDB17EFBF1BDE2E6339CE631B77F236E3B80BDC0C72440B457B650ED35EA1B5E21388310F531FCE9B745EC8B61F728687EACB9DD00C0BB7E745A83867C0D55BA1331")
	fileInfo := FileInfo{}
	if f.name == path.Join("/", f.pctx.Meta.Name) {
		fileInfo.isDir = true
		return fileInfo, nil
	}
	if val, ok := f.env[f.name]; ok {
		fileInfo.isDir = val.(FileInfo).isDir
	}
	return fileInfo, nil
}

func (f File) Write(p []byte) (n int, err error) {
	panic("implement me")
}
