package utils

import (
	"math/rand"
	"regexp"
	"time"
)

const letterBytes = "abcdefghijkl012346789mnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

//RandStringBytesMaskImprSrc returns random string
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

//Sanitize a string without removing dots: [^a-zA-Z0-9.]+
func Sanitize(filename string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9.]+")
	if err != nil {
		return filename, err
	}
	return reg.ReplaceAllString(filename, ""), nil
}
