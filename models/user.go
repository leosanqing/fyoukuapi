package models

import (
	redisClient "fyoukuapi/service/redis"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"strconv"
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

type UserInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	AddTime int64  `json:"addTime"`
	Avatar  string `json:"avatar"`
}

func init() {
	orm.RegisterModel(new(User), new(UserInfo))
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

func GetUserInfo(uid int) (UserInfo, error) {
	o := orm.NewOrm()
	var userInfo UserInfo

	err := o.Raw("SELECT id, name, add_time, avatar\n"+
		"FROM user\n"+
		"WHERE id = ? \n"+
		"LIMIT 1", uid).
		QueryRow(&userInfo)
	return userInfo, err
}

func RedisGetUserInfo(uid int) (UserInfo, error) {
	var userInfo UserInfo
	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "user:id:" + strconv.Itoa(uid)

	exists, err := redis.Bool(conn.Do(redisClient.Exists, redisKey))
	if exists {
		values, _ := redis.Values(conn.Do(redisClient.HGetAll, redisKey))
		err = redis.ScanStruct(values, &userInfo)
	} else {
		userInfo, err = GetUserInfo(uid)
		if err == nil {
			// 保存数据到redis
			_, err := conn.Do(
				redisClient.HmSet,
				redis.Args{redisKey}.AddFlat(userInfo)...,
			)
			if err == nil {
				conn.Do(redisClient.Expire, redisKey, 86400)
			}
		}
	}
	return userInfo, err
}
