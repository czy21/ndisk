package _189

import (
	"encoding/json"
	"errors"
	"fmt"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/util"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

type API struct{}

func getJsonAndTokenHeader(req *resty.Request) *resty.Request {
	req.SetHeader("accept", "application/json;charset=UTF-8")
	req = getTokenHeader(req)
	return req
}

func getTokenHeader(req *resty.Request) *resty.Request {
	req.SetHeader("cookie", viper.GetString("cloud189.cookie"))
	return req
}

func (a API) getFolderById(folderId string) (FileListAO, error) {
	var ret FileListAORes
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	params := map[string]string{
		"noCache":    QueryParamNoCache,
		"pageNum":    "1",
		"pageSize":   "100",
		"mediaType":  "0",
		"folderId":   folderId,
		"iconOption": "5",
		"orderBy":    "lastOpTime",
		"descending": "true",
	}
	req.SetQueryParams(params)
	res, err := req.Get(fmt.Sprintf("https://cloud.189.cn/api/open/file/listFiles.action"))
	err = http2.GetClient().JSONUnmarshal(res.Body(), &ret)
	strBody := string(res.Body())
	log.Debugf(strBody)
	if ret.ResMsg != ResSuccessMsg {
		log.Error(strBody)
		err = errors.New(ret.ErrorMsg)
	}
	return ret.FileListAO, err
}

func (a API) CreateFolder(parentFolderId string, name string) (FolderMetaRes, error) {
	var ret FolderMetaRes
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	queryParam := map[string]string{
		"noCache": QueryParamNoCache,
	}
	formData := map[string]string{
		"parentFolderId": parentFolderId,
		"folderName":     name,
	}
	req.SetQueryParams(queryParam)
	req.SetFormData(formData)
	res, err := req.Post("https://cloud.189.cn/api/open/file/createFolder.action")
	err = http2.GetClient().JSONUnmarshal(res.Body(), &ret)
	strBody := string(res.Body())
	log.Debugf(strBody)
	if ret.ResMsg != ResSuccessMsg {
		log.Error(strBody)
		err = errors.New(ret.ErrorMsg)
	}
	return ret, err
}

func (a API) Delete(fileId string, fileName string, isFolder bool) error {
	var (
		err error
		ret TaskRes
	)
	const taskType = "DELETE"
	taskInfosBytes, err := json.Marshal([]map[string]string{
		{
			"fileId":   fileId,
			"fileName": fileName,
			"isFolder": strconv.Itoa(util.BoolToInt(isFolder)),
		},
	})
	queryParam := map[string]string{
		"noCache": QueryParamNoCache,
	}
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	formParam := map[string]string{
		"type":      taskType,
		"taskInfos": string(taskInfosBytes),
	}
	req.SetFormData(formParam)
	req.SetQueryParams(queryParam)
	res, err := req.Post("https://cloud.189.cn/api/open/batch/createBatchTask.action")
	err = http2.GetClient().JSONUnmarshal(res.Body(), &ret)
	strBody := string(res.Body())
	log.Debugf(strBody)
	if ret.ResMsg != ResSuccessMsg {
		log.Error(strBody)
		err = errors.New(ret.ErrorMsg)
		return err
	}
	var taskStatus int
	i := 0
	for {
		if taskStatus == 4 || i >= 2 {
			break
		}
		time.Sleep(500 * time.Millisecond)
		taskStatus = a.CheckTask(ret.TaskId, taskType)
		i++
	}
	return err
}
func (a API) CheckTask(taskId string, kind string) int {
	var (
		err error
		ret TaskRes
	)
	queryParam := map[string]string{
		"noCache": QueryParamNoCache,
	}
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	formParam := map[string]string{
		"taskId": taskId,
		"type":   kind,
	}
	req.SetQueryParams(queryParam)
	req.SetFormData(formParam)
	res, err := req.Post("https://cloud.189.cn/api/open/batch/checkBatchTask.action")
	err = http2.GetClient().JSONUnmarshal(res.Body(), &ret)
	strBody := string(res.Body())
	log.Debugf(strBody)
	if ret.ResMsg != ResSuccessMsg {
		log.Error(strBody)
		err = errors.New(ret.ErrorMsg)
	}
	if err != nil {
		log.Error(err)
	}
	return ret.TaskStatus
}
func (a API) RenameFolder(folderId string, destName string) error {
	var (
		err           error
		folderMetaRes FolderMetaRes
	)
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	formParams := map[string]string{
		"folderId":       folderId,
		"destFolderName": destName,
	}
	req.SetFormData(formParams)
	res, err := req.Post(fmt.Sprintf("https://cloud.189.cn/api/open/file/renameFolder.action?noCache=%s", QueryParamNoCache))
	err = http2.GetClient().JSONUnmarshal(res.Body(), &folderMetaRes)
	strBody := string(res.Body())
	log.Debugf(strBody)
	if folderMetaRes.ResMsg != ResSuccessMsg {
		log.Error(strBody)
		err = errors.New(folderMetaRes.ErrorMsg)
	}
	return err
}

func (a API) DownloadFile(fileId string) error {
	var (
		err                error
		fileDownloadUrlRes FileDownloadUrlRes
	)
	getDownloadUrlParams := map[string]string{
		"noCache": QueryParamNoCache,
		"fileId":  fileId,
	}
	getDownloadUrlReq := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	getDownloadUrlReq.SetQueryParams(getDownloadUrlParams)
	getDownloadUrlReq.SetResult(&fileDownloadUrlRes)
	res, err := getDownloadUrlReq.Get("https://cloud.189.cn/api/open/file/getFileDownloadUrl.action")
	err = http2.GetClient().JSONUnmarshal(res.Body(), &fileDownloadUrlRes)
	strBody := string(res.Body())
	log.Debugf(strBody)
	if fileDownloadUrlRes.ResMsg != ResSuccessMsg {
		log.Error(strBody)
		err = errors.New(fileDownloadUrlRes.ErrorMsg)
	}
	return err
}
