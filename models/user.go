package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id       int
	Name     string
	Password string
	AddTime  int64
	Status   int
	Mobile   string
	Avatar   string
}

func init() {
	orm.RegisterModel(new(User))
}

// 根据手机号判断用户是否存在
func IsUserExistByMobile(mobile string) bool {
	o := orm.NewOrm()
	user := User{Mobile: mobile}
	err := o.Read(&user, "Mobile")
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		return false
	}
	return true
}

// 保存用户
func UserSave(mobile string, password string) error {
	o := orm.NewOrm()
	_, err := o.Insert(
		&User{
			Name:     "",
			Mobile:   mobile,
			Password: password,
			Status:   1,
			AddTime:  time.Now().Unix()})
	return err
}

func UserLogin(mobile string, password string) (int, string) {
	o := orm.NewOrm()
	var user User
	err := o.QueryTable("user").
		Filter("mobile", mobile).
		Filter("password", password).
		One(&user)
	if err == orm.ErrMissPK || err == orm.ErrNoRows {
		return 0, ""
	}
	return user.Id, user.Name
}
