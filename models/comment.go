package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Comment struct {
	Id          int
	Content     string
	AddTime     int64
	UserId      int
	Status      int
	Stamp       int
	PraiseCount int
	EpisodesId  int
	VideoId     int
}

func init() {
	orm.RegisterModel(new(Comment))
}

func GetCommentList(episodesId int, limit int, offset int) (int64, []Comment, error) {
	o := orm.NewOrm()
	var comments []Comment

	rows, _ := o.Raw("SELECT id \n"+
		"FROM comment \n"+
		"WHERE status =1 \n"+
		"AND episodes_id = ?", episodesId).
		QueryRows(&comments)

	_, err := o.QueryTable("comment").
		Filter("episodes_id", episodesId).
		Filter("status", 1).
		Limit(limit, offset).
		OrderBy("-add_time").
		All(&comments)

	return rows, comments, err

}

func SaveComment(content string, uid int, episodesId int, videoId int) error {
	o := orm.NewOrm()
	var comment Comment
	comment.UserId = uid
	comment.Content = content
	comment.Stamp = 0
	comment.AddTime = time.Now().Unix()
	comment.Status = 1
	comment.EpisodesId = episodesId
	comment.VideoId = videoId

	_, err := o.Insert(&comment)
	if err == nil {
		o.Raw("UPDATE video SET comment = comment+1 WHERE id = ?", videoId).Exec()
		// 修改 视频剧情的评论数
		o.Raw("UPDATE video_episodes SET comment = comment+1 WHERE id=?", episodesId).Exec()
	}
	return err
}
