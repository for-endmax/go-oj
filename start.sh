#!/bin/bash

# 数组用于保存Go程序的PID
declare -a pids

# 函数用于关闭程序
cleanup() {
    echo "Cleaning up and shutting down..."

    # 向保存的PID发送Ctrl+C信号，捕获输出
    for pid in "${pids[@]}"; do
        output=$(kill -SIGINT "$pid" 2>&1)
        echo "关闭程序： PID $pid: $output"
    done

    # 等待一段时间，确保程序有足够的时间处理信号
    sleep 8

    exit 0
}
# 注册 cleanup 函数，以捕获 Ctrl+C 信号
trap cleanup SIGINT

# 进入 go-oj 目录
cd /home/endmax/go-oj/

# 启动 Jaeger Docker 容器
echo "启动容器..."
docker run -d --rm --name jaeger -p6831:6831/udp -p16686:16686 jaegertracing/all-in-one:latest

sleep 2

# 运行 init_config 程序，并保存PID
echo "向consul写配置..."
cd init_config
/usr/local/go/bin/go run main.go >/dev/null 2>&1 &
cd ..

echo "启动user_srv..."
cd user_srv
/usr/local/go/bin/go run main.go >/dev/null 2>&1 &
pids+=($!)
cd ..

echo "启动 record_srv..."
cd record_srv
/usr/local/go/bin/go run main.go >/dev/null 2>&1 &
pids+=($!)
cd ..

echo "启动 question_srv..."
cd question_srv
/usr/local/go/bin/go run main.go >/dev/null 2>&1 &
pids+=($!)
cd ..

sleep 5

echo "启动 user_web..."
cd user_web
/usr/local/go/bin/go run main.go >/dev/null 2>&1 &
pids+=($!)
cd ..

echo "启动 question_web..."
cd question_web
/usr/local/go/bin/go run main.go >/dev/null 2>&1 &
pids+=($!)
cd ..

echo "启动 submit_web..."
cd submit_web
/usr/local/go/bin/go run main.go >/dev/null 2>&1 &
pids+=($!)
cd ..

echo "启动 2个judge_srv..."
cd judge_srv
/usr/local/go/bin/go run main.go >/dev/null 2>&1  &
pids+=($!)

/usr/local/go/bin/go run main.go >/dev/null 2>&1  &
pids+=($!)

# 进入无限循环，保持脚本运行
while true; do
    # 在这里可以添加一些逻辑，以确保脚本一直在运行
    sleep 1
done
