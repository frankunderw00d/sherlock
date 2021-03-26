package encrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type ()

const ()

var ()

func init() {}

// MD5 摘要加密
func MD5(str ...string) string {
	if len(str) <= 0 {
		return ""
	}

	targetStr := ""
	for _, s := range str {
		targetStr += s
	}
	data := []byte(targetStr)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

// SHA1哈希
func Sha1(str string) string {
	o := sha1.New()
	o.Write([]byte(str))
	return hex.EncodeToString(o.Sum(nil))
}

// base64 加密
func Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
