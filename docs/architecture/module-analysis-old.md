# Qio v1.x 模块分析

基于现有 qio-frontend 和 qio-backend 代码梳理，作为 v2 重构的参考基线。

## 前端模块（qio-frontend）

代码规模：26 个 Vue 文件 + 19 个 JS 文件，约 17,800 行代码。

### 1. 首页展示（index）

| 文件 | 功能 |
|------|------|
| `views/index/index.vue` | 主页 |
| `views/index/introduce.vue` | 侨批文化介绍 |
| `views/index/letter.vue` | 信件入口 |
| `views/index/game.vue` | 游戏入口 |
| `views/index/shop.vue` | 商店入口 |
| `views/index/marketing.vue` | 营销活动入口 |

### 2. 信件（letter）

| 文件 | 功能 |
|------|------|
| `views/letter/WriteLetter.vue` | 写信 |
| `views/letter/ReceiveLetter.vue` | 收信 |
| `views/letter/DriftingBottle.vue` | 漂流瓶 |
| `api/letter.js` | 信件接口 |
| `api/drifting.js` | 漂流瓶接口 |

### 3. AI 对话（ai）

| 文件 | 功能 |
|------|------|
| `views/ai/Qiaobao.vue` | AI 聊天主页面 |
| `views/ai/plan1.vue` | AI 方案 1 |
| `views/ai/plan2.vue` | AI 方案 2 |
| `views/ai/planSucess.vue` | 方案成功页 |
| `views/chat.vue` | WebSocket 对话 |

### 4. 游戏（game）

| 文件 | 功能 |
|------|------|
| `views/game/explore.vue` | 探索游戏 |
| `views/game/question.vue` | 答题 |
| `views/game/know.vue` | 知识问答 |
| `views/game/fanfanle.vue` | 翻翻乐 |
| `views/game/import.vue` | 游戏入口 |
| `api/fanfanle.js` | 翻翻乐接口 |
| `api/know.js` | 知识问答接口 |

### 5. 商店 & 营销（shop / marketing）

| 文件 | 功能 |
|------|------|
| `api/shop.js` | 商店接口 |
| `api/card.js` | 明信片/卡片接口 |
| `api/marketing.js` | 营销接口 |

### 6. 用户（user）

| 文件 | 功能 |
|------|------|
| `views/user/userLogin.vue` | 登录 |
| `views/user/register.vue` | 注册 |
| `views/user/forgetCode.vue` | 忘记密码 |
| `views/profile/index.vue` | 个人中心 |
| `api/user.js` | 用户接口 |
| `store/modules/user.js` | 用户状态 |

### 7. 基础设施（infra）

| 文件 | 功能 |
|------|------|
| `utils/request.js` | Axios 封装 |
| `utils/auth.js` | Token 管理 |
| `utils/qiaopi.js` | 业务工具函数 |
| `utils/errorCode.js` | 错误码映射 |
| `plugins/cache.js` | 缓存 |
| `plugins/canvas-trail-plugin.js` | Canvas 轨迹动画插件 |
| `router/index.js` | 路由 |
| `store/index.js` | 状态管理入口 |
| `store/money.js` | 金币/货币状态 |

---

## 后端模块（qio-backend）

代码规模：196 个 Java 文件，约 17,400 行代码。

### 1. 用户（User）

| 层 | 文件 |
|----|------|
| Controller | `UserController` |
| Service | `UserService` / `UserServiceImpl` |
| Mapper | `UserMapper`, `AvatarMapper` |
| Entity | `User`, `Avatar`, `UserStatistics`, `UserSignAward` |
| DTO | `UserLoginDTO`, `UserRegisterDTO`, `UserResetPasswordDTO`, `UserUpdateDTO` |

### 2. 信件（Letter）

| 层 | 文件 |
|----|------|
| Controller | `LetterController` |
| Service | `LetterService` / `LetterServiceImpl`, `G2dService` / `G2dServiceImpl` |
| Mapper | `LetterMapper` |
| Handler | `LetterSocketHandler`（WebSocket 推送信件状态） |
| Task | `LetterTask`（定时任务，信件送达计时） |
| Entity | `Letter`, `Address`, `Country` |

### 3. 漂流瓶（Bottle）

| 层 | 文件 |
|----|------|
| Controller | `BottleController` |
| Service | `BottleService` / `BottleServiceImpl` |
| Mapper | `BottleMapper` |
| Entity | `Bottle` |

### 4. AI 对话（Chat）

| 层 | 文件 |
|----|------|
| Service | `ChatService` / `ChatServiceImpl` |
| Handler | `ChatSocketHandler`（WebSocket 流式输出） |
| Listener | `AiInteractiveListener`（RabbitMQ 消费 AI 消息） |
| Task | `AiTask`（定时任务） |
| Pojo | `AiData`, `AiInteractData`, `ChatCompletionRequest`, `ChatTool`, `Choice`, `Delta`, `MyModelData`, `Usage`, `WebSearch`, `WebSearchResponse` |

### 5. 游戏 & 答题（Game / Question）

| 层 | 文件 |
|----|------|
| Controller | `GameController`, `QuestionController` |
| Service | `GameService` / `GameServiceImpl`, `QuestionService` / `QuestionServiceImpl` |
| Mapper | `QuestionsMapper`, `QuestionUserStatusMapper` |
| Entity | `Questions`, `QuestionUserStatus` |

### 6. 商店 & 营销（Shop / Marketing）

| 层 | 文件 |
|----|------|
| Controller | `CardController`, `MarketingController` |
| Service | `CardService` / `CardServiceImpl`, `MarketingService` / `MarketingServiceImpl` |
| Mapper | `CardMapper`, `MarketingMapper`, `SignetMapper` |
| Entity | `Commodity`, `FunctionCard`, `Signet`, `TaskTable` |
| Task | `SignTask`（签到任务） |

### 7. 字体 & 信纸（Font / Paper）

| 层 | 文件 |
|----|------|
| Controller | `FontController`, `PaperController` |
| Service | `FontService` / `FontServiceImpl`, `PaperService` / `PaperServiceImpl` |
| Mapper | `FontMapper`, `FontColorMapper`, `FontPaperMapper`, `PaperMapper` |
| Entity | `Font`, `FontColor`, `FontPaper`, `Paper` |

### 8. 好友（Friend）

| 层 | 文件 |
|----|------|
| Service | `FriendService` / `FriendServiceImpl` |
| Mapper | `FriendMapper`, `FriendRequestMapper` |
| Entity | `Friend`, `FriendRequest` |

### 9. 基础设施（Infra）

| 类别 | 文件 |
|------|------|
| 配置 | `AppRedisConfig`, `MqConfig`, `MybatisConfig`, `WebConfig`, `WebMvcConfiguration`, `WebSocketConfig` |
| 拦截器 | `JwtTokenAdminInterceptor`, `JwtTokenUserInterceptor`, `UserHandshakeInterceptor` |
| 异常处理 | `GlobalExceptionHandler` + 各业务异常类（15 个） |
| 工具类 | `JwtUtil`, `AESUtil`, `CacheClient`, `RedisUtils`, `IpUtils`, `StringUtils` 等 |
| 通用 | `Result`, `PageResult`, `BaseEntity`, `UserContext` |

---

## 统计概览

| 维度 | 前端 | 后端 |
|------|------|------|
| 文件数 | 45 | 196 |
| 代码行数 | ~17,800 | ~17,400 |
| 业务模块 | 7 | 9 |
| 页面/Controller | 26 个 Vue | 12 个 Controller |
| API/Service | 8 个 API 文件 | 12 个 Service |
| 数据模型 | — | 20 个 Entity |
