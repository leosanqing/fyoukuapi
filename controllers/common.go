package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
)

type CommonController struct {
	beego.Controller
}

type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Items interface{} `json:"items"`
	Count int64       `json:"count"`
}

func ReturnSuccess(code int, msg interface{}, items interface{}, count int64) *JsonStruct {
	return &JsonStruct{Code: code, Msg: msg, Items: items, Count: count}
}

func ReturnError(code int, msg interface{}) *JsonStruct {
	return &JsonStruct{Code: code, Msg: msg}
}

// 密码加密
func EncryptPassword(password string) string {
	hash := md5.New()
	hash.Write([]byte(password + beego.AppConfig.String("md5code")))
	return hex.EncodeToString(hash.Sum(nil))
}
