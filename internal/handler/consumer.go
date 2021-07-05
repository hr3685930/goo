package handler

import (
	"goo/internal/job"
	"goo/pkg/queue"
	"fmt"
)

type Consumer struct {

}

func (*Consumer) Handler() (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("Consumer error : %+v", err)
		}
	}()
	go queue.NewConsumer("wechat", &job.WechatTextMsg{}, 0, 0)
	go queue.NewConsumer("wechat", &job.WechatNewsMsg{}, 0, 0)

	return nil
}