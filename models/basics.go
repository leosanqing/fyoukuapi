package models

import "github.com/astaxie/beego/orm"

type ChannelRegion struct {
	Id        int
	Name      string
	Status    int
	AddTime   int64
	ChannelId int
	Sort      int
}

func init() {
	orm.RegisterModel(new(ChannelRegion))
}

type Type struct {
	Id   int
	Name string
}

func GetChannelRegion(channelId int) (int64, []ChannelRegion, error) {
	o := orm.NewOrm()
	var channelRegion []ChannelRegion
	count, err := o.QueryTable("channel_region").
		Filter("channel_id", channelId).
		Filter("status", 1).
		All(&channelRegion)
	return count, channelRegion, err
}

func GetChannelType(channelId int) (int64, []Type, error) {
	o := orm.NewOrm()
	var types []Type
	count, err := o.Raw("SELECT id,name\n"+
		"FROM channel_type\n"+
		"WHERE channel_id = ?\n"+
		"AND status = 1", channelId).
		QueryRows(&types)
	return count, types, err
}
