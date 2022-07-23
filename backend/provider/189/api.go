package _189

import (
	"bytes"
	"crypto/aes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/util"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type API struct{}

func getJsonAndTokenHeader(req *resty.Request) *resty.Request {
	req.SetHeader("accept", "application/json;charset=UTF-8")
	req = getTokenHeader(req)
	return req
}

func logRes(funcName string, strBody string, ret ResponseVO, err error) {
	fmtMsg := fmt.Sprintf("%s %s", funcName, strBody)
	log.Debugf(fmtMsg)
	if ret.ResMsg != ResSuccessMsg {
		log.Error(fmtMsg)
		err = errors.New(ret.ErrorMsg)
	}
}

func getTokenHeader(req *resty.Request) *resty.Request {
	req.SetHeader("cookie", viper.GetString("cloud189.cookie"))
	return req
}

func (a API) GetFolderById(folderId string) (FileListAO, error) {
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
	res, err := req.Get(fmt.Sprintf("%s/open/file/listFiles.action", ApiUrl))
	err = http2.GetClient().JSONUnmarshal(res.Body(), &ret)
	logRes("GetFolderById", res.String(), ret.ResponseVO, err)
	return ret.FileListAO, err
}

func (a API) CreateFolder(parentFolderId string, name string) (FolderRes, error) {
	var ret FolderRes
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
	res, err := req.Post(fmt.Sprintf("%s/open/file/createFolder.action", ApiUrl))
	err = http2.GetClient().JSONUnmarshal(res.Body(), &ret)
	logRes("CreateFolder", res.String(), ret.ResponseVO, err)
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
	res, err := req.Post(fmt.Sprintf("%s/open/batch/createBatchTask.action", ApiUrl))
	err = http2.GetClient().JSONUnmarshal(res.Body(), &ret)
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
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	formParam := map[string]string{
		"taskId": taskId,
		"type":   kind,
	}
	req.SetQueryParams(queryParam)
	req.SetFormData(formParam)
	res, err := req.Post(fmt.Sprintf("%s/open/batch/checkBatchTask.action", ApiUrl))
	err = http2.GetClient().JSONUnmarshal(res.Body(), &ret)
	logRes("CheckTask", res.String(), ret.ResponseVO, err)
	return ret.TaskStatus
}
func (a API) RenameFolder(folderId string, destName string) error {
	var (
		err error
		ret FolderRes
	)
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	formParams := map[string]string{
		"folderId":       folderId,
		"destFolderName": destName,
	}
	req.SetFormData(formParams)
	res, err := req.Post(fmt.Sprintf("%s/open/file/renameFolder.action?noCache=%s", ApiUrl, QueryParamNoCache))
	err = http2.GetClient().JSONUnmarshal(res.Body(), &ret)
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
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	req.SetQueryParams(queryParam)
	req.SetResult(&ret)
	res, err := req.Get(fmt.Sprintf("%s/open/file/getFileInfo.action", ApiUrl))
	err = http2.GetClient().JSONUnmarshal(res.Body(), &ret)
	logRes("GetFileInfoById", res.String(), ret.ResponseVO, err)
	return ret.FileInfoVO, err
}

func (a API) GetRSAKey() (RSAKeyRes, error) {
	var (
		err error
		ret RSAKeyRes
	)
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	res, err := req.Get(fmt.Sprintf("%s/security/generateRsaKey.action?noCache=%s", ApiUrl, QueryParamNoCache))
	logRes("GetRSAKey", res.String(), ret.ResponseVO, err)
	return ret, err
}

func hmacSha1(data string, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
func aesEncrypt(data, key []byte) string {
	block, _ := aes.NewCipher(key)
	if block == nil {
		return hex.EncodeToString([]byte{})
	}
	data = pkcs7Padding(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}
	return hex.EncodeToString(decrypted)
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func random(v string) string {
	reg := regexp.MustCompilePOSIX("[xy]")
	data := reg.ReplaceAllFunc([]byte(v), func(msg []byte) []byte {
		var i int64
		t := int64(16*rand.Float32()) | 0
		if msg[0] == 120 {
			i = t
		} else {
			i = 3&t | 8
		}
		return []byte(strconv.FormatInt(i, 16))
	})
	return string(data)
}

func (a API) GetUserBriefInfo() UserBriefInfoVO {
	var (
		err error
		ret UserBriefInfoVORes
	)
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	res, err := req.Get(fmt.Sprintf("%s/portal/v2/getUserBriefInfo.action?noCache=%s", ApiUrl, QueryParamNoCache))
	logRes("GetUserBriefInfo", res.String(), ret.ResponseVO, err)
	return ret.UserBriefInfoVO
}

func (a API) UploadRequest(uri string, queryParam map[string]string) (InitUploadVO, error) {
	var (
		err error
		ret ResponseDataVO[InitUploadVO]
	)

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
	rsa, err := a.GetRSAKey()
	if err != nil {
		return ret.Data, err
	}
	req := http2.GetClient().NewRequest().
		SetHeader("accept", "application/json;charset=UTF-8").
		SetHeader("SessionKey", sessionKey).
		SetHeader("Signature", signature).
		SetHeader("X-Request-Date", c).
		SetHeader("X-Request-ID", r).
		SetHeader("EncryptionText", "").
		SetHeader("PkId", rsa.PKId)
	req.SetQueryParam("param", encryptParam)
	req.SetResult(&ret)
	_, err = req.Get("https://upload.cloud.189.cn" + uri)
	return ret.Data, err
}

func (a API) Upload(parentFolderId string, fileName string, fileSize, bytes []byte) {
	const chunkSize uint64 = 10485760
	slices := math.Max(1, math.Ceil(float64(len(bytes))/float64(chunkSize)))
	fmt.Println(slices)
}
