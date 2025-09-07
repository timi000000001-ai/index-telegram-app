#!/bin/bash

# 优化后的调试脚本，用于重启 management-service 和 bot-service。

echo "--- 正在停止服务 ---"
fuser -k 8090/tcp 2>/dev/null
fuser -k 8081/tcp 2>/dev/null
echo "端口 8090 和 8081 上的服务已停止。"

sleep 2

PROJECT_ROOT="/root/新建文件夹/index-telegram-app"

# 清理旧的日志文件
rm -f "$PROJECT_ROOT/telegram-bot-services/management-service/management.log"
rm -f "$PROJECT_ROOT/telegram-bot-services/bot-service/bot.log"

echo ""
echo "--- 正在启动 management-service (端口 8090) ---"
cd "$PROJECT_ROOT/telegram-bot-services/management-service"
# 使用 --http 标志指定端口
nohup /root/sdk/go1.25.0/bin/go run . serve --dev --dir pb_data --http="127.0.0.1:8090" > management.log 2>&1 &
sleep 2

echo ""
echo "--- 正在启动 bot-service (端口 8081) ---"
cd "$PROJECT_ROOT/telegram-bot-services/bot-service"
nohup /root/sdk/go1.25.0/bin/go run ./cmd/bot/main.go > bot.log 2>&1 &
sleep 2

echo ""
echo "--- 服务已在后台启动 ---"

# 检查是否提供了 '--logs' 参数
if [ "$1" == "--logs" ]; then
    echo "--- 正在输出日志 ---"
    trap 'kill $(jobs -p)' EXIT
    tail -f "$PROJECT_ROOT/telegram-bot-services/management-service/management.log" &
    tail -f "$PROJECT_ROOT/telegram-bot-services/bot-service/bot.log" &
    wait
else
    echo "你可以使用 'tail -f /root/新建文件夹/index-telegram-app/telegram-bot-services/management-service/management.log' 查看 management-service 日志"
    echo "你可以使用 'tail -f /root/新建文件夹/index-telegram-app/telegram-bot-services/bot-service/bot.log' 查看 bot-service 日志"
fi  