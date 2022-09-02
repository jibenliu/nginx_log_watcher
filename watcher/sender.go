package watcher

import (
	"github.com/idoubi/goz"
)

type sendStruct struct {
	name   string
	level  string
	detail string
	time   string
}

func (s *sendStruct) String() string {
	return s.name + "\n" + s.level + "\n" + s.detail + "\n" + s.time
}

var cli = goz.NewClient()

func sendAlarm(send *sendStruct) {
	resp, err := cli.Post("https://api.telegram.org/bot5362907303:AAH9xZI8UBi_EIjGKIlTYcPOtiexif0oN5w/sendMessage", goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
		JSON: map[string]interface{}{
			"chat_id": "-644611182",
			"text":    send.String(),
		},
	})
	if err != nil {
		Error("send msg to robot error " + err.Error())
		return
	}
	_, err = resp.GetBody()
	if err != nil {
		Error("get curl body msg error " + err.Error())
		return
	}
}
