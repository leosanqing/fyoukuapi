package controllers

import (
	"fyoukuapi/models"
	"github.com/astaxie/beego"
	"regexp"
	"strconv"
	"strings"
)

type UserController struct {
	beego.Controller
}

// 用户注册
// @router /register/save [post]
func (c *UserController) SaveRegister() {
	mobile := c.GetString("mobile")
	password := c.GetString("password")

	if mobile == "" {
		c.Data["json"] = ReturnError(4001, "手机号不能为空")
		c.ServeJSON()
	}
	isMatch, _ := regexp.MatchString(`^1([34578])[0-9]\d{8}$`, mobile)
	if !isMatch {
		c.Data["json"] = ReturnError(4002, "手机号码格式不正确")
		c.ServeJSON()
	}
	if "" == password {
		c.Data["json"] = ReturnError(4003, "密码不能为空")
		c.ServeJSON()
	}

	userExistByMobile := models.IsUserExistByMobile(mobile)
	if userExistByMobile {
		c.Data["json"] = ReturnError(4004, "该手机号已注册")
		c.ServeJSON()
	} else {
		err := models.UserSave(mobile, EncryptPassword(password))
		if err == nil {
			c.Data["json"] = ReturnSuccess(0, "注册成功", nil, 0)
			c.ServeJSON()
		} else {
			c.Data["json"] = ReturnError(5000, err)
			c.ServeJSON()
		}
	}
}

// @router /login/do [*]
func (c *UserController) LoginDo() {
	var (
		mobile   string
		password string
	)
	mobile = c.GetString("mobile")
	password = c.GetString("password")

	if mobile == "" {
		c.Data["json"] = ReturnError(4001, "手机号不能为空")
		c.ServeJSON()
	}
	isMatch, _ := regexp.MatchString(`^1(3|4|5|7|8)[0-9]\d{8}$`, mobile)
	if !isMatch {
		c.Data["json"] = ReturnError(4002, "手机号码格式不正确")
		c.ServeJSON()
	}
	if "" == password {
		c.Data["json"] = ReturnError(4003, "密码不能为空")
		c.ServeJSON()
	}

	id, name := models.UserLogin(mobile, EncryptPassword(password))
	if 0 != id {
		c.Data["json"] = ReturnSuccess(0, "登录成功",
			map[string]interface{}{"uid": id, "username": name}, 1)
	} else {
		c.Data["json"] = ReturnError(4004, "手机号或密码不正确")
	}
	c.ServeJSON()
}

// 发送消息
// @router /send/message [post]
func (c *UserController) SendMessage() {
	uids := c.GetString("uids")
	content := c.GetString("content")
	if uids == "" {
		c.Data["json"] = ReturnError(4001, "uid 不能为空")
		c.ServeJSON()
	}

	if content == "" {
		c.Data["json"] = ReturnError(4002, "发送内容 不能为空")
		c.ServeJSON()
	}

	messageId, err := models.SendMessage(content)

	if err == nil {
		uid := strings.Split(uids, ",")
		for _, s := range uid {
			userId, _ := strconv.Atoi(s)
			//models.SendMessageUser(userId, messageId)
			models.SendMessageUserMq(userId, messageId)
		}
		c.Data["json"] = ReturnSuccess(0, "发送成功", nil, 1)
	} else {
		c.Data["json"] = ReturnError(4004, "发送消息失败")
	}
	c.ServeJSON()
}
