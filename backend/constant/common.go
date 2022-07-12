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
