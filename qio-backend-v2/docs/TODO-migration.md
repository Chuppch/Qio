# v1 → v2 迁移待办

本文档记录从 v1（Java / Spring Boot / MyBatis-Plus）迁移到 v2（Go / Gin / GORM）过程中
发现的问题。

当前阶段的原则是**等价迁移**：只做语言与架构的转换，不改业务逻辑、不优化性能。
因此下列问题在 v2 代码中被如实保留，待语言迁移完成后再统一处理。

每条记录包含：问题所在、v1 的做法、v2 当前的处置、建议方向。

---

## 一、安全

### 1.1 密码使用裸 MD5

- **位置**：v1 `UserServiceImpl`
- **v1 做法**：密码经单次 MD5 摘要后入库，无加盐、无慢哈希
- **v2 处置**：保留，`user.User.PasswordHash` 仍存 MD5 摘要
- **建议**：改为 bcrypt 或 argon2id。需要一次密码迁移：登录成功时按旧算法校验、
  按新算法重写摘要，双写一段时间后移除旧分支

### 1.2 答题接口靠对称加密隐藏答案

- **位置**：v1 `QuestionServiceImpl.allAnswerToFront` / `decode`
- **v1 做法**：整套题目（含 `correct_answer`）用 AES 加密后下发前端，前端解密判题；
  密钥来自配置项 `secretKey`
- **v2 处置**：尚未迁移该接口，仓储层只提供带答案与不带答案两种查询
- **建议**：判题留在服务端。`startAnswer` 已经在做列投影不下发答案，
  `allAnswerToFront` 与之矛盾，二者应择一

---

## 二、数据库结构

### 2.1 `create_user` / `update_user` 类型不统一

- **位置**：全库
- **现状**：`letter`、`user`、`friend`、`friend_request` 是 `bigint`；
  `bottle`、`questions`、`question_user_status` 是 `varchar(50)`。存的都是用户 ID
- **v2 处置**：`internal/infrastructure/mysql/base.go` 拆成 `auditFields` 与
  `auditFieldsStrUser` 两个变体如实映射
- **建议**：统一为 `bigint`，合并两个变体

### 2.2 `question_user_status` 把进度平铺成十列

- **位置**：`question_user_status.question_set_1_id` ~ `question_set_10_id`
- **v1 做法**：十个 `int` 列。写入用 `switch(setId)` 逐个 `setQuestionSetN(1)`，
  读取用反射拼 `"getQuestionSet" + i` 方法名
- **附带问题**：列名叫 `..._id`，但实际存的是 0 未完成 / 1 已完成的状态值，并非 ID
- **v2 处置**：PO 如实映射十列，在 `po_explore.go` 中折叠为 `map[int64]int`；
  `explore.SetCount` 常量记录「套数被固定在表结构里」这一约束
- **建议**：改为 `(user_id, set_id, state)` 纵表，增删题库不再需要改表，
  `SetCount` 常量随之移除

### 2.3 `paper` 表无主键

- **位置**：`paper` 表
- **v2 处置**：PO 如实映射
- **建议**：补主键。无主键会影响 GORM 的按主键更新、以及主从复制的行定位

### 2.4 用数值语义的字段存成 varchar

- **位置**：`function_card.speed_rate`、`function_card.reduce_time`
- **v1 做法**：列类型是 varchar，存的是数值文本
- **v2 处置**：`po_shop.go` 中用 `parseFloatDefault` 解析，解析失败取零值
- **建议**：改为 `decimal` / `int`，去掉解析环节

### 2.5 布尔语义的字段存成字符串

- **位置**：`user.address` JSON 列中的 `isDefault`
- **v1 做法**：JSON 里存字符串 `"true"` / `"false"`，而非 JSON 布尔
- **v2 处置**：`json.go` 的 `decodeAddress` 兼容解析
- **建议**：改为 JSON 布尔值，或随地址簿一起从 JSON 列改为关联表

---

## 三、并发与一致性

### 3.1 余额读—改—写

- **位置**：v1 `UserServiceImpl`、`GameServiceImpl.winFfl` 等多处
- **v1 做法**：`selectById` 取出 User，`setMoney(getMoney() + n)`，再 `updateById`
- **v2 处置**：`user.Repository.UpdateMoney` 已改为 `gorm.Expr("money + ?")`，
  这是本轮迁移中少数偏离等价原则的地方——读改写会静默丢失并发的加减
- **建议**：其余走余额的路径统一收敛到 `UpdateMoney`

### 3.2 抽奖次数扣减非原子

- **位置**：v1 `GameServiceImpl.winFfl`
- **v1 做法**：Redis `GET` → 判断 `> 0` → `SET limit-1`
- **v2 处置**：`redis/repo_explore_draw.go` 的 `Consume` 保留同样的三步
- **建议**：改为 Lua 脚本或 `DECR` + 负值回滚，使判断与扣减原子化

### 3.3 每日任务状态更新读写整个 JSON 数组

- **位置**：v1 任务模块
- **v1 做法**：个人任务副本以 JSON 数组整体存 Redis，改一条任务的状态需要
  取出整个数组、改完再整体写回
- **v2 处置**：`redis/repo_user_task.go` 的 `setTaskStatus` 保留该做法
- **建议**：改为 Hash，一个 field 一条任务，单条更新不再影响其他任务

---

## 四、查询与性能

### 4.1 可捞漂流瓶全表扫描

