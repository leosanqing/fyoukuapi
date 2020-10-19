package models

import "github.com/astaxie/beego/orm"

type VideoData struct {
	Id                 int
	Title              string
	SubTitle           string
	Status             int
	Img                string
	Img1               string
	AddTime            int64
	ChannelId          int
	TypeId             int
	RegionId           int
	UserId             int
	EpisodesCount      int
	EpisodesUpdateTime int64
	IsHot              int
	IsEnd              int
	IsRecommend        int
	Comment            int
}

func init() {
	orm.RegisterModel(new(VideoData))
}

// 获取热播频道列表
func GetChannelHotList(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var video []VideoData
	rows, err := o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end\n"+
		"FROM video\n"+
		"WHERE channel_id = ? \n"+
		"AND is_hot = 1\n"+
		"AND status = 1\n"+
		"ORDER BY episodes_update_time DESC\n"+
		"LIMIT 9", channelId).
		QueryRows(&video)
	return rows, video, err
}

// 获取地区推荐
func GetChannelRegionRecommend(channelId int, regionId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var video []VideoData
	rows, err := o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end\n"+
		"FROM video\n"+
		"WHERE channel_id = ? \n"+
		"AND region_id = ?\n"+
		"AND is_recommend = 1\n"+
		"AND status = 1\n"+
		"ORDER BY episodes_update_time DESC\n"+
		"LIMIT 9", channelId, regionId).
		QueryRows(&video)
	return rows, video, err
}

func GetChannelTypeRecommend(channelId int, typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var video []VideoData
	rows, err := o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end\n"+
		"FROM video\n"+
		"WHERE channel_id = ? \n"+
		"AND type_id = ?\n"+
		"AND is_recommend = 1\n"+
		"AND status = 1\n"+
		"ORDER BY episodes_update_time DESC\n"+
		"LIMIT 9", channelId, typeId).
		QueryRows(&video)
	return rows, video, err
}
