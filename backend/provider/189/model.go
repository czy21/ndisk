package _189

type Response struct {
	Code    int    `json:"res_code"`
	Message string `json:"res_message"`
}

type FolderMeta struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type FileMeta struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
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
