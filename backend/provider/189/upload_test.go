package _189

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestMd5Slice(t *testing.T) {
	const sliceItemSize uint64 = 10485760
	var sliceTotal = math.Ceil(float64(1024*1024*10) / float64(sliceItemSize))
	a := math.Max(1, sliceTotal)
	fmt.Println(a)
}

func Test1(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Float32())
}
