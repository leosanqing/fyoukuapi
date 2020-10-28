package models

import (
	"encoding/json"
	redisClient "fyoukuapi/service/redis"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"strconv"
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

// 获取视频详情
func GetVideoInfo(videoId int) (Video, error) {
	o := orm.NewOrm()
	var video Video
	err := o.QueryTable("video").
		Filter("id", videoId).
		One(&video)
	return video, err
}

// 使用 Redis 缓存 ，改造获取视频详情
func RedisGetVideoInfo(videoId int) (Video, error) {
	var video Video
	conn := redisClient.PoolConnect()
	defer conn.Close()

	// 定义 RedisKey
	redisKey := "video:id:" + strconv.Itoa(videoId)

	exists, err := redis.Bool(conn.Do(redisClient.Exists, redisKey))
	if exists {
		values, _ := redis.Values(conn.Do(redisClient.HGetAll, redisKey))
		err = redis.ScanStruct(values, &video)
	} else {
		o := orm.NewOrm()
		err := o.QueryTable("video").
			Filter("id", videoId).
			One(&video)
		if err == nil {
			// 保存数据到redis
			_, err := conn.Do(
				redisClient.HmSet,
				redis.Args{redisKey}.AddFlat(video)...,
			)
			if err == nil {
				conn.Do(redisClient.Expire, redisKey, 86400)
			}
		}
	}
	return video, err
}

// 获取视频剧集列表
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

func RedisGetVideoEpisodesList(videoId int) (int64, []VideoEpisodes, error) {
	var (
		episodes []VideoEpisodes
		num      int64
		err      error
	)
	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "video:episodes:videoId" + strconv.Itoa(videoId)

	exists, err := redis.Bool(conn.Do(redisClient.Exists, redisKey))
	if exists {
		num, err = redis.Int64(conn.Do(redisClient.LLen, redisKey))
		if err == nil {
			values, _ := redis.Values(conn.Do(redisClient.LRange, redisKey, "0", "-1"))
			var episodesInfo VideoEpisodes
			for _, v := range values {
				err = json.Unmarshal(v.([]byte), &episodesInfo)
				if err == nil {
					episodes = append(episodes, episodesInfo)
				}
			}
		}
	} else {
		num, episodes, err = GetVideoEpisodesList(videoId)
		if err == nil {
			// 保存数据到redis
			for _, episode := range episodes {
				marshal, err := json.Marshal(episode)
				if err == nil {
					conn.Do(redisClient.RPush, redisKey, marshal)
				}
			}
			conn.Do(redisClient.Expire, redisKey, time.Hour*24)
		}
	}

	return num, episodes, err
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

func RedisGetChannelTop(channelId int) (int64, []Video, error) {
	var (
		videos []Video
		num    int64
		err    error
	)

	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "video:top:channel:channelId:" + strconv.Itoa(channelId)

	exists, err := redis.Bool(conn.Do(redisClient.Exists, redisKey))
	if exists {
		num = 0
		res, _ := redis.Values(conn.Do(redisClient.ZRevRange, redisKey, 0, 10, redisClient.WithScores))
		for k, v := range res {
			// 获取到id
			if k%2 == 0 {
				videoId, err := strconv.Atoi(string(v.([]byte)))
				if err == nil {
					videoInfo, err := RedisGetVideoInfo(videoId)
					if err == nil {
						videos = append(
							videos,
							Video{
								Id:            videoInfo.Id,
								Img:           videoInfo.Img,
								Img1:          videoInfo.Img1,
								IsEnd:         videoInfo.IsEnd,
								SubTitle:      videoInfo.SubTitle,
								Title:         videoInfo.Title,
								AddTime:       videoInfo.AddTime,
								Comment:       videoInfo.Comment,
								EpisodesCount: videoInfo.EpisodesCount,
							})
						num++
					}
				}
			}
		}
	} else {
		num, videos, err = GetChannelTop(channelId)
		if err == nil {
			// 保存数据到redis
			for _, v := range videos {
				conn.Do(redisClient.ZAdd, redisKey, v.Comment, v.Id)
			}
			conn.Do(redisClient.Expire, redisKey, 30*time.Second)
		}
	}
	return num, videos, err
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

func RedisGetTypeTop(typeId int) (int64, []Video, error) {
	var (
		videos []Video
		num    int64
		err    error
	)

	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "video:top:type:typeId:" + strconv.Itoa(typeId)

	exists, err := redis.Bool(conn.Do(redisClient.Exists, redisKey))
	if exists {
		num = 0
		res, _ := redis.Values(conn.Do(redisClient.ZRevRange, redisKey, 0, 10, redisClient.WithScores))
		for k, v := range res {
			// 获取到id
			if k%2 == 0 {
				videoId, err := strconv.Atoi(string(v.([]byte)))
				if err == nil {
					videoInfo, err := RedisGetVideoInfo(videoId)
					if err == nil {
						videos = append(
							videos,
							Video{
								Id:            videoInfo.Id,
								Img:           videoInfo.Img,
								Img1:          videoInfo.Img1,
								IsEnd:         videoInfo.IsEnd,
								SubTitle:      videoInfo.SubTitle,
								Title:         videoInfo.Title,
								AddTime:       videoInfo.AddTime,
								Comment:       videoInfo.Comment,
								EpisodesCount: videoInfo.EpisodesCount,
							})
						num++
					}
				}
			}
		}
	} else {
		num, videos, err = GetTypeTop(typeId)
		if err == nil {
			// 保存数据到redis
			for _, v := range videos {
				conn.Do(redisClient.ZAdd, redisKey, v.Comment, v.Id)
			}
			conn.Do(redisClient.Expire, redisKey, 30*time.Second)
		}
	}
	return num, videos, err
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
