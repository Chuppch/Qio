# Docker 服务部署文档

## 快速启动

所有服务使用阿里云镜像源，按以下顺序启动：

### 1. 基础环境服务（必须）

```bash
cd docker-environment
docker-compose -f docker-compose-environment-aliyun.yml up -d
```

### 2. MCP 服务（可选）

```bash
cd docker-mcp
docker-compose -f docker-compose-mcp-aliyun.yml up -d
```

### 3. 监控服务（可选）

```bash
cd docker-monitoring
docker-compose -f docker-compose-grafana-aliyun.yml up -d
```

### 4. 日志服务（可选）

```bash
cd docker-logging
docker-compose -f docker-compose-elk-aliyun.yml up -d
```

---

## 服务端口说明

### 基础环境服务

| 端口 | 服务 | 说明 | 访问地址 |
|------|------|------|---------|
| 20001 | MySQL | 关系型数据库 | `localhost:20001` |
| 20002 | phpMyAdmin | MySQL 管理界面 | `http://localhost:20002` |
| 20003 | Redis | 内存数据库/缓存 | `localhost:20003` |
| 20004 | RedisAdmin | Redis 管理界面 | `http://localhost:20004` |
| 20005 | PostgreSQL | 向量数据库 | `localhost:20005` |
| 20006 | pgAdmin | PostgreSQL 管理界面 | `http://localhost:20006` |

**默认账号：**
- MySQL: `root` / `123456`
- phpMyAdmin: 使用 MySQL root 账号
- Redis: 无密码
- RedisAdmin: `admin` / `admin`
- PostgreSQL: `postgres` / `postgres`
- pgAdmin: `admin@qq.com` / `admin`

---

### MCP 服务

| 端口 | 服务 | 说明 | 访问地址 |
|------|------|------|---------|
| 8101 | CSDN MCP | CSDN 自动发帖服务 | `http://localhost:8101` |
| 8102 | 微信 MCP | 微信公众号通知服务 | `http://localhost:8102` |

---

### 监控服务

| 端口 | 服务 | 说明 | 访问地址 |
|------|------|------|---------|
| 3000 | Grafana | 监控面板（标准端口） | `http://localhost:3000` |
| 4000 | Grafana | 监控面板（备用端口） | `http://localhost:4000` |
| 8000 | Grafana MCP | Grafana MCP 服务 | `http://localhost:8000` |
| 9090 | Prometheus | 指标采集服务 | `http://localhost:9090` |
| 9100 | Node Exporter | 系统指标采集 | `http://localhost:9100` |

---

### 日志服务

| 端口 | 服务 | 说明 | 访问地址 |
|------|------|------|---------|
| 4560 | Logstash | 日志收集服务 | `localhost:4560` |
| 50000 | Logstash | 日志接收端口（TCP/UDP） | `localhost:50000` |
| 5601 | Kibana | 日志可视化界面 | `http://localhost:5601` |
| 9200 | Elasticsearch | 日志存储服务 | `http://localhost:9200` |
| 9300 | Elasticsearch | 集群通信端口 | `localhost:9300` |
| 9600 | Logstash | 监控端口 | `http://localhost:9600` |

---

## 服务依赖关系

```
基础环境服务
  ├─ MySQL (20001)
  │   └─ phpMyAdmin (20002)
  │
  ├─ Redis (20003)
  │   └─ RedisAdmin (20004)
  │
  └─ PostgreSQL (20005)
      └─ pgAdmin (20006)

MCP 服务（依赖基础环境网络）
  ├─ CSDN MCP (8101)
  └─ 微信 MCP (8102)

监控服务（独立）
  ├─ Prometheus (9090)
  ├─ Grafana (3000/4000)
  ├─ Grafana MCP (8000)
  └─ Node Exporter (9100)

日志服务（独立）
  ├─ Elasticsearch (9200/9300)
  ├─ Logstash (4560/50000/9600)
  └─ Kibana (5601)
```

---

## 停止服务

```bash
# 停止基础环境
cd docker-environment
docker-compose -f docker-compose-environment-aliyun.yml down

# 停止 MCP 服务
cd docker-mcp
docker-compose -f docker-compose-mcp-aliyun.yml down

# 停止监控服务
cd docker-monitoring
docker-compose -f docker-compose-grafana-aliyun.yml down

# 停止日志服务
cd docker-logging
docker-compose -f docker-compose-elk-aliyun.yml down
```

---

## 注意事项

1. **启动顺序**：必须先启动基础环境服务，再启动其他服务
2. **网络配置**：所有服务使用 `my-network` 网络，基础环境会自动创建
3. **端口占用**：确保端口未被占用，特别是 20001-20006 端口段
4. **应用配置**：应用连接数据库时需要使用新的端口（如 PostgreSQL: 20005）

