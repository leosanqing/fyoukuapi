package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Message struct {
	Id      int
	Content string
	AddTime int64
}

type MessageUser struct {
	Id        int
	MessageId int64
	AddTime   int64
	Status    int
}

func init() {
	orm.RegisterModel(new(Message), new(MessageUser))
}

func SendMessage(content string) (int64, error) {
	o := orm.NewOrm()
	var message Message
	message.Content = content
	message.AddTime = time.Now().Unix()
	return o.Insert(&message)
}

func SendMessageUser(userId int, messageId int64) error {
	o := orm.NewOrm()
	var messageUser MessageUser
	messageUser.AddTime = time.Now().Unix()
	messageUser.MessageId = messageId
	messageUser.Status = 1
	messageUser.Id = userId
	_, err := o.Insert(&messageUser)
	return err
}
