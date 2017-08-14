package wechatpay

import (
	"fmt"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func Atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return i
}

func Atoi64(a string) int64 {
	i, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func Atof64(a string) float64 {
	i, err := strconv.ParseFloat(a, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func SliceAtoi(a []string) []int {
	var r = make([]int, len(a))
	for i, j := 0, len(a); i < j; i++ {
		r[i] = Atoi(a[i])
	}

	return r
}

func UUID(length int) string {
	u4 := uuid.NewV4()
	s := strings.Replace(fmt.Sprintf("%s", u4), "-", "", -1)
	return s[:length]
}
