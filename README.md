# goo

微服务框架

1. HTTP/RPC
2. Rabbitmq/Kafka队列
3. Redis缓存
4. 事件驱动
5. 日志告警
6. 文件系统
7. 限流熔断降级
8. 链路监控
9. k8s/istio部署方便


启动步骤
```go
1. 拷贝config.example.toml至config.toml 根据需求修改其内容
2. 拷贝docker-compose.example.yml至docker-compose.yml 根据需求修改其内容
3. 运行docker-compose pull && docker-compose up -d启动环境
4. 运行make api启动服务
3. 关闭环境请运行 docker-compose down
```


目录说明
```go
|____build                      编译所需文件
|____cmd                        入口
| |____api                      http/rpc
| |____queue                    队列消费
| |____task                     任务脚本
|____config                     默认配置
|____google                     GRPC所需文件
|____internal                   内部文件
| |____client                   任务脚本
| |____errors                   错误定义
| |____events                   事件处理
| |____handler                  处理程序
| |____job                      生产队列
| |____repo                     repository层
| |____server                   服务入口
| |____svc                      context
| |____types                    类型定义
| |____utils                    辅助方法
|____pkg                        外部包
| |____cache                    缓存
| |____client                   客户端脚本任务
| |____db                       数据库
| |____event                    事件
| |____file                     文件系统
| |____grpcgw                   grpc-gateway
| |____http                     echo/gin
| |____log                      日志
| |____pool                     池
| |____queue                    队列
|____proto                      protobuf文件
|____storage                    存储


```

