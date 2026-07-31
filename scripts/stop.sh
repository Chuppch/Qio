#!/bin/bash
# 一键停止开发环境（前端 + 后端 + 容器 + Colima）

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "🌐 停止前端..."
pkill -f "vue-cli-service serve" 2>/dev/null

echo "☕ 停止后端..."
pkill -f "spring-boot:run" 2>/dev/null
pkill -f "QiaopiServerApplication" 2>/dev/null

echo "🐳 停止中间件容器..."
cd "$PROJECT_ROOT" && docker compose down

echo "💤 停止 Colima..."
colima stop

echo ""
echo "✅ 开发环境已停止，内存已释放"
