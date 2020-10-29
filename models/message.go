package models

import (
	"encoding/json"
	mqconn "fyoukuapi/service/mq"
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
	UserId    int
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
	messageUser.UserId = userId
	_, err := o.Insert(&messageUser)
	return err
}

// 保存消息接收人到队列中
func SendMessageUserMq(userId int, messageId int64) {
	type Data struct {
		UserId    int
		MessageId int64
	}
	dataJson, _ := json.Marshal(Data{
		UserId:    userId,
		MessageId: messageId,
	})
	mqconn.Publish("", "fyouku_send_message_user", string(dataJson))
}
