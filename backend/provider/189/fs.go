package _189

import (
	"fmt"
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/util"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type FileSystem struct {
}

const localDir = "data"

func (fs FileSystem) Mkdir(ctx context.Context, pctx model.ProviderContext, name string, perm os.FileMode) error {
	return webdav.Dir(localDir).Mkdir(ctx, name, perm)
}
func (fs FileSystem) OpenFile(ctx context.Context, pctx model.ProviderContext, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return File{name: name, pctx: pctx, env: ctx.Value("env").(map[string]interface{})}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, pctx model.ProviderContext, name string) error {
	return webdav.Dir(localDir).RemoveAll(ctx, name)
}
func (fs FileSystem) Rename(ctx context.Context, pctx model.ProviderContext, oldName, newName string) error {
	return webdav.Dir(localDir).Rename(ctx, oldName, newName)
}
func (fs FileSystem) Stat(ctx context.Context, pctx model.ProviderContext, name string) (os.FileInfo, error) {
	env := ctx.Value("env").(map[string]interface{})
	ret := FileListAORes{}
	var client = util.HttpUtil{}.NewClient()
	client.SetHeader("accept", "application/json;charset=UTF-8")
	client.SetHeader("cookie", "s_fid=1F2141B769232BD6-27945D9DC425F8FF; lvid=a8761a577d0946ea770ac65cdf877c2f; nvid=1; trkId=645D4484-F660-49CE-9983-355F77E5D334; _gscu_1708861450=45501776bz4o6613; svid=40D644AB56B89B6BEED64A023263A993; userId=201%7C20170100000261869905; apm_ua=8B11E0A1C25A29CA8CD6B530E64C5294; apm_ct=20220620141608000; apm_ip=116.247.110.46; apm_uid=35A656C0E78BB334950E945E5DFFC2E1; apm_sid=2FEFD7ABE7002318AE3829E902CACA81; JSESSIONID=866FF0F3B3B373DA48810B9AB109A9F9; COOKIE_LOGIN_USER=81ACFDB17EFBF1BDE2E6339CE631B77F236E3B80BDC0C72440B457B650ED35EA1B5E21388310F531FCE9B745EC8B61F728687EACB9DD00C0BB7E745A83867C0D55BA1331")
	fileInfo := FileInfo{}
	if name == path.Join("/", pctx.Meta.Name)+"/" {
		client.Get(fmt.Sprintf("https://cloud.189.cn/api/open/file/listFiles.action?noCache=0.7362081385378736&pageSize=60&pageNum=1&mediaType=0&folderId=%s&iconOption=5&orderBy=lastOpTime&descending=true", pctx.Meta.RemoteName), &ret)
		fileInfo.isDir = true
	} else {
		client.Get(fmt.Sprintf("https://cloud.189.cn/api/open/file/listFiles.action?noCache=0.7362081385378736&pageSize=60&pageNum=1&mediaType=0&folderId=%s&iconOption=5&orderBy=lastOpTime&descending=true", env["remoteName"]), &ret)
		for _, t := range ret.FileListAO.Folders {
			if t.Name == filepath.Base(name) {
				fileInfo.isDir = true
				env["isDir"] = true
				env["remoteName"] = strconv.FormatInt(t.Id, 10)
			}
		}
	}
	return fileInfo, nil
}
