package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

// MD5 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	_, _ = ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func Base64Encode(text string) string {
	encodeString := base64.StdEncoding.EncodeToString([]byte(text))
	return encodeString
}

func Base64Decode(text string) (string, error) {
	if encodeString, err := base64.StdEncoding.DecodeString(text); err != nil {
		return "", err
	} else {
		return string(encodeString[:]), nil
	}
}
