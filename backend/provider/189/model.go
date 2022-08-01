package _189

import (
	"github.com/czy21/ndisk/model"
)

type ResponseVO struct {
	ResCode   any    `json:"res_code"`
	ResMsg    string `json:"res_message"`
	Success   string `json:"success"`
	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

type BaseTrackModel[TID any, TCreateDate any, TUpdateDate any] struct {
	Id         TID         `json:"id"`
	Name       string      `json:"name"`
	CreateDate TCreateDate `json:"createDate"`
	UpdateDate TUpdateDate `json:"lastOpTime"`
}

type FolderVO struct {
	BaseTrackModel[int64, model.LocalTime, model.LocalTime]
}

type FileVO struct {
	BaseTrackModel[int64, model.LocalTime, model.LocalTime]
	Size int64 `json:"size"`
}

type FolderRes struct {
	ResponseVO
	FolderVO
}

type FileRes struct {
	ResponseVO
	FileVO
}
type FileListAO struct {
	Count   int        `json:"count"`
	Files   []FileVO   `json:"fileList"`
	Folders []FolderVO `json:"folderList"`
}
type FileListAORes struct {
	ResponseVO
	FileListAO FileListAO `json:"fileListAO"`
}

type TaskRes struct {
	ResponseVO
	TaskId     string `json:"taskId"`
	TaskStatus int    `json:"taskStatus"`
}

type FileInfoVO struct {
	BaseTrackModel[int64, model.LocalTime, model.UnixTime]
	MediaType       int    `json:"mediaType"`
	FileDownloadUrl string `json:"fileDownloadUrl"`
	Size            int64  `json:"size"`
}

type FileInfoVORes struct {
	ResponseVO
	FileInfoVO
}

type RSAKeyRes struct {
	ResponseVO
	Expire  model.UnixTime `json:"expire"`
	PKId    string         `json:"pkId"`
	PubKey  string         `json:"pubKey"`
	Version string         `json:"ver"`
}

type UserBriefInfoVO struct {
	EncryptAccount string `json:"encryptAccount"`
	UserAccount    string `json:"userAccount"`
	SessionKey     string `json:"sessionKey"`
}

type UserBriefInfoVORes struct {
	ResponseVO
	UserBriefInfoVO
}

type InitUploadVO struct {
	FileDataExists int    `json:"fileDataExists"`
	UploadFileId   string `json:"uploadFileId"`
	UploadHost     string `json:"uploadHost"`
	UploadType     int    `json:"uploadType"`
}

type ResponseDataVO[TData any, TFile any] struct {
	Code string `json:"code"`
	Data TData  `json:"data"`
	File TFile  `json:"file"`
}

type CommitFileVO struct {
	Name string `json:"fileName"`
}

type UploadPartVO struct {
	RequestURL    string `json:"requestURL"`
	RequestHeader string `json:"requestHeader"`
}

type UploadUrlVORes struct {
	Code       string                  `json:"code"`
	UploadUrls map[string]UploadPartVO `json:"uploadUrls"`
}
