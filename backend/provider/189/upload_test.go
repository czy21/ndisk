package _189

import (
	"fmt"
	"math"
	"testing"
)

func TestMd5Slice(t *testing.T) {
	const sliceItemSize uint64 = 10485760
	var sliceTotal = int64(math.Ceil(float64(1024*10) / float64(sliceItemSize)))
	fmt.Println(sliceTotal)
}