- **位置**：v1 `BottleServiceImpl.getNotIsPickedBottles`
- **v1 做法**：把所有可捞的瓶子全量查出，在应用层 `Collections.shuffle` 取一个
- **v2 处置**：`mysql/repo_bottle.go` 的 `ListAvailable` 保留全量查询，
  随机选择留在 service
- **建议**：改为 `ORDER BY RAND() LIMIT 1`，或按主键区间随机定位

### 4.2 `ListAvailable` 对 NULL 的处理待确认

- **位置**：`mysql/repo_bottle.go` 的 `ListAvailable`
- **v1 做法**：用 MyBatis-Plus 的 `notIn("update_user", ...)`，生成
  `update_user NOT IN (...)`。MySQL 三值逻辑下该条件会**过滤掉** `update_user IS NULL`
  的行，即从未被捞过的瓶子反而捞不到
- **v2 处置**：当前写成 `update_user IS NULL OR update_user <> ?`，
  会把 NULL 行纳入结果——**这是行为改变，不是性能优化**
- **待决策**：
  - A：严格照搬 v1，只用 `update_user <> ?`
  - B：保持现状，视 v1 为缺陷并在此处修正
- **状态**：等待确认

### 4.3 查询方法内做写操作

- **位置**：v1 `UserServiceImpl.getMyFriends`
- **v1 做法**：名为 get 的查询方法内部执行了更新
- **v2 处置**：`friend` 域尚未迁移
- **建议**：迁移时把写操作剥离到独立方法，由调用方显式触发

---

## 五、硬编码与死代码

### 5.1 注册赠品的道具 ID 写死

- **位置**：v1 `UserServiceImpl` 注册流程
- **v1 做法**：注册时赠送 7 件道具，ID 直接写在代码里；其中一处是
  `selectById(0L)`——依赖库里真实存在一条 `id = 0` 的记录
- **v2 处置**：按要求**未迁移**该逻辑
- **附带影响**：`shop.Repository.FindFunctionCard` 用 `Where("id = ?", id)`
  而非 `First(&po, id)`，因为 GORM 的主键简写会忽略零值，查不到 `id = 0` 的功能卡
- **建议**：改为可配置的新人礼包，并清理 `id = 0` 这条数据

### 5.2 抽奖次数与奖励金额写死

- **位置**：v1 `GameServiceImpl`
- **v1 做法**：每日 10 次、每次奖励 10 猪仔钱，均为字面量
- **v2 处置**：提为 `explore.DefaultDrawLimit` 与 `explore.DrawReward` 两个常量
- **建议**：改为运营可配

### 5.3 抽奖次数的过期时间会顺延

- **位置**：v1 `GameServiceImpl`
- **v1 做法**：每次读或写都以 `Duration.ofDays(1)` 重置 TTL，因此次数窗口是
  「自最后一次操作起 24 小时」，而非自然日归零
- **v2 处置**：`drawLimitTTL` 保留 24 小时，且在每次写入时重置
- **建议**：若产品语义是「每日 10 次」，应改为按自然日的键（如
  `game:ffl:user:<id>:<yyyyMMdd>`）或对齐到次日零点的 TTL

### 5.4 `QuestionServiceImpl` 存在大段注释掉的死代码

- **位置**：v1 `QuestionServiceImpl` 第 41–149、183–207 行
- **现状**：`genQuestion` 整个方法及两段进度读取逻辑被注释掉，其中的
  Mapper 调用容易被误读为生效逻辑
- **v2 处置**：未迁移。`explore.Repository` 的方法集只覆盖生效的
  `userLoginPage`、`startAnswer`、`submitAnswers`、`allAnswerToFront`
- **建议**：v1 侧直接删除

### 5.5 信件版式常量命名与取值不符

- **位置**：v1 信件模块的 `letterType` 常量
- **现状**：横排与竖排两个常量的名称与实际取值对应关系是反的
- **v2 处置**：`letter` 域尚未迁移
- **建议**：迁移 `letter` 域时修正命名，注意存量数据按取值而非名称解读

---

## 六、领域建模

### 6.1 持久化实体反向依赖视图对象

- **位置**：v1 `User` 实体
- **v1 做法**：用户拥有的道具以 `List<XxxVO>` 存进 JSON 列，实体因此依赖 VO
- **v2 处置**：`user.User.OwnedItems` 只表达「拥有哪些道具」，落库形态由仓储层决定
- **建议**：无需额外处理，迁移时已解决

### 6.2 用户名长度约束自相矛盾

- **位置**：v1 `UserConstants.USERNAME_MIN_LENGTH` 与 `AccountValidator`
- **现状**：常量声明最短 6 位，正则实际允许 4 位（首字母 + 3 位）
- **v2 处置**：`user/validation.go` 以正则为准统一为 4 位，并把常量与正则收在一处
- **建议**：确认产品预期长度后统一

### 6.3 `ItemType` 在 `user` 与 `shop` 两域重复定义

- **位置**：`internal/domain/user/model.go`、`internal/domain/shop/model.go`
- **现状**：取值相同、语义相同，两处各定义一份
- **v2 处置**：暂未处理
- **建议**：抽到共享位置。注意 `Address` 是另一种情况——两域结构相同但语义不同
  （资料 vs 投递凭据），刻意不共享

### 6.4 `questions` 表四个选项列

- **位置**：`questions.option_a` ~ `option_d`
- **现状**：同一件东西的四份拷贝，任何遍历都要写死四次
- **v2 处置**：`po_explore.go` 折叠为 `map[explore.Option]string`
- **建议**：若要支持选项数量可变，需改为子表；当前四选一场景下折叠已足够
