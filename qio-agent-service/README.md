# Qio Agent Service

侨缘信使智能中台后端，负责模型接入、Prompt、RAG、MCP 工具和 Agent 执行编排。

<p align="center">
	<!-- <img alt="logo" src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/logo.png"> -->
	<img alt="logo" src="https://raw.githubusercontent.com/Chuppch/qiaopi-master-frontend/refs/heads/master/docs/image/主页面照片.png">
</p>
<h4 align="center">跨越四海，侨缘线牵——侨缘信使，让世界没有距离。</h4>
<p align="center">
	<a href="https://github.com/Chuppch/qiaopi-master-frontend"><img src="https://img.shields.io/badge/%E5%89%8D%E7%AB%AF%E5%B7%A5%E7%A8%8B-github?logo=github&label=github&color=%23181717"></a>
    <a href="https://github.com/Chuppch/Qiaopi-master"><img src="https://img.shields.io/badge/%E5%90%8E%E7%AB%AF%E5%B7%A5%E7%A8%8B-github?logo=github&label=github&color=%23181717"></a>
    <a href="https://github.com/Chuppch/qio-agent-service"><img src="https://img.shields.io/badge/Agent%E5%90%8E%E7%AB%AF%E5%B7%A5%E7%A8%8B-github?logo=github&label=github&color=%23181717"></a>
	<a href="https://github.com/Chuppch/Qiaopi-master"><img src="https://img.shields.io/badge/Qiaopi-v1.0.1-brightgreen.svg"></a>
	<a href="https://github.com/Chuppch/Qiaopi-master?tab=MIT-1-ov-file"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
</p>



## 项目介绍

《**侨缘信使**》是一个旨在宣传和传承侨批文化的互动网站。文化内容的数字化呈现、互动性与教育性的结合、情感共鸣的构建以及用户参与的文化共创。我们致力于通过网站，让更多人了解并参与到侨批文化的保护与传承中来，同时为文化带来新的活力。学习侨批文化、体验写侨批、收侨批和漂流瓶等功能，感受慢信文化的魅力，"跨越四海，侨缘线牵——侨缘信使，让世界没有距离。"

平台还引入了智能交互能力作为支撑，使 AI 能够在不同文化语境下参与用户的创作与交流，在情感引导、文化背景补充与内容辅助等方面提供恰当支持，提升整体互动的连贯性与沉浸感。同时，该能力也被应用于平台的管理端，用于辅助内容整理与日常运营，让使平台在保持文化温度的同时更加高效、灵活。



<p align = "right">—— 五灵威力小队</p>

**技术栈**：***Java 17、Spring Boot 3.4.3、Spring AI、MyBatis、MySQL、PostgreSQL (pgvector)、Redis、Maven、Docker***

## **在线体验**

**[侨缘信使🎉](http://qiaoyuanxinshi.com/)**

## 演示图

<table>
    <tr>
        <td align="center" width="50%">
            <img src="docs/dev-ops/image/对齐.gif" alt="对齐功能演示" style="max-width: 100%; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"/>
            <br/><small><b>对齐功能</b></small>
        </td>
        <td align="center" width="50%">
            <img src="docs/dev-ops/image/copy.gif" alt="复制功能演示" style="max-width: 100%; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"/>
            <br/><small><b>快捷键功能</b></small>
        </td>
    </tr>
    <tr>
        <td align="center" colspan="2">
            <img src="docs/dev-ops/image/高亮动画.gif" alt="高亮动画演示" style="max-width: 100%; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"/>
            <br/><small><b>高亮动画</b></small>
        </td>
    </tr>
    <tr>
        <td align="center" width="50%">
            <img src="docs/dev-ops/image/agent列表.png" alt="Agent列表" style="max-width: 100%; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"/>
            <br/><small><b>Agent列表</b></small>
        </td>
        <td align="center" width="50%">
            <img src="docs/dev-ops/image/配置页面.png" alt="配置页面" style="max-width: 100%; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"/>
            <br/><small><b>配置页面</b></small>
        </td>
    </tr>
</table>

## 快速部署

### 环境要求

- **JDK**: 17 或更高版本
- **Maven**: 3.6+ 
- **MySQL**: 8.0+ (业务数据库)
- **PostgreSQL**: 14+ with pgvector (向量数据库，用于 RAG)
- **Redis**: 6.0+ (缓存)
- **Docker**: 20.10+ (可选，用于容器化部署)

### 数据库初始化

1. **MySQL 数据库初始化**

   执行 SQL 脚本创建数据库和表结构：

   ```bash
   # 初始化数据库（标准版本）
   mysql -u root -p < docs/dev-ops/mysql/sql/qio-agent-init.sql
   
   # 或使用滑动窗口版本（优化 token 消耗）
   mysql -u root -p < docs/dev-ops/mysql/sql/qio-agent-init-window.sql
   ```

2. **PostgreSQL 向量数据库初始化**

   ```bash
   # 初始化 pgvector 扩展和表结构
   psql -U postgres -d ai-rag-knowledge -f docs/dev-ops/pgvector/sql/init.sql
   ```

### 配置文件

项目使用多环境配置，主要配置文件位于 `qio-agent-app/src/main/resources/`：

- `application.yml` - 主配置文件
- `application-dev.yml` - 开发环境配置
- `application-test.yml` - 测试环境配置
- `application-prod.yml` - 生产环境配置

#### 数据库配置

在 `application-dev.yml` 中配置数据库连接信息：

```yaml
spring:
  datasource:
    # MySQL 配置（业务数据库）
    mysql:
      url: jdbc:mysql://localhost:20001/ai-agent-station?useUnicode=true&characterEncoding=UTF-8&useSSL=false&serverTimezone=Asia/Shanghai
      username: root
      password: your-mysql-password  # 替换为您的 MySQL 密码
    # PostgreSQL 配置（向量数据库）
    pgvector:
      url: jdbc:postgresql://localhost:20005/ai-rag-knowledge?tcpKeepAlive=true&socketTimeout=60
      username: postgres
      password: your-postgres-password  # 替换为您的 PostgreSQL 密码
  # Redis 配置
  data:
    redis:
      host: localhost
      port: 20003
      password: # Redis 密码（如无密码可留空）
```

#### Spring AI 配置

配置 AI 模型服务（支持阿里云百炼、OpenAI 等）：

```yaml
spring:
  ai:
    openai:
      base-url: https://dashscope.aliyuncs.com/compatible-mode  # 阿里云百炼
      api-key: your-api-key-here  # 替换为您的 AI 服务 API Key
      embedding:
        options:
          model: text-embedding-v4
          dimensions: 1536
```

## 未完待续

- [ ] 完善 API 文档和接口说明
- [ ] 添加单元测试和集成测试
- [ ] 优化数据库连接池配置，提升并发性能
- [ ] 完善监控和告警机制（Prometheus + Grafana）
- [ ] 添加分布式链路追踪（SkyWalking/Jaeger）
- [ ] 优化 RAG 服务检索性能，支持更大规模知识库
- [ ] 完善滑动窗口机制，进一步优化 token 消耗
- [ ] 添加多租户支持和权限管理
- [ ] 完善日志收集和分析（ELK Stack）
