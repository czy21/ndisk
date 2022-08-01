package constant

const (
	DefaultPageIndex = 1
	DefaultPageSize  = 10
)

var WebDavMethods = [...]string{
	"PROPFIND",
	"MKCOL",
	"LOCK",
	"UNLOCK",
	"PROPPATCH",
	"COPY",
	"MOVE",
}

const HttpExtra = "extra"
const HttpExtraMethod = "method"
const HttpExtraFileSize = "fileSize"
const HttpExtraChunks = "chunks"
const HttpExtraChunkL = "chunkL"
const HttpExtraChunkI = "chunkI"
const HttpExtraRangeS = "rangeS"
const HttpExtraRangeE = "rangeE"
const HttpExtraDownloadUrl = "downloadUrl"
