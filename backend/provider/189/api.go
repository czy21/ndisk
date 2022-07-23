package _189

import (
	"bytes"
	"crypto/aes"
	"crypto/hmac"
	rand2 "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/util"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var b64map = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
var BI_RM = "0123456789abcdefghijklmnopqrstuvwxyz"

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
	req.SetResult(&ret)
	res, err := req.Get(fmt.Sprintf("%s/open/file/listFiles.action", ApiUrl))
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
	req.SetResult(&ret)
	res, err := req.Post(fmt.Sprintf("%s/open/file/createFolder.action", ApiUrl))
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
	req.SetResult(&ret)
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
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	formParam := map[string]string{
		"taskId": taskId,
		"type":   kind,
	}
	req.SetQueryParams(queryParam)
	req.SetFormData(formParam)
	req.SetResult(&ret)
	res, err := req.Post(fmt.Sprintf("%s/open/batch/checkBatchTask.action", ApiUrl))
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
	req.SetResult(&ret)
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
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	req.SetQueryParams(queryParam)
	req.SetResult(&ret)
	res, err := req.Get(fmt.Sprintf("%s/open/file/getFileInfo.action", ApiUrl))
	logRes("GetFileInfoById", res.String(), ret.ResponseVO, err)
	return ret.FileInfoVO, err
}

func (a API) GetRSAKey() (RSAKeyRes, error) {
	var (
		err error
		ret RSAKeyRes
	)
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	req.SetResult(&ret)
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
func int2char(a int) string {
	return strings.Split(BI_RM, "")[a]
}
func b64tohex(a string) string {
	d := ""
	e := 0
	c := 0
	for i := 0; i < len(a); i++ {
		m := strings.Split(a, "")[i]
		if m != "=" {
			v := strings.Index(b64map, m)
			if 0 == e {
				e = 1
				d += int2char(v >> 2)
				c = 3 & v
			} else if 1 == e {
				e = 2
				d += int2char(c<<2 | v>>4)
				c = 15 & v
			} else if 2 == e {
				e = 3
				d += int2char(c)
				d += int2char(v >> 2)
				c = 3 & v
			} else {
				e = 0
				d += int2char(c<<2 | v>>4)
				d += int2char(15 & v)
			}
		}
	}
	if e == 1 {
		d += int2char(c << 2)
	}
	return d
}
func rsaEncode(origData []byte, j_rsakey string, hex bool) string {
	publicKey := []byte("-----BEGIN PUBLIC KEY-----\n" + j_rsakey + "\n-----END PUBLIC KEY-----")
	block, _ := pem.Decode(publicKey)
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	pub := pubInterface.(*rsa.PublicKey)
	b, err := rsa.EncryptPKCS1v15(rand2.Reader, pub, origData)
	if err != nil {
		log.Errorf("err: %s", err.Error())
	}
	res := base64.StdEncoding.EncodeToString(b)
	if hex {
		return b64tohex(res)
	}
	return res
}

func (a API) GetUserBriefInfo() UserBriefInfoVO {
	var (
		err error
		ret UserBriefInfoVORes
	)
	req := getJsonAndTokenHeader(http2.GetClient().NewRequest())
	req.SetResult(&ret)
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
	rsaRes, err := a.GetRSAKey()
	if err != nil {
		return ret.Data, err
	}
	b := rsaEncode([]byte(l), rsaRes.PubKey, false)
	req := http2.GetClient().NewRequest().
		SetHeader("accept", "application/json;charset=UTF-8").
		SetHeader("SessionKey", sessionKey).
		SetHeader("Signature", signature).
		SetHeader("X-Request-Date", c).
		SetHeader("X-Request-ID", r).
		SetHeader("EncryptionText", b).
		SetHeader("PkId", rsaRes.PKId)
	req.SetQueryParam("params", encryptParam)
	req.SetResult(&ret)
	res, err := req.Get("https://upload.cloud.189.cn" + uri)
	log.Debug(res.String())
	return ret.Data, err
}

func (a API) CreateUpload(parentFolderId, fileName string, fileSize, sliceSize int64) (InitUploadVO, error) {
	res, err := a.UploadRequest(
		"/person/initMultiUpload",
		map[string]string{
			"parentFolderId": parentFolderId,
			"fileName":       url.QueryEscape(fileName),
			"fileSize":       strconv.FormatInt(fileSize, 10),
			"sliceSize":      strconv.FormatInt(sliceSize, 10),
			"lazyCheck":      "1",
		})
	return res, err
}
