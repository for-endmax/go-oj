# go-oj
这是一个使用go语言编写的在线oj平台（后端部分）


## 主要技术栈
- Gin
- GRPC
- GORM
- Redis
- MySQL
- RabbitMQ

## 功能
### 用户模块
普通用户 注册、登录、修改自己的信息

管理员 添加普通用户、修改用户的信息、查看用户列表

> 使用Validator进行参数的校验，JWT控制登录状态
### 题目模块
普通用户 查看题目列表，查看特定题目信息，查看测试数据
管理员   增删改查题目和测试信息
## 预想架构

![架构图](pics/img.png)


## 快速启动demo
[快速启动](build/introduction.md)