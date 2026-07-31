#!/bin/bash
# 一键启动开发环境（Colima + 中间件 + 后端 + 前端）

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "启动 Colima..."
colima start --cpu 2 --memory 4

echo "启动中间件容器..."
cd "$PROJECT_ROOT" && docker compose up -d

echo "等待 MySQL 就绪..."
sleep 8

echo "启动后端（端口 3031）..."
cd "$PROJECT_ROOT/qio-backend" && mvn spring-boot:run -pl Qiaopi-server &
BACKEND_PID=$!

echo "等待后端启动..."
sleep 15

echo "启动前端（端口 3010）..."
cd "$PROJECT_ROOT/qio-frontend" && npm run serve &
FRONTEND_PID=$!

echo ""
echo "开发环境已启动"
echo ""
echo "  前端: http://localhost:3010"
echo "  后端: http://localhost:3031"
echo "  RabbitMQ 管理面板: http://localhost:15672 (root/root)"
echo "  Minio 控制台: http://localhost:9001 (minioadmin/minioadmin)"
echo ""
echo "  后端 PID: $BACKEND_PID"
echo "  前端 PID: $FRONTEND_PID"
echo ""
echo "  停止所有服务: ./scripts/stop.sh"

wait
