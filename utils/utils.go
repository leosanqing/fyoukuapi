package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

func GetVideoName(uid string) string {
	hash := md5.New()
	hash.Write([]byte(uid + strconv.FormatInt(time.Now().Unix(), 10)))
	return hex.EncodeToString(hash.Sum(nil))
}
