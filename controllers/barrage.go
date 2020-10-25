package controllers

import (
	"encoding/json"
	"fyoukuapi/models"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"net/http"
)

type BarrageController struct {
	beego.Controller
}

type WsData struct {
	CurrentTime int
	EpisodesId  int
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// @router /barrage/ws [*]
func (c *BarrageController) BarrageWs() {
	var (
		conn     *websocket.Conn
		err      error
		data     []byte
		barrages []models.BarrageData
	)

	if conn, err = upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil); err != nil {
		goto ERR
	}
	for {
		if _, data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}

		var wsData WsData
		json.Unmarshal(data, &wsData)
		endTime := wsData.CurrentTime + 60

		// 获取弹幕数据
		_, barrages, err = models.GetBarrageList(wsData.EpisodesId, wsData.CurrentTime, endTime)
		if err == nil {
			if err := conn.WriteJSON(barrages); err != nil {
				goto ERR
			}
		}
	}

ERR:
	conn.Close()
}

// @router /barrage/save [post]
func (c *BarrageController) Save() {
	episodesId, _ := c.GetInt("episodesId")
	videoId, _ := c.GetInt("videoId")
	uid, _ := c.GetInt("uid")
	content := c.GetString("content")
	currentTime, _ := c.GetInt64("currentTime")

	if 0 == episodesId {
		c.Data["json"] = ReturnError(4001, "必须指定剧情集数")
		c.ServeJSON()
	}

	if 0 == videoId {
		c.Data["json"] = ReturnError(4002, "必须指定视频类型")
		c.ServeJSON()
	}
	if 0 == uid {
		c.Data["json"] = ReturnError(4003, "请登录")
		c.ServeJSON()
	}
	if "" == content {
		c.Data["json"] = ReturnError(4004, "请输入内容")
		c.ServeJSON()
	}

	if 0 == currentTime {
		c.Data["json"] = ReturnError(4003, "请指定时间")
		c.ServeJSON()
	}

	err := models.Save(content, currentTime, uid, episodesId, videoId)
	if err != nil {
		c.Data["json"] = ReturnError(0, "发送弹幕失败")
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnSuccess(0, "发送成功", nil, 0)
		c.ServeJSON()
	}
}
