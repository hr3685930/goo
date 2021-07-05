package job

import (
    "goo/pkg/queue"
    "fmt"
)

type WechatNewsMsg struct {
    AppID     string                `json:"app_id"`
    AppSecret string                `json:"app_secret"`
    OpenID    string                `json:"open_id"`
}

func (w *WechatNewsMsg) Handler() (queueErr *queue.Error) {
    defer func() {
        if err := recover(); err != nil {
            queueErr = queue.Err(fmt.Errorf("error send news msg: %+v", err))
        }
    }()

    return
}
