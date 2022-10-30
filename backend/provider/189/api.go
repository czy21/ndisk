package _189

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/czy21/ndisk/cache"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/util"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"io/fs"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type API struct {
	File model.ProviderFile
}

func (a API) getRequestWithJsonAndToken(req *resty.Request) *resty.Request {
	req.SetHeader("accept", "application/json;charset=UTF-8")
	a.setTokenHeader(req)
	return req
}

func (a API) checkError(funcName string, res *resty.Response, ret ResponseVO) (err error) {
	fmtMsg := fmt.Sprintf("%s %s", funcName, res.String())
	if ret.ResMsg != SuccessMsg {
		err = errors.New(fmtMsg)
	}
	return err
}

func (a API) setTokenHeader(req *resty.Request) {
	req.SetHeader("cookie", viper.GetString("cloud189.cookie"))
}

func (a API) GetObjectsById(folderId string, fileName string) ([]*FileVO, error) {
	var (
		files []*FileVO
		err   error
	)
	pageIndex := 1
	pageSize := 60
	for {
		params := map[string]string{
			"noCache":    QueryParamNoCache,
			"pageNum":    strconv.Itoa(pageIndex),
			"pageSize":   strconv.Itoa(pageSize),
			"mediaType":  "0",
			"folderId":   folderId,
			"iconOption": "5",
			"orderBy":    "lastOpTime",
			"descending": "true",
		}
		if fileName != "" {
			params["filename"] = fileName
		}
		var pageRet FileListAORes
		req := a.getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
			SetQueryParams(params).
			SetResult(&pageRet)
		res, _ := req.Get(fmt.Sprintf("%s/open/file/searchFiles.action", ApiUrl))
		err = a.checkError("GetObjectsById", res, pageRet.ResponseVO)
		if err != nil || pageRet.Count == 0 {
			break
		}
		for _, t := range pageRet.Folders {
			t.IsDir = true
		}
		for _, t := range pageRet.Files {
			t.IsDir = false
		}
		files = append(append(files, pageRet.Folders...), pageRet.Files...)
		pageIndex++
	}
	return files, err
}

