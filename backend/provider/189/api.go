package _189

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/czy21/ndisk/exception"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/util"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type API struct{}

func getJsonAndTokenHeader(req *resty.Request) *resty.Request {
	req.SetHeader("accept", "application/json;charset=UTF-8")
	req = getTokenHeader(req)
	return req
}

func getTokenHeader(req *resty.Request) *resty.Request {
	req.SetHeader("cookie", "apm_ua=74B7FDC3A7244BA6FBFD4FC6669EFAFF; apm_ct=20220620210641000; apm_ip=218.81.3.182; apm_uid=9F0A8C1B5282C049DAF6ED52FD27EA97; apm_sid=77FEBF7DACFE57F91D14160020EBDC57; JSESSIONID=476F1741FB09C8B8E6CD0B35348587F5; COOKIE_LOGIN_USER=38851EEE136AA00502CEDE8732B290CD0BC9C911A14DFB6C5DAE526EB5AF8015E2AA10626CF8BF3F88C6E4E342196460994DF26498EDFA1B45CA0C8A6ED85E02A9BE75BF")
	return req
}

func (a API) getFolderById(folderId string) (FileListAO, error) {
	var ret FileListAORes
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	params := map[string]string{
		"noCache":    "0.7362081385378736",
		"pageNum":    "1",
		"pageSize":   "100",
		"mediaType":  "0",
		"folderId":   folderId,
		"iconOption": "5",
		"orderBy":    "lastOpTime",
		"descending": "true",
	}
	req.SetQueryParams(params)
	req.SetResult(&ret)
	res, err := req.Get(fmt.Sprintf("https://cloud.189.cn/api/open/file/listFiles.action"))
	log.Debugf(string(res.Body()))
	if ret.ResMsg != "成功" {
		err = errors.New(ret.ErrorMsg)
	}
	return ret.FileListAO, err
}

func (a API) CreateFolder(parentFolderId string, name string) (FolderMetaRes, error) {
	var ret FolderMetaRes
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	queryParam := map[string]string{
		"noCache": "0.7362081385378736",
	}
	formData := map[string]string{
		"parentFolderId": parentFolderId,
		"folderName":     name,
	}
	req.SetQueryParams(queryParam)
	req.SetFormData(formData)
	req.SetResult(&ret)
	res, err := req.Post("https://cloud.189.cn/api/open/file/createFolder.action")
	log.Debugf(string(res.Body()))
	if ret.ResMsg != "成功" {
		err = errors.New(ret.ErrorMsg)
	}
	return ret, err
}

func (a API) Delete(fileId string, fileName string, isFolder bool) error {
	var (
		err           error
		createTaskRes TaskRes
	)
	const taskType = "DELETE"
	taskInfosBytes, err := json.Marshal([]map[string]string{
		{
			"fileId":   fileId,
			"fileName": fileName,
			"isFolder": strconv.Itoa(util.BoolToInt(isFolder)),
		},
	})
	exception.Check(err)
	createTaskQueryParams := map[string]string{
		"noCache": "0.42321547761726697",
	}
	createTaskReq := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	createTaskFormParams := map[string]string{
		"type":      taskType,
		"taskInfos": string(taskInfosBytes),
	}
	createTaskReq.SetFormData(createTaskFormParams)
	createTaskReq.SetQueryParams(createTaskQueryParams)
	createTaskReq.SetResult(&createTaskRes)
	res, err := createTaskReq.Post("https://cloud.189.cn/api/open/batch/createBatchTask.action")
	log.Debugf(string(res.Body()))
	if createTaskRes.ResMsg != "成功" {
		err = errors.New(createTaskRes.ErrorMsg)
		return err
	}
	commitTaskQueryParam := map[string]string{
		"noCache": "0.42321547761726697",
	}
	commitTaskReq := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	commitTaskFormParams := map[string]string{
		"taskId": createTaskRes.TaskId,
		"type":   taskType,
	}
	commitTaskReq.SetQueryParams(commitTaskQueryParam)
	commitTaskReq.SetFormData(commitTaskFormParams)
	res, err = commitTaskReq.Post("https://cloud.189.cn/api/open/batch/checkBatchTask.action")
	log.Debugf(string(res.Body()))
	if createTaskRes.ResMsg != "成功" {
		err = errors.New(createTaskRes.ErrorMsg)
	}
	return err
}

func (a API) RenameFolder(folderId string, destName string) error {
	var (
		err        error
		folderMeta FolderMetaRes
	)
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	formParams := map[string]string{
		"folderId":       folderId,
		"destFolderName": destName,
	}
	req.SetFormData(formParams)
	req.SetResult(&folderMeta)
	res, err := req.Post(fmt.Sprintf("https://cloud.189.cn/api/open/file/renameFolder.action?noCache=%s", "0.33825729434675056"))
	log.Debugf(string(res.Body()))
	if folderMeta.ResMsg != "成功" {
		err = errors.New(folderMeta.ErrorMsg)
	}
	return err
}

func (a API) DownloadFile(fileId string) error {
	var (
		err                error
		fileDownloadUrlRes FileDownloadUrlRes
	)
	getDownloadUrlParams := map[string]string{
		"noCache": "0.42321547761726697",
		"fileId":  fileId,
	}
	getDownloadUrlReq := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	getDownloadUrlReq.SetQueryParams(getDownloadUrlParams)
	getDownloadUrlReq.SetResult(&fileDownloadUrlRes)
	res, err := getDownloadUrlReq.Get("https://cloud.189.cn/api/open/file/getFileDownloadUrl.action")
	log.Debugf(string(res.Body()))
	if fileDownloadUrlRes.ResMsg != "成功" {
		err = errors.New(fileDownloadUrlRes.ErrorMsg)
	}
	return err
}
