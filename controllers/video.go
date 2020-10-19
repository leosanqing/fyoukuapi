package controllers

import (
	"fyoukuapi/models"
	"github.com/astaxie/beego"
	"log"
)

type VideoController struct {
	beego.Controller
}

// 频道页 - 获取顶部广告
// @router /channel/advert [*]
func (c *VideoController) ChannelAdvert() {
	channelId, _ := c.GetInt("channelId")
	if 0 == channelId {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
	}
	num, adverts, err := models.GetChannelAdvert(channelId)
	if err != nil {
		c.Data["json"] = ReturnError(4004, "请求数据失败，请稍后重试")
	} else {
		c.Data["json"] = ReturnSuccess(0, "请求成功", adverts, num)
	}
	c.ServeJSON()
}

// 获取热播频道
// @router /channel/hot [*]
func (c *VideoController) ChannelHotList() {
	channelId, _ := c.GetInt("channelId")
	if 0 == channelId {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
	}
	count, data, err := models.GetChannelHotList(channelId)
	if err != nil {
		c.Data["json"] = ReturnError(4001, "没有相关内容")
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnSuccess(0, "Success", data, count)
		c.ServeJSON()
	}
}

// 频道页-根据频道地区获取推荐的视频
// @router /channel/recommend/region [*]
func (c *VideoController) ChannelRegionRecommendList() {
	var (
		channelId int
		regionId  int
	)
	channelId, _ = c.GetInt("channelId")
	if 0 == channelId {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
	}

	regionId, _ = c.GetInt("regionId")
	if 0 == regionId {
		c.Data["json"] = ReturnError(4002, "必须指定频道地区")
		c.ServeJSON()
	}

	count, data, err := models.GetChannelRegionRecommend(channelId, regionId)
	if err != nil {
		c.Data["json"] = ReturnError(4001, "没有相关内容")
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnSuccess(0, "Success", data, count)
		c.ServeJSON()
	}
}

// 频道页-根据类型获取推荐的视频
// @router /channel/recommend/type [*]
func (c *VideoController) ChannelTypeRecommendList() {
	var (
		channelId int
		typeId    int
	)
	channelId, _ = c.GetInt("channelId")
	if 0 == channelId {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
	}

	typeId, _ = c.GetInt("typeId")
	if 0 == typeId {
		c.Data["json"] = ReturnError(4002, "必须指定频道类型")
		c.ServeJSON()
	}

	count, data, err := models.GetChannelTypeRecommend(channelId, typeId)
	if err != nil {
		log.Println(err)
		c.Data["json"] = ReturnError(4001, "没有相关内容")
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnSuccess(0, "Success", data, count)
		c.ServeJSON()
	}
}

// 获取视频列表
// @router /channel/video [*]
func (c *VideoController) ChannelVideo() {
	channelId, _ := c.GetInt("channelId")
	regionId, _ := c.GetInt("regionId")
	typeId, _ := c.GetInt("typeId")
	end := c.GetString("end")
	sort := c.GetString("sort")
	limit, _ := c.GetInt("limit")
	offset, _ := c.GetInt("offset")

	if 0 == channelId {
		c.Data["json"] = ReturnError(4001, "必须指定频道")
		c.ServeJSON()
	}

	if 0 == limit {
		limit = 12
	}
	list, params, err := models.GetChannelVideoList(channelId, regionId, typeId, end, sort, offset, limit)

	if err != nil {
		c.Data["json"] = ReturnError(4001, "没有相关内容")
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnSuccess(0, "Success", params, list)
		c.ServeJSON()
	}
}
