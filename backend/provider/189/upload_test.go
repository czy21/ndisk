package _189

import (
	"fmt"
	"math"
	"testing"
)

func TestMd5Slice(t *testing.T) {
	const sliceItemSize uint64 = 10485760
	var sliceTotal = math.Ceil(float64(1024*1024*30) / float64(sliceItemSize))
	a := math.Max(1, sliceTotal)
	fmt.Println(a)
}
