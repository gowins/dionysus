package algs

import (
	"math/rand"
	"time"
)

const (
	upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower = "abcdefghijklmnopqrstuvwxyz"
	num   = "0123456789"

	chars        = upper + lower
	charsWithNum = chars + num
)

func init() {
	rand.Seed(time.Now().Unix())
}

func FirstNotEmpty(strs ...string) string {
	for _, str := range strs {
		if str != "" {
			return str
		}
	}

	return ""
}

func RandStr(length int, withNum bool) string {
	if length < 1 {
		return ""
	}

	c := len(chars)
	if withNum {
		c = len(charsWithNum)
	}

	ret := make([]byte, length)

	for i := 0; i < length; i++ {
		ret[i] = charsWithNum[rand.Intn(c)]
	}
	return string(ret)
}
