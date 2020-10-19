package controllers

import (
	"fyoukuapi/models"
	"github.com/astaxie/beego"
)

type BasicController struct {
	beego.Controller
}

//获取频道下地区信息
// @router /channel/region [*]
func (c *BasicController) GetChannelRegion() {
	channelId, _ := c.GetInt("channelId")
	if 0 == channelId {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
	}

	count, data, err := models.GetChannelRegion(channelId)
	if err != nil {
		c.Data["json"] = ReturnError(4001, "没有相关内容")
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnSuccess(0, "Success", data, count)
		c.ServeJSON()
	}
}

//获取频道下类型信息
// @router /channel/type [*]
func (c *BasicController) GetChannelType() {
	channelId, _ := c.GetInt("channelId")
	if 0 == channelId {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
	}

	count, data, err := models.GetChannelType(channelId)
	if err != nil {
		c.Data["json"] = ReturnError(4001, "没有相关内容")
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnSuccess(0, "Success", data, count)
		c.ServeJSON()
	}
}
