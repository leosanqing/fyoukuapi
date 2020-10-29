package main

import (
	"encoding/json"
	"fmt"
	"fyoukuapi/models"
	mqconn "fyoukuapi/service/mq"
	redisconn "fyoukuapi/service/redis"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func main() {
	// 注册数据库
	beego.LoadAppConfig("ini", "../../conf/app.conf")
	defaultdb := beego.AppConfig.String("defaultdb")

	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		panic(err)
	}
	err = orm.RegisterDataBase("default", "mysql", defaultdb)
	if err != nil {
		panic(err)
	}
	mqconn.Consumer("", "fyouku_top", callback)
}

func callback(msg string) {
	type Data struct {
		VideoId int `json:"videoId"`
	}
	var data Data
	err := json.Unmarshal([]byte(msg), &data)
	if err != nil {
		fmt.Println("解析 Json 失败")
	}
	videoInfo, err := models.RedisGetVideoInfo(data.VideoId)
	if err == nil {
		// 更新排行榜
		conn := redisconn.PoolConnect()
		defer conn.Close()
		redisChannelKey := "video:top:channel:channelId:" + strconv.Itoa(videoInfo.ChannelId)
		redisTypeKey := "video:top:type:typeId:" + strconv.Itoa(videoInfo.TypeId)
		conn.Do("zincrby", redisChannelKey, 1, data.VideoId)
		conn.Do("zincrby", redisTypeKey, 1, data.VideoId)
	}
	fmt.Println("msg is :" + msg)
}
