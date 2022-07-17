package _189

import (
	"errors"
	"fmt"
	http2 "github.com/czy21/ndisk/http"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type API struct{}

func (a API) getFolderById(folderId string) (FileListAO, error) {
	var ret FileListAORes
	req := getHeader(http2.GetClient().NewRequest())
	req.SetResult(&ret)
	res, err := req.Get(fmt.Sprintf("https://cloud.189.cn/api/open/file/listFiles.action?noCache=0.7362081385378736&pageSize=60&pageNum=1&mediaType=0&folderId=%s&iconOption=5&orderBy=lastOpTime&descending=true", folderId))
	log.Debugf(string(res.Body()))
	if ret.ResMsg != "成功" {
		err = errors.New(ret.ErrorMsg)
	}
	return ret.FileListAO, err
}

func (a API) CreateFolder(parentFolderId string, name string) (FolderMetaAddRes, error) {
	var ret FolderMetaAddRes
	req := getHeader(http2.GetClient().NewRequest())
	reqBody := make(map[string]string)
	reqBody["parentFolderId"] = parentFolderId
	reqBody["folderName"] = name
	req.SetFormData(reqBody)
	req.SetResult(&ret)
	res, err := req.Post("https://cloud.189.cn/api/open/file/createFolder.action?noCache=0.40908345822366")
	log.Debugf(string(res.Body()))
	if ret.ResMsg != "成功" {
		err = errors.New(ret.ErrorMsg)
	}
	return ret, err
}

func getHeader(req *resty.Request) *resty.Request {
	req.SetHeader("accept", "application/json;charset=UTF-8")
	req.SetHeader("cookie", "apm_ua=74B7FDC3A7244BA6FBFD4FC6669EFAFF; apm_ct=20220620210641000; apm_ip=218.81.3.182; apm_uid=9F0A8C1B5282C049DAF6ED52FD27EA97; apm_sid=77FEBF7DACFE57F91D14160020EBDC57; JSESSIONID=476F1741FB09C8B8E6CD0B35348587F5; COOKIE_LOGIN_USER=38851EEE136AA00502CEDE8732B290CD0BC9C911A14DFB6C5DAE526EB5AF8015E2AA10626CF8BF3F88C6E4E342196460994DF26498EDFA1B45CA0C8A6ED85E02A9BE75BF")
	return req
}
