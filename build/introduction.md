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

### consul安装：
```shell
docker run -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0
```

### redis安装：
```shell
docker run -p 6379:6379 -d redis:latest redis-server
```

## 2.在MySQL中创建数据库

以交互方式进入容器，启动shell
```shell
docker exec -it mysql-5.7-1 /bin/bash
```
登录mysql
```shell
mysql -uroot -padmin123
```
创建数据库
```mysql
CREATE DATABASE `go-oj_user_srv` DEFAULT CHARACTER SET utf8mb4
CREATE DATABASE `go-oj_question_srv` DEFAULT CHARACTER SET utf8mb4
```
建表
> 运行每个模块中test/build文件

## 3.向consul写入远程配置
```shell
cd build
go run build.go
```

## 4.依次启动srv服务与web服务

