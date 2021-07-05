# goo

微服务框架

1. HTTP/RPC
2. Rabbitmq/Kafka队列
3. Redis缓存
4. 事件驱动
5. 监控日志
6. 文件系统
7. 链路监控
8. k8s/istio部署方便


启动步骤
```go
1. 拷贝config.example.toml至config.toml 根据需求修改其内容
2. 拷贝docker-compose.example.yml至docker-compose.yml 根据需求修改其内容
3. 运行`docker-compose pull && docker-compose up -d`启动环境
4. 运行`make api` 启动环境
3. 关闭环境请运行`docker-compose down`
```


目录说明
