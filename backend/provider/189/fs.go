package _189

import (
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/util"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
)

type FileSystem struct{}

const localDir = "data"

func (FileSystem) Mkdir(ctx context.Context, provider model.Provider, name string, perm os.FileMode) error {
	return webdav.Dir(localDir).Mkdir(ctx, name, perm)
}
func (FileSystem) OpenFile(ctx context.Context, provider model.Provider, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return webdav.Dir(localDir).OpenFile(ctx, name, flag, perm)
}
func (FileSystem) RemoveAll(ctx context.Context, provider model.Provider, name string) error {
	return webdav.Dir(localDir).RemoveAll(ctx, name)
}
func (FileSystem) Rename(ctx context.Context, provider model.Provider, oldName, newName string) error {
	return webdav.Dir(localDir).Rename(ctx, oldName, newName)
}
func (FileSystem) Stat(ctx context.Context, provider model.Provider, name string) (os.FileInfo, error) {
	ret := make(map[string]interface{})
	var client = util.HttpUtil{}.NewClient()
	client.SetHeader("accept", " application/json;charset=UTF-8")
	client.SetHeader("cookie", "apm_ua=74B7FDC3A7244BA6FBFD4FC6669EFAFF; apm_ct=20220620210641000; apm_ip=218.81.3.182; apm_uid=9F0A8C1B5282C049DAF6ED52FD27EA97; apm_sid=C021764EA5E92F466B1B4EF2BE79537D; JSESSIONID=DC089721DA138A4DB20577C715696F95; COOKIE_LOGIN_USER=D21169CA90978527C84D625098572A1DE654D95111B01060E1255D575174CEC9C605FFA3E8508D51226F97103B83443D6EB6FA331CBFF2D4A1102DCDFFC10FE8240FCEC9")
	client.Get("https://cloud.189.cn/api/open/file/listFiles.action?noCache=0.7362081385378736&pageSize=60&pageNum=1&mediaType=0&folderId=81419116578555537&iconOption=5&orderBy=lastOpTime&descending=true", &ret)
	fileInfo := FileInfo{}
	if name == "/"+provider.Name {
		fileInfo.isDir = true
	}
	return fileInfo, nil
}
