package _189

import (
	"github.com/czy21/ndisk/model"
)

type Response struct {
	ResCode   int    `json:"res_code"`
	ResMsg    string `json:"res_message"`
	Success   string `json:"success"`
	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

type BaseTrackModel[TID any] struct {
	Id         TID             `json:"id"`
	CreateDate model.LocalTime `json:"createDate"`
	UpdateDate model.LocalTime `json:"lastOpTime"`
}

type FolderMeta struct {
	BaseTrackModel[int64]
	Name string `json:"name"`
}

type FileMeta struct {
	BaseTrackModel[int64]
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type FolderMetaRes struct {
	Response
	FolderMeta
}

type FileListAO struct {
	Count   int          `json:"count"`
	Files   []FileMeta   `json:"fileList"`
	Folders []FolderMeta `json:"folderList"`
}
type FileListAORes struct {
	Response
	FileListAO FileListAO `json:"fileListAO"`
}

type TaskRes struct {
	Response
	TaskId     string `json:"taskId"`
	TaskStatus int    `json:"taskStatus"`
}

type FileDownloadUrlRes struct {
	Response
	Url string `json:"fileDownloadUrl"`
}

type RSAKeyRes struct {
	Response
	Expire  int64  `json:"expire"`
	PKId    string `json:"pkId"`
	PubKey  string `json:"publicKey"`
	Version string `json:"ver"`
}
