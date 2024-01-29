# 启动项目demo
## 1.环境搭建

### mysql安装：
```shell
docker run \
-p 33061:3306  \
--restart=always \
--name mysql-5.7-1 \
--privileged=true  \
-v /root/mysql-5.7-1/conf:/etc/mysql/conf.d \
-v /root/mysql-5.7-1/data:/var/lib/mysql \
-v /root/mysql-5.7-1/log:/var/log/mysql \
-v /root/mysql-5.7-1/mysql-files:/var/lib/mysql-files \
-e MYSQL_ROOT_PASSWORD=admin123 \
-d mysql:5.7 \
--character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

### Consul安装：
```shell
docker run -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0
```

### Redis 安装：
```shell
docker run -p 6379:6379 -d redis:latest redis-server
```

### RabbitMQ 安装
```shell
docker run -d --hostname my-rabbit --name rabbit -p 15672:15672 -p 5672:5672 rabbitmq
docker exec -it 容器id /bin/bash
rabbitmq-plugins enable rabbitmq_management
```

### Jaeger安装
```shell
docker run -d --rm --name jaeger -p6831:6831/udp -p16686:16686 jaegertracing/all-in-one:latest
```

## 2.在MySQL中创建相应的数据库
运行xx_srv/build/main文件建表


## 3.构建docker镜像
使用 judge_srv/docker中的DockerFile

## 4.启动
 ./start.sh

