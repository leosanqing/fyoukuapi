package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Video struct {
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

type VideoData struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	SubTitle      string `json:"subTitle"`
	Img           string `json:"img"`
	Img1          string `json:"img1"`
	AddTime       int64  `json:"addTime"`
	IsEnd         int    `json:"isEnd"`
	EpisodesCount int    `json:"episodesCount"`
}

type VideoEpisodes struct {
	Id            int
	Title         string
	AddTime       int64
	Num           int
	VideoId       int
	PlayUrl       string
	Status        int
	Comment       int
	AliyunVideoId int
}

func init() {
	orm.RegisterModel(new(Video))
}

// 获取热播频道列表
func GetChannelHotList(channelId int) (int64, []Video, error) {
	o := orm.NewOrm()
	var video []Video
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
func GetChannelRegionRecommend(channelId int, regionId int) (int64, []Video, error) {
	o := orm.NewOrm()
	var video []Video
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

func GetChannelTypeRecommend(channelId int, typeId int) (int64, []Video, error) {
	o := orm.NewOrm()
	var video []Video
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

func GetChannelVideoList(
	channelId int, regionId int, typeId int,
	end string, sort string, offset int, limit int) (int64, []orm.Params, error) {
	o := orm.NewOrm()
	var videos []orm.Params

	qs := o.QueryTable("video").
		Filter("channel_id", channelId).
		Filter("status", 1)

	if regionId > 0 {
		qs = qs.Filter("region_id", regionId)
	}
	if typeId > 0 {
		qs = qs.Filter("type_id", typeId)
	}
	if end == "n" {
		qs = qs.Filter("is_end", 0)
	} else if "y" == end {
		qs = qs.Filter("is_end", 1)
	}

	if "episodesUpdateTime" == sort {
		qs = qs.OrderBy("-episodes_update_time")
	} else if "comment" == sort {
		qs = qs.OrderBy("-comment")
	} else if "addTime" == sort {
		qs = qs.OrderBy("-add_time")
	} else {
		qs = qs.OrderBy("-add_time")
	}

	values, _ := qs.Values(&videos, "id", "title", "sub_title", "img", "img1", "add_time", "episodes_count", "is_end")
	qs = qs.Limit(limit, offset)
	_, err := qs.Values(&videos, "id", "title", "sub_title", "img", "img1", "add_time", "episodes_count", "is_end")

	return values, videos, err
}

func GetVideoInfo(videoId int) (Video, error) {
	o := orm.NewOrm()
	var video Video
	err := o.QueryTable("video").
		Filter("id", videoId).
		One(&video)
	return video, err
}

func GetVideoEpisodesList(videoId int) (int64, []VideoEpisodes, error) {
	o := orm.NewOrm()
	var episodes []VideoEpisodes
	rows, err := o.Raw("SELECT id,title,add_time,num,play_url,comment\n"+
		"FROM video_episodes\n"+
		"WHERE video_id = ?\n"+
		"AND status = 1\n"+
		"ORDER BY num ASC", videoId).
		QueryRows(&episodes)

	return rows, episodes, err
}

func GetChannelTop(channelId int) (int64, []Video, error) {
	o := orm.NewOrm()
	var video []Video
	rows, err := o.QueryTable("video").
		Filter("channelId", channelId).
		Filter("status", 1).
		OrderBy("-comment").
		Limit(10).
		All(&video)

	return rows, video, err
}

func GetTypeTop(typeId int) (int64, []Video, error) {
	o := orm.NewOrm()
	var video []Video
	rows, err := o.QueryTable("video").
		Filter("typeId", typeId).
		Filter("status", 1).
		OrderBy("-comment").
		Limit(10).
		All(&video)
	return rows, video, err
}

func GetUserVideo(uid int) (int64, []Video, error) {
	o := orm.NewOrm()
	var videos []Video
	rowsCount, err := o.Raw(
		"SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end\n"+
			"FROM video\n"+
			"WHERE user_id = ? \n"+
			"ORDER BY add_time DESC", uid).
		QueryRows(&videos)
	return rowsCount, videos, err
}

func SaveVideo(title string, subTitle string, channelId int, regionId int, typeId int, playUrl string, userId int) error {
	o := orm.NewOrm()
	id, err := o.Insert(&Video{
		Title:              title,
		SubTitle:           subTitle,
		AddTime:            time.Now().Unix(),
		Img:                "",
		Img1:               "",
		EpisodesCount:      1,
		IsEnd:              1,
		ChannelId:          channelId,
		Status:             1,
		RegionId:           regionId,
		TypeId:             typeId,
		EpisodesUpdateTime: time.Now().Unix(),
		Comment:            0,
		UserId:             userId,
	})
	if err == nil {
		o.Raw("INSERT INTO video_episodes \n"+
			"(title,add_time,num,video_id,play_url,status,comment)\n"+
			"VALUES (?,?,?,?,?,?,?)",
			subTitle, time.Now().Unix(), 1, id, playUrl, 1, 0).Exec()
	}
	return err
}
