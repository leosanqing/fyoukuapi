package main

import (
	"encoding/json"
	"fmt"
	"fyoukuapi/models"
	mqconn "fyoukuapi/service/mq"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 注册数据库
	beego.LoadAppConfig("ini", "../../conf/app.conf")
	defaultdb := beego.AppConfig.String("defaultdb")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", defaultdb, 30, 30)

	mqconn.Consumer("", "fyouku_send_message_user", callback)
}

func callback(msg string) {
	type Data struct {
		UserId    int
		MessageId int64
	}
	var data Data
	err := json.Unmarshal([]byte(msg), &data)
	if err != nil {
		fmt.Println("解析 Json 失败")
	}
	models.SendMessageUser(data.UserId, data.MessageId)
}
