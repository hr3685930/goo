package config

import (
    "github.com/spf13/viper"
)

func Config(viper *viper.Viper) {
    //app
    viper.SetDefault("app.name", "goo")
    viper.SetDefault("app.env", "local")
    viper.SetDefault("app.debug", true)

    //mysql
    viper.SetDefault("db.connection", "mysql")
    viper.SetDefault("db.host", "127.0.0.1")
    viper.SetDefault("db.port", "3306")
    viper.SetDefault("db.database", "default")
    viper.SetDefault("db.username", "default")
    viper.SetDefault("db.password", "123456")

    viper.SetDefault("cache.drive", "redis")

    //redis
    viper.SetDefault("redis.host", "127.0.0.1")
    viper.SetDefault("redis.port", "6379")
    viper.SetDefault("redis.auth", "")
    viper.SetDefault("redis.database", 0)

    //Queue
    viper.SetDefault("queue.drive", "kafka")

    //rabbitmq
    viper.SetDefault("rabbitmq.host", "127.0.0.1")
    viper.SetDefault("rabbitmq.port", "5672")
    viper.SetDefault("rabbitmq.vhost", "/")
    viper.SetDefault("rabbitmq.user", "admin")
    viper.SetDefault("rabbitmq.pass", "admin")

    //kafka
    viper.SetDefault("kafka.addr", "127.0.0.1:9092")

    //filesystem  cos,oss
    viper.SetDefault("fs.drive", "cos")

    viper.SetDefault("fs.bucket", "")
    viper.SetDefault("fs.region", "")
    viper.SetDefault("fs.secret_id", "")
    viper.SetDefault("fs.secret_key", "")

    //errors
    viper.SetDefault("error.report", "")


    //wechat-center rpc
    viper.SetDefault("wechat-center.addr", "127.0.0.1")
    viper.SetDefault("wechat-center.port", "8081")

}

func LoadConf() {
    viper.Set("db.connections", map[string]map[string]string{
        "mysql": {
            "driver":   "mysql",
            "host":     viper.GetString("db.host"),
            "port":     viper.GetString("db.port"),
            "database": viper.GetString("db.database"),
            "username": viper.GetString("db.username"),
            "password": viper.GetString("db.password"),
        },
    })

    viper.Set("filesystem", map[string]map[string]string{
        "cos": {
            "driver":     "cos",
            "bucket":     viper.GetString("fs.bucket"),
            "region":     viper.GetString("fs.region"),
            "secret_id":  viper.GetString("fs.secret_id"),
            "secret_key": viper.GetString("fs.secret_key"),
        },
    })

    viper.Set("caches", map[string]map[string]string{
        "redis": {
            "driver":   "redis",
            "host":     viper.GetString("redis.host"),
            "port":     viper.GetString("redis.port"),
            "database": viper.GetString("redis.database"),
            "auth":     viper.GetString("redis.auth"),
        },
        "sync": {
           "driver": "sync",
        },
    })

    viper.Set("queues", map[string]map[string]string{
        "rabbitmq": {
            "driver": "rabbitmq",
            "host":   viper.GetString("rabbitmq.host"),
            "port":   viper.GetString("rabbitmq.port"),
            "vhost":  viper.GetString("rabbitmq.vhost"),
            "user":   viper.GetString("rabbitmq.user"),
            "pass":   viper.GetString("rabbitmq.pass"),
        },
        "kafka": {
            "driver": "kafka",
            "addr":   viper.GetString("kafka.addr"),
        },
    })
}
