package util

import "io"

func CopyBuffer(dst io.Writer, src io.Reader, bf int) (n int64, err error) {
	if bf < 8 || bf > 64 {
		bf = 8
	}
	return io.CopyBuffer(dst, src, make([]byte, 1024*1024*bf))
}
