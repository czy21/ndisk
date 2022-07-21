package webdav

import (
	"io"
	"net/http"
)

type Writer struct {
	http.ResponseWriter
}

func (w Writer) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}
