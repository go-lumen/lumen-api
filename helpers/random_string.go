package helpers

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIDxBits = 6
	letterIDxMask = 1<<letterIDxBits - 1
	letterIDxMax  = 63 / letterIDxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIDxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIDxMax
		}
		if idx := int(cache & letterIDxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIDxBits
		remain--
	}

	return string(b)
}
