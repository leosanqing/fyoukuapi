package models

import "github.com/astaxie/beego/orm"

type Barrage struct {
	Id          int
	Content     string
	AddTime     int64
	UserId      int
	Status      int
	CurrentTime int64
	EpisodesId  int
	VideoId     int
}

type BarrageData struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	CurrentTime int    `json:"currentTime"`
}

func init() {
	orm.RegisterModel(new(Barrage))
}

func GetBarrageList(episodesId int, startTime int, endTime int) (int64, []BarrageData, error) {
	o := orm.NewOrm()
	var barrageDatas []BarrageData
	rows, err := o.Raw("SELECT id,content,`current_time` FROM barrage\n"+
		"WHERE status = 1 AND episodes_id =? AND `current_time` >= ? AND `current_time` <?\n"+
		"ORDER BY `current_time` ASC",
		episodesId, startTime, endTime).
		QueryRows(&barrageDatas)

	return rows, barrageDatas, err
}
