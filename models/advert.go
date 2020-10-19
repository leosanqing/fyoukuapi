package models

import "github.com/astaxie/beego/orm"

type Advert struct {
	Id        int
	Title     string
	SubTitle  string
	ChannelId int
	Img       string
	Sort      string
	AddTime   int64
	Url       string
	Status    int
}

func init() {
	orm.RegisterModel(new(Advert))
}

func GetChannelAdvert(channelId int) (int64, []Advert, error) {
	o := orm.NewOrm()
	var advert []Advert

	rows, err := o.Raw(
		"SELECT id, title, sub_title, img, add_time, url \n" +
			"FROM advert \n" +
			"WHERE status = 1 \n" +
			"AND channel_id =? \n" +
			"ORDER BY sort DESC \n" +
			"LIMIT 1").
		SetArgs(channelId).
		QueryRows(&advert)

	return rows, advert, err
}