func (a API) GetFoldersById(folderId string) ([]FolderNodeRes, error) {
	var ret []FolderNodeRes
	queryParam := map[string]string{
		"noCache": QueryParamNoCache,
	}
	formParam := map[string]string{
		"id":      folderId,
		"orderBy": strconv.Itoa(1),
		"order":   "ASC",
	}
	req := a.getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetFormData(formParam).
		SetQueryParams(queryParam).
		SetResult(&ret)
	_, err := req.Post(fmt.Sprintf("%s/portal/getObjectFolderNodes.action", ApiUrl))
	return ret, err
}
func (a API) GetFileById(fileId string) (FileInfoVO, error) {
	var (
		ret FileInfoVORes
		err error
	)
	queryParam := map[string]string{
		"noCache": QueryParamNoCache,
		"fileId":  fileId,
	}
	req := a.getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetQueryParams(queryParam).
		SetResult(&ret)
	res, _ := req.Get(fmt.Sprintf("%s/open/file/getFileInfo.action", ApiUrl))
	err = a.checkError("GetFileById", res, ret.ResponseVO)
	if ret.ResCode == ResFileNotFoundCode {
		err = fs.ErrNotExist
	}
	return ret.FileInfoVO, err
}
func (a API) CreateFolder(parentFolderId string, name string) (err error) {
	var ret FolderRes
	queryParam := map[string]string{
		"noCache": QueryParamNoCache,
	}
	formData := map[string]string{
		"parentFolderId": parentFolderId,
		"folderName":     name,
	}
	req := a.getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetQueryParams(queryParam).
		SetFormData(formData).
		SetResult(&ret)
	res, _ := req.Post(fmt.Sprintf("%s/open/file/createFolder.action", ApiUrl))
	err = a.checkError("CreateFolder", res, ret.ResponseVO)
	return err
}
func (a API) Delete(fileId string, fileName string, isFolder bool) (err error) {
	return a.CreateTask("DELETE", fileId, fileName, isFolder, nil)
}
func (a API) Copy(fileId string, fileName string, isFolder bool, targetFolderId string) (err error) {
	return a.CreateTask("COPY", fileId, fileName, isFolder, map[string]string{"targetFolderId": targetFolderId})
}
func (a API) Move(fileId string, fileName string, isFolder bool, targetFolderId string) (err error) {
	return a.CreateTask("MOVE", fileId, fileName, isFolder, map[string]string{"targetFolderId": targetFolderId})
}
func (a API) CreateTask(kind string, fileId string, fileName string, isFolder bool, extraFormParam map[string]string) error {
	var (
		ret TaskRes
		err error
	)
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
	formParam := map[string]string{
		"type":      kind,
		"taskInfos": string(taskInfosBytes),
	}
	if extraFormParam != nil {
		for k, v := range extraFormParam {
			formParam[k] = v
		}
	}
	req := a.getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetFormData(formParam).
		SetQueryParams(queryParam).
		SetResult(&ret)
	res, _ := req.Post(fmt.Sprintf("%s/open/batch/createBatchTask.action", ApiUrl))
	err = a.checkError(fmt.Sprintf("CreateTask %s", kind), res, ret.ResponseVO)
	return err
}
func (a API) CheckTask(taskId string, kind string) (s int, err error) {
	var (
		ret TaskRes
	)
	queryParam := map[string]string{
		"noCache": QueryParamNoCache,
	}
	formParam := map[string]string{
		"taskId": taskId,
		"type":   kind,
	}
	req := a.getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetQueryParams(queryParam).
		SetFormData(formParam).
		SetResult(&ret)
	res, _ := req.Post(fmt.Sprintf("%s/open/batch/checkBatchTask.action", ApiUrl))
	err = a.checkError("CheckTask", res, ret.ResponseVO)
	return ret.TaskStatus, err
}
func (a API) RenameFile(fileId string, destName string) (err error) {
	var ret FileRes
	formParam := map[string]string{
		"fileId":       fileId,
		"destFileName": destName,
	}
	req := a.getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetFormData(formParam).
		SetResult(&ret)
	res, _ := req.Post(fmt.Sprintf("%s/open/file/renameFile.action?noCache=%s", ApiUrl, QueryParamNoCache))
	err = a.checkError("RenameFile", res, ret.ResponseVO)
	return err
}
func (a API) RenameFolder(folderId string, destName string) (err error) {
	var ret FolderRes
	formParams := map[string]string{
		"folderId":       folderId,
		"destFolderName": destName,
	}
	req := a.getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetFormData(formParams).
		SetResult(&ret)
	res, _ := req.Post(fmt.Sprintf("%s/open/file/renameFolder.action?noCache=%s", ApiUrl, QueryParamNoCache))
	err = a.checkError("RenameFolder", res, ret.ResponseVO)
	return err
}
func (a API) GetRSAKey() (ret RSAKeyRes, err error) {
	const rsaCacheKey = "e:rsa:189"
	req := a.getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetResult(&ret)
	cache.Client.GetObj(context.Background(), rsaCacheKey, &ret)
	if ret.PKId == "" {
		res, _ := req.Get(fmt.Sprintf("%s/security/generateRsaKey.action?noCache=%s", ApiUrl, QueryParamNoCache))
		err = a.checkError("GetRSAKey", res, ret.ResponseVO)
		if err != nil {
			return ret, err
		}
		expire := time.Time(ret.Expire).Sub(time.Now())
		cache.Client.SetObjEX(context.Background(), rsaCacheKey, ret, expire)
	}
	return ret, err
}
func (a API) GetUserBrief() (UserBriefInfoVO, error) {
	var (
		err error
		ret UserBriefInfoVORes
	)
	req := a.getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetResult(&ret)
	res, _ := req.Get(fmt.Sprintf("%s/portal/v2/getUserBriefInfo.action?noCache=%s", ApiUrl, QueryParamNoCache))
	err = a.checkError("GetUserBrief", res, ret.ResponseVO)
	return ret.UserBriefInfoVO, err
}
func (a API) UploadRequest(uri string, queryParam map[string]string, resVO interface{}, errPredicate func() bool) (err error) {
	rand.Seed(time.Now().UnixNano())
	c := strconv.FormatInt(time.Now().UnixMilli(), 10)
	r := random("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx")
	l := random("xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx")
	l = l[0 : 16+int(16*rand.Float32())|0]
	var u []string
	for k, v := range queryParam {
		u = append(u, k+"="+v)
	}
	userBriefInfo, err := a.GetUserBrief()
	if err != nil {
		return err
	}
	encryptParam := aesEncrypt([]byte(strings.Join(u, "&")), []byte(l[0:16]))
	signature := hmacSha1(fmt.Sprintf("SessionKey=%s&Operate=GET&RequestURI=%s&Date=%s&params=%s", userBriefInfo.SessionKey, uri, c, encryptParam), l)
	rsaRes, err := a.GetRSAKey()
	if err != nil {
		return err
	}
	b := rsaEncode([]byte(l), rsaRes.PubKey)
	req := http2.GetClient().NewRequest().
		SetHeader("accept", "application/json;charset=UTF-8").
		SetHeader("SessionKey", userBriefInfo.SessionKey).
		SetHeader("Signature", signature).
		SetHeader("X-Request-Date", c).
		SetHeader("X-Request-ID", r).
		SetHeader("EncryptionText", b).
		SetHeader("PkId", rsaRes.PKId).
		SetQueryParam("params", encryptParam).
		SetResult(resVO)
	res, err := req.Get("https://upload.cloud.189.cn" + uri)
	if errPredicate() {
		err = errors.New(res.String())
	}
	return err
}
func (a API) CreateFile(parentFolderId, fileName string, fileSize int64, fileMd5 string) (InitUploadVO, error) {
	var initUploadVO ResponseDataVO[InitUploadVO, any]
	lazyCheck := 1
	if fileSize == 0 {
		lazyCheck = 0
	}
	queryParam := map[string]string{
		"parentFolderId": parentFolderId,
		"fileName":       url.QueryEscape(fileName),
		"fileSize":       strconv.FormatInt(fileSize, 10),
		"sliceSize":      strconv.FormatInt(1024*1024*10, 10),
		"lazyCheck":      fmt.Sprintf("%d", lazyCheck),
	}
	if fileMd5 != "" {
		queryParam["fileMd5"] = fileMd5
		queryParam["sliceMd5"] = fileMd5
	}
	err := a.UploadRequest("/person/initMultiUpload", queryParam, &initUploadVO, func() bool {
		return initUploadVO.Code != SuccessCode
	})
	return initUploadVO.Data, err
}
func (a API) CommitFile(fileId string, fileSize int64, fileMd5 string, sliceMd5 string) (err error) {
	var ret ResponseDataVO[any, CommitFileVO]
	lazyCheck := 1
	if fileSize == 0 {
		lazyCheck = 0
	}
	err = a.UploadRequest(
		"/person/commitMultiUploadFile",
		map[string]string{
			"uploadFileId": fileId,
			"fileMd5":      fileMd5,
			"sliceMd5":     sliceMd5,
			"lazyCheck":    fmt.Sprintf("%d", lazyCheck),
			"opertype":     "3",
		}, &ret, func() bool {
			return ret.Code != SuccessCode
		})
	return err
}
func (a API) UploadChunk(fileId string, b []byte, md5Bytes []byte, index int) error {
	md5Base64 := base64.StdEncoding.EncodeToString(md5Bytes)
	var uploadUrlsRes UploadUrlVORes
	err := a.UploadRequest("/person/getMultiUploadUrls",
		map[string]string{
			"partInfo":     fmt.Sprintf("%d-%s", index, md5Base64),
			"uploadFileId": fileId,
		}, &uploadUrlsRes, func() bool {
			return uploadUrlsRes.Code != SuccessCode
		})
	if err != nil {
		return err
	}
	uploadData := uploadUrlsRes.UploadUrls[fmt.Sprintf("partNumber_%d", index)]
	uploadHeader, _ := url.PathUnescape(uploadData.RequestHeader)
	uploadHeaders := strings.Split(uploadHeader, "&")
	uploadRequest := http2.GetClient().NewRequest().SetBody(bytes.NewReader(b))
	for _, t := range uploadHeaders {
		i := strings.Index(t, "=")
		uploadRequest.Header.Set(t[0:i], t[i+1:])
	}
	_, err = uploadRequest.Put(uploadData.RequestURL)
	log.Debugf("fileId: %s request: %s", fileId, uploadData)
	return err
}
