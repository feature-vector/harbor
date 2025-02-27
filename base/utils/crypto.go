package utils

import (
	"crypto"
	"encoding/hex"
)

func MD5(src string) string {
	md5 := crypto.MD5.New()
	md5.Write([]byte(src))
	return hex.EncodeToString(md5.Sum(nil))
}

func SHA256(src string) string {
	sha256 := crypto.SHA256.New()
	sha256.Write([]byte(src))
	return hex.EncodeToString(sha256.Sum(nil))
}
