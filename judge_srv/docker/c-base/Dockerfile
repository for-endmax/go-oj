# 使用alpine:3.7作为基础镜像
FROM alpine:latest

# 替换Alpine的软件源为清华大学的镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

# 更新Alpine的包索引并安装基本的构建工具
RUN apk update && \
    apk upgrade && \
    apk add --no-cache build-base

# 在这里添加你的其他Dockerfile指令，例如COPY、RUN等

# 定义容器启动时执行的命令
CMD ["/bin/sh"]
