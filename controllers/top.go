package controllers

import (
	"fyoukuapi/models"
	"github.com/astaxie/beego"
)

type TopController struct {
	beego.Controller
}

// 频道排行榜
// @router /channel/top  [*]
func (c *TopController) ChannelTop() {
	channelId, _ := c.GetInt("channelId")
	if 0 == channelId {
		c.Data["json"] = ReturnError(4004, "必须指定频道Id")
		c.ServeJSON()
	}

	count, videos, err := models.RedisGetChannelTop(channelId)

	if err != nil {
		c.Data["json"] = ReturnError(4001, "没有相关内容")
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnSuccess(0, "success", videos, count)
		c.ServeJSON()
	}
}

// 类型排行榜
// @router /type/top [*]
func (c *TopController) TypeTop() {
	typeId, _ := c.GetInt("typeId")
	if 0 == typeId {
		c.Data["json"] = ReturnError(4004, "必须指定类型Id")
		c.ServeJSON()
	}

	count, videos, err := models.RedisGetTypeTop(typeId)

	if err != nil {
		c.Data["json"] = ReturnError(4001, "没有相关内容")
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnSuccess(0, "success", videos, count)
		c.ServeJSON()
	}
}
