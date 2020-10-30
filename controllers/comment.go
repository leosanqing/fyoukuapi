package controllers

import (
	"fyoukuapi/models"
	"github.com/astaxie/beego"
)

type CommentController struct {
	beego.Controller
}

type CommentInfo struct {
	Id           int             `json:"id"`
	Content      string          `json:"content"`
	AddTime      int64           `json:"addTime"`
	AddTimeTitle string          `json:"addTimeTitle"`
	UserId       int             `json:"userId"`
	Stamp        int             `json:"stamp"`
	PraiseCount  int             `json:"praiseCount"`
	UserInfo     models.UserInfo `json:"userinfo"`
}

// 根据剧集数获取评论列表
// @router /comment/list [*]
func (c *CommentController) CommentList() {
	episodesId, _ := c.GetInt("episodesId")
	limit, _ := c.GetInt("limit")
	offset, _ := c.GetInt("offset")

	if 0 == episodesId {
		c.Data["json"] = ReturnError(4004, "必须指定视频剧情Id")
		c.ServeJSON()
	}
	if 0 == limit {
		limit = 12
	}

	count, comments, err := models.GetCommentList(episodesId, limit, offset)

	if err != nil {
		c.Data["json"] = ReturnError(4001, "没有相关内容")
		c.ServeJSON()
	} else {
		var commentInfo CommentInfo
		var commentInfos []CommentInfo

		// 获取 uid channel
		uidChan := make(chan int, 12)
		closeChan := make(chan bool, 5)
		resChan := make(chan models.UserInfo, 12)

		go func() {
			for _, comment := range comments {
				uidChan <- comment.UserId
			}
			close(uidChan)
		}()

		// 处理 uidChannel 中的信息
		for i := 0; i < 5; i++ {
			go func(uidChan chan int, resChan chan models.UserInfo, closeChan chan bool) {
				for uid := range uidChan {
					info, err := models.RedisGetUserInfo(uid)
					if err == nil {
						resChan <- info
					}
				}
				closeChan <- true
			}(uidChan, resChan, closeChan)
		}

		// 信息聚合
		go func() {
			for i := 0; i < 5; i++ {
				<-closeChan
			}
			close(resChan)
			close(closeChan)
		}()

		userInfoMap := make(map[int]models.UserInfo)
		for userInfo := range resChan {
			userInfoMap[userInfo.Id] = userInfo
		}

		for _, comment := range comments {
			commentInfo.Id = comment.Id
			commentInfo.Content = comment.Content
			commentInfo.AddTime = comment.AddTime
			commentInfo.AddTimeTitle = DateFormat(comment.AddTime)
			commentInfo.UserId = comment.UserId
			commentInfo.Stamp = comment.Stamp
			commentInfo.PraiseCount = comment.PraiseCount
			commentInfo.UserInfo, _ = userInfoMap[comment.UserId]
			commentInfos = append(commentInfos, commentInfo)
		}

		c.Data["json"] = ReturnSuccess(0, "Success", commentInfos, count)
		c.ServeJSON()
	}
}

//  发表评论
// @router /comment/save [post]
func (c *CommentController) SaveComment() {
	episodesId, _ := c.GetInt("episodesId")
	videoId, _ := c.GetInt("videoId")
	content := c.GetString("content")
	uid, _ := c.GetInt("uid")

	if 0 == episodesId {
		c.Data["json"] = ReturnError(4004, "必须指定视频剧情Id")
		c.ServeJSON()
	}

	if 0 == videoId {
		c.Data["json"] = ReturnError(4005, "必须指定视频Id")
		c.ServeJSON()
	}

	if "" == content {
		c.Data["json"] = ReturnError(4006, "必须输入内容")
		c.ServeJSON()
	}

	if 0 == uid {
		c.Data["json"] = ReturnError(4007, "请先登录")
		c.ServeJSON()
	}

	err := models.SaveComment(content, uid, episodesId, videoId)
	if err != nil {
		c.Data["json"] = ReturnError(4001, "没有相关内容")
		c.ServeJSON()
	} else {
		c.Data["json"] = ReturnSuccess(0, "success", nil, 0)
		c.ServeJSON()
	}

}
