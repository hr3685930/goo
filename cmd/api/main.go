package main

import (
	"goo/cmd"
	"goo/internal/handler"
	"goo/internal/server"
	"goo/internal/svc"
	pb "goo/proto"
	"context"
	"fmt"
	"github.com/aaronjan/hunch"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
					return nil, boot.Log("api")
				},
				func(ctx context.Context) (interface{}, error) {
					return nil, boot.Fs()
				},
				func(ctx context.Context) (interface{}, error) {
					return nil, boot.Metrics("9001")
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
			return nil, boot.GW(":8081", ":8080", handler.Middleware() ,func(s *grpc.Server) {
				if viper.GetBool("app.debug") {
					reflection.Register(s)
				}
				ctxs := svc.NewServiceContext()
				pb.RegisterUserServer(s, server.NewUser(ctxs))
			}, func(mux *runtime.ServeMux, e *echo.Echo, endpoint string)  error {
				dialOption := []grpc.DialOption{grpc.WithInsecure()}
				handler.Route(e, mux)
				return pb.RegisterUserGWFromEndpoint(ctx, mux, endpoint, dialOption)
			})
		},
	)
	if err != nil {
		panic(err)
	}

	boot.Signal()
}
