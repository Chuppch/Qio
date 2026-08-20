# Qio Agent Console

侨缘信使智能中台的内部配置与运营控制台。

<p align="center">
	<img alt="logo" src="../public/brand/qio-hero.png">
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
            <img src="../public/previews/agent-hub/align.gif" alt="对齐功能演示" style="max-width: 100%; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"/>
            <br/><small><b>对齐功能</b></small>
        </td>
        <td align="center" width="50%">
            <img src="../public/previews/agent-hub/copy.gif" alt="复制功能演示" style="max-width: 100%; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"/>
            <br/><small><b>快捷键功能</b></small>
        </td>
    </tr>
    <tr>
        <td align="center" colspan="2">
            <img src="../public/previews/agent-hub/highlight.gif" alt="高亮动画演示" style="max-width: 100%; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"/>
            <br/><small><b>高亮动画</b></small>
        </td>
    </tr>
    <tr>
        <td align="center" width="50%">
            <img src="../public/previews/agent-hub/agent-list.png" alt="Agent列表" style="max-width: 100%; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"/>
            <br/><small><b>Agent列表</b></small>
        </td>
        <td align="center" width="50%">
            <img src="../public/previews/agent-hub/configuration.png" alt="配置页面" style="max-width: 100%; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);"/>
            <br/><small><b>配置页面</b></small>
        </td>
    </tr>
</table>


## 快速部署

### 环境要求

- **Node.js**: 20 或更高版本
- **npm/yarn**: 包管理器（npm 9+ 或 yarn 1.22+）
- **Docker**: 20.10+ (可选，用于容器化部署)

### 启动命令

#### 开发环境

```bash
# 安装依赖
npm install

# 启动开发服务器（会自动打开浏览器）
npm run dev
# 或
npm start
```

开发服务器默认运行在 `http://localhost:3000`（具体端口以实际输出为准）

#### 生产构建

```bash
# 构建生产版本
npm run build

# 构建产物位于 dist/ 目录
```

#### Docker 部署

```bash
# 使用 Docker Compose 启动（推荐）
docker-compose up -d

# 或使用 Dockerfile 构建镜像
docker build -t qio-agent-console:latest .
docker run -d -p 3002:3002 qio-agent-console:latest
```

生产环境默认运行在 `http://localhost:3002`

### 其他命令

```bash
# 代码检查
npm run lint

# 自动修复代码格式问题
npm run lint:fix

# 清理构建产物
npm run clean
```
