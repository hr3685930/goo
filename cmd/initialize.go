package cmd

import (
	"goo/config"
	"goo/internal/handler"
	"goo/pkg/cache"
	"goo/pkg/cache/redis"
	"goo/pkg/cache/sync"
	"goo/pkg/db"
	"goo/pkg/db/mysql"
	"goo/pkg/event"
	"goo/pkg/file/cos"
	"goo/pkg/grpcgw"
	echo2 "goo/pkg/http/echo"
	"goo/pkg/log/zap"
	"goo/pkg/queue"
	"goo/pkg/queue/kafka"
	"goo/pkg/queue/rabbitmq"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Initialize struct {
}

func (*Initialize) Config(filename string) error {
	// 优先级 显式调用Set函数>命令行参数>环境变量>配置文件
	config.Config(viper.GetViper())
	_, err := os.Stat(filename)
	if (err == nil) || (os.IsExist(err)) {
		viper.SetConfigFile(filename)
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			config.LoadConf()
		})
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.LoadConf()
	return nil
}

func (*Initialize) DB() error {
	con := viper.Get("db.connections").(map[string]map[string]string)
	defaultCon := viper.GetString("db.connection")
	debug := viper.GetBool("app.debug")
	for key, v := range con {
		switch v["driver"] {
		case "mysql":
			mysqls := mysql.NewMysql(key, v["database"], v["host"], v["port"], v["username"], v["password"], debug)
			err := mysqls.Connect()
			if err != nil {
				return err
			}
		}

		//set default conn
		if defaultCon == key {
			db.SetDefault(db.ORM(key))
		}
	}
	return nil
}

func (*Initialize) Log(typeName string) error {
	filepath := "./storage/log/"
	filename := typeName + ".log"
	//create file
	_, err := os.Stat(filepath + filename)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath, 0777)
			if err != nil {
				return err
			}
			f, err := os.Create(filepath + filename)
			if err != nil {
				return err
			}
			_ = f.Close()
		}
	}

	zap.NewZaps([]string{filepath + filename, "stdout"}).InitLog()
	return nil
}

func (*Initialize) Fs() error {
	fs := viper.Get("filesystem").(map[string]map[string]string)
	for _, v := range fs {
		switch v["driver"] {
		case "cos":
			co := cos.NewCosFile(v["secret_id"], v["secret_key"], v["region"], v["bucket"], 100*time.Second)
			co.InitCos()
		}
	}

	return nil
}

func (*Initialize) Cache() error {
	caches := viper.Get("caches").(map[string]map[string]string)
	defaultCache := viper.GetString("cache.drive")
	for key, v := range caches {
		switch v["driver"] {
		case "redis":
			database, _ := strconv.Atoi(v["database"])
			c, err := redis.New(v["host"], v["port"], database, v["auth"])
			if err != nil {
				return err
			}
			cache.CacheMap.Store(key, c)

		case "sync":
			c := sync.New()
			cache.CacheMap.Store(key, c)
		}

		//set default conn
		if defaultCache == key {
			cache.SetDefault(cache.GetCache(key))
		}
	}

	return nil
}

func (*Initialize) Event() error {
	e := &handler.Event{}
	e.Handler()
	return event.Register()
}

func (*Initialize) Consumer() error {
	c := &handler.Consumer{}
	return c.Handler()
}

func (*Initialize) Queue() error {
	queues := viper.Get("queues").(map[string]map[string]string)
	defaultQueue := viper.GetString("queue.drive")
	prefix := viper.GetString("app.name")
	for key, v := range queues {
		switch v["driver"] {
		case "rabbitmq":
			r := rabbitmq.NewRabbitMQ(v["user"], v["pass"], v["host"], v["port"], v["vhost"], prefix)
			queue.QueueMap.Store(key, r)
		case "kafka":
			k := kafka.NewKafka(v["addr"], prefix)
			queue.QueueMap.Store(key, k)
		}

		//set default queue
		if defaultQueue == key {
			queue.SetDefault(queue.GetMQ(key))
			err := queue.MQ.Connect()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (*Initialize) Metrics(port string) error {
	go func() {
		// create a new mux server
		server := http.NewServeMux() // register a new handler for the /metrics endpoint
		server.Handle("/metrics", promhttp.Handler())
		// start an http server using the mux server
		_ = http.ListenAndServe(":"+port, server)
	}()
	return nil
}

func (*Initialize) HTTP(port string, route func(e *echo.Echo)) error {
	e := &echo2.EchoHTTP{}
	s := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}
	return e.HTTP(s, route)
}

func (*Initialize) GW(rpcAddr, httpAddr string,grpcOpt []grpc.UnaryServerInterceptor, rpcHandler func(s *grpc.Server),
	httpHandler func(mux *runtime.ServeMux, e *echo.Echo, endpoint string) error) error {
	gw := &grpcgw.GW{}
	return gw.Register(rpcAddr, httpAddr, func(serverOpt []grpc.UnaryServerInterceptor) *grpc.Server {
		grpcOpt = append(grpcOpt, serverOpt...)
		s := grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(grpcOpt...)),
		)

		rpcHandler(s)
		return s

	}, func(gemux *runtime.ServeMux,endpoint string) (i http.Handler, e error) {
		ec := echo.New()
		echo.NotFoundHandler = func(c echo.Context) error {
			return echo.ErrMethodNotAllowed
		}
		ec.HTTPErrorHandler = echo2.CustomHTTPErrorHandler
		ec.Validator = echo2.NewCustomValidator()
		ec.Binder = echo2.NewCustomBinder()
		return ec, httpHandler(gemux, ec, endpoint)
	})
}

func (*Initialize) Cmd(cmds []cli.Command) error {
	app := cli.NewApp()
	app.Commands = cmds
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	return app.Run(os.Args)
}

func (*Initialize) Signal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			fmt.Print("exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
