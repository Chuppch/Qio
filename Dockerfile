# ============================================================
# Qio - Docker Compose 开发环境
# 使用方式: docker compose up -d
# ============================================================

version: '3.8'

services:
  # --------------------------------------------------------
  # MySQL 8.3
  # --------------------------------------------------------
  mysql:
    image: mysql:8.4
    container_name: qio-mysql
    restart: unless-stopped
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: qiaopi
      TZ: Asia/Shanghai
    volumes:
      - qio-mysql-data:/var/lib/mysql
    command: --default-authentication-plugin=mysql_native_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

  # --------------------------------------------------------
  # Redis 7
  # --------------------------------------------------------
  redis:
    image: redis:8-alpine
    container_name: qio-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    command: redis-server --requirepass root
    volumes:
      - qio-redis-data:/data

  # --------------------------------------------------------
  # RabbitMQ 3 (with management plugin)
  # --------------------------------------------------------
  rabbitmq:
    image: rabbitmq:4-management
    container_name: qio-rabbitmq
    restart: unless-stopped
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: root
      RABBITMQ_DEFAULT_PASS: root
      RABBITMQ_DEFAULT_VHOST: /

  # --------------------------------------------------------
  # Minio (S3-compatible object storage)
  # --------------------------------------------------------
  minio:
    image: minio/minio
    container_name: qio-minio
    restart: unless-stopped
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - qio-minio-data:/data
    command: server /data --console-address ":9001"

volumes:
  qio-mysql-data:
  qio-redis-data:
  qio-minio-data:
