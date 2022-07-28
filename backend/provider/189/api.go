package _189

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/czy21/ndisk/cache"
	http2 "github.com/czy21/ndisk/http"
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

type API struct{}

func getRequestWithJsonAndToken(req *resty.Request) *resty.Request {
	req.SetHeader("accept", "application/json;charset=UTF-8")
	setTokenHeader(req)
	return req
}

func logRes(funcName string, strBody string, ret ResponseVO, err error) {
	fmtMsg := fmt.Sprintf("%s %s", funcName, strBody)
	if ret.ResMsg != SuccessMsg {
		log.Error(fmtMsg)
		err = errors.New(ret.ErrorMsg)
	}
}

func setTokenHeader(req *resty.Request) {
	req.SetHeader("cookie", viper.GetString("cloud189.cookie"))
}

func (a API) GetFolderById(folderId string) (FileListAO, error) {
	var (
		ret FileListAORes
		err error
	)
	pageIndex := 1
	for {
		params := map[string]string{
			"noCache":    QueryParamNoCache,
			"pageNum":    strconv.Itoa(pageIndex),
			"pageSize":   "100",
			"mediaType":  "0",
			"folderId":   folderId,
			"iconOption": "5",
			"orderBy":    "lastOpTime",
			"descending": "true",
		}
		var pageRet FileListAORes
		req := getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
			SetQueryParams(params).
			SetResult(&pageRet)
		res, err := req.Get(fmt.Sprintf("%s/open/file/listFiles.action", ApiUrl))
		if pageRet.FileListAO.Count == 0 {
			break
		}
		ret.FileListAO.Folders = append(ret.FileListAO.Folders, pageRet.FileListAO.Folders...)
		ret.FileListAO.Files = append(ret.FileListAO.Files, pageRet.FileListAO.Files...)
		logRes("GetFolderById", res.String(), pageRet.ResponseVO, err)
		pageIndex++
	}
	return ret.FileListAO, err
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
	req := getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetQueryParams(queryParam).
		SetFormData(formData).
		SetResult(&ret)
	res, err := req.Post(fmt.Sprintf("%s/open/file/createFolder.action", ApiUrl))
	logRes("CreateFolder", res.String(), ret.ResponseVO, err)
	return err
}

func (a API) Delete(fileId string, fileName string, isFolder bool) (err error) {
	var (
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
	formParam := map[string]string{
		"type":      taskType,
		"taskInfos": string(taskInfosBytes),
	}
	req := getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetFormData(formParam).
		SetQueryParams(queryParam).
		SetResult(&ret)
	res, err := req.Post(fmt.Sprintf("%s/open/batch/createBatchTask.action", ApiUrl))
	logRes("Delete", res.String(), ret.ResponseVO, err)
	if err != nil {
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
	formParam := map[string]string{
		"taskId": taskId,
		"type":   kind,
	}
	req := getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetQueryParams(queryParam).
		SetFormData(formParam).
		SetResult(&ret)
	res, err := req.Post(fmt.Sprintf("%s/open/batch/checkBatchTask.action", ApiUrl))
	logRes("CheckTask", res.String(), ret.ResponseVO, err)
	return ret.TaskStatus
}
func (a API) RenameFolder(folderId string, destName string) (err error) {
	var (
		ret FolderRes
	)
	formParams := map[string]string{
		"folderId":       folderId,
		"destFolderName": destName,
	}
	req := getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetFormData(formParams).
		SetResult(&ret)
	res, err := req.Post(fmt.Sprintf("%s/open/file/renameFolder.action?noCache=%s", ApiUrl, QueryParamNoCache))
	logRes("RenameFolder", res.String(), ret.ResponseVO, err)
	return err
}

func (a API) GetFileInfoById(fileId string) (FileInfoVO, error) {
	var (
		err error
		ret FileInfoVORes
	)
	queryParam := map[string]string{
		"noCache": QueryParamNoCache,
		"fileId":  fileId,
	}
	req := getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetQueryParams(queryParam).
		SetResult(&ret)
	res, err := req.Get(fmt.Sprintf("%s/open/file/getFileInfo.action", ApiUrl))
	logRes("GetFileInfoById", res.String(), ret.ResponseVO, err)
	if ret.ResCode == ResFileNotFoundCode {
		err = fs.ErrNotExist
	}
	return ret.FileInfoVO, err
}

func (a API) GetRSAKey() (RSAKeyRes, error) {
	const rsaCacheKey = "e:rsa:189"
	var (
		err error
		ret RSAKeyRes
	)
	req := getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetResult(&ret)
	cache.Client.GetObj(context.Background(), rsaCacheKey, &ret)
	if ret.PKId == "" {
		res, err := req.Get(fmt.Sprintf("%s/security/generateRsaKey.action?noCache=%s", ApiUrl, QueryParamNoCache))
		logRes("GetRSAKey", res.String(), ret.ResponseVO, err)
		expire := time.Time(ret.Expire).Sub(time.Now())
		cache.Client.SetObjEX(context.Background(), rsaCacheKey, ret, expire)
	}
	return ret, err
}

func (a API) GetUserBriefInfo() UserBriefInfoVO {
	var (
		err error
		ret UserBriefInfoVORes
	)
	req := getRequestWithJsonAndToken(http2.GetClient().NewRequest()).
		SetResult(&ret)
	res, err := req.Get(fmt.Sprintf("%s/portal/v2/getUserBriefInfo.action?noCache=%s", ApiUrl, QueryParamNoCache))
	logRes("GetUserBriefInfo", res.String(), ret.ResponseVO, err)
	return ret.UserBriefInfoVO
}

func (a API) UploadRequest(uri string, queryParam map[string]string, resVO interface{}, errPredicate func() bool) error {
	var err error
	rand.Seed(time.Now().UnixNano())
	c := strconv.FormatInt(time.Now().UnixMilli(), 10)
	r := random("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx")
	l := random("xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx")
	l = l[0 : 16+int(16*rand.Float32())|0]
	var u []string
	for k, v := range queryParam {
		u = append(u, k+"="+v)
	}
	sessionKey := a.GetUserBriefInfo().SessionKey
	encryptParam := aesEncrypt([]byte(strings.Join(u, "&")), []byte(l[0:16]))
	signature := hmacSha1(fmt.Sprintf("SessionKey=%s&Operate=GET&RequestURI=%s&Date=%s&params=%s", sessionKey, uri, c, encryptParam), l)
	rsaRes, err := a.GetRSAKey()
	if err != nil {
		return err
	}
	b := rsaEncode([]byte(l), rsaRes.PubKey)
	req := http2.GetClient().NewRequest().
		SetHeader("accept", "application/json;charset=UTF-8").
		SetHeader("SessionKey", sessionKey).
		SetHeader("Signature", signature).
		SetHeader("X-Request-Date", c).
		SetHeader("X-Request-ID", r).
		SetHeader("EncryptionText", b).
		SetHeader("PkId", rsaRes.PKId).
		SetQueryParam("params", encryptParam).
		SetResult(resVO)
	res, err := req.Get("https://upload.cloud.189.cn" + uri)
	if errPredicate() {
		log.Error(res)
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
	if len(b) == 0 {
		return nil
	}
	var uploadUrlsRes UploadUrlVORes
	err := API{}.UploadRequest("/person/getMultiUploadUrls",
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
