package job

import (
    "goo/pkg/queue"
    "fmt"
)

type WechatTextMsg struct {
    AppID     string `json:"app_id"`
    AppSecret string `json:"app_secret"`
    OpenID    string `json:"open_id"`
    Content   string `json:"content"`
}

func (w *WechatTextMsg) Handler() (queueErr *queue.Error) {
    defer func() {
        if err := recover(); err != nil {
            queueErr = queue.Err(fmt.Errorf("error send text msg: %+v", err))
        }
    }()

    return
}
