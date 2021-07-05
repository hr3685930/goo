package main

import (
    "goo/cmd"
    "context"
    "fmt"
    "github.com/aaronjan/hunch"
    _ "go.uber.org/automaxprocs"
    "time"
)

func main() {
    boot := &cmd.Initialize{}
    ctx := context.Background()
    _, err := hunch.Waterfall(
        ctx,
        func(ctx context.Context, n interface{}) (interface{}, error) {
            filename := `config.toml`
            return nil, boot.Config(filename)
        },
        func(ctx context.Context, n interface{}) (interface{}, error) {
            return hunch.All(
                ctx,
                func(ctx context.Context) (interface{}, error) {
                    return nil, boot.Log("queue")
                },
                func(ctx context.Context) (interface{}, error) {
                    return nil, boot.Fs()
                },
                func(ctx context.Context) (interface{}, error) {
                    return hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
                        err := boot.Cache()
                        if err != nil {
                            fmt.Println("缓存重连中...", err)
                            time.Sleep(time.Second * 2)
                        }
                        return nil, err
                    })
                },
                func(ctx context.Context) (interface{}, error) {
                    return hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
                        err := boot.Queue()
                        if err != nil {
                            fmt.Println("队列重连中...", err)
                            time.Sleep(time.Second * 2)
                        }
                        return nil, err
                    })
                },
                func(ctx context.Context) (interface{}, error) {
                    return hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
                        err := boot.DB()
                        if err != nil {
                            fmt.Println("数据库重连中...", err)
                            time.Sleep(time.Second * 2)
                        }
                        return nil, err
                    })
                },
            )
        },
        func(ctx context.Context, n interface{}) (interface{}, error) {
            return nil, boot.Event()
        },
        func(ctx context.Context, n interface{}) (interface{}, error) {
            return nil, boot.Consumer()
        },
    )
    if err != nil {
        panic(err)
    }

    boot.Signal()
}
