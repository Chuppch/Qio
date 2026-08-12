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

### 4.2 可捞漂流瓶查询会排除 NULL 行

- **位置**：v1 `BottleServiceImpl.getNotIsPickedBottles`
- **v1 做法**：
  ```java
  new LambdaQueryWrapper<Bottle>()
      .eq(Bottle::isPicked, false)
      .notIn(Bottle::getUserId, UserContext.getUserId())
      .notIn(Bottle::getUpdateUser, UserContext.getUserId());
  ```
  两个 `notIn` 生成 `NOT IN (?)`。MySQL 三值逻辑下该条件对 NULL 行求值为 NULL，
  因此 `user_id` 或 `update_user` 为 NULL 的瓶子会被**静默排除**，永远捞不到
- **实际影响范围**：有限。`Bottle extends BaseEntity`，`update_user` 带
  `@TableField(fill = FieldFill.INSERT_UPDATE)`，`insertFill` 在插入时会把它填成
  投放者 ID，所以经应用投放的瓶子 `update_user` 不为 NULL。NULL 只出现在
  非应用途径写入的数据里
- **v2 处置**：照搬，用 `update_user <> ?`（与 `NOT IN (?)` 的 NULL 语义等价）。
  判断这属于 v1 的数据逻辑缺陷，不是语言或架构层面的转换问题，因此不在等价迁移
  阶段修正
- **建议**：确认是否存在 `update_user IS NULL` 的存量行。若有，要么补数据、
  要么把条件改为 `update_user IS NULL OR update_user <> ?`。
  同时考虑给 bottle 表加「捞起者」独立列，不再复用审计列承载业务语义

### 4.3 查询方法内做写操作

- **位置**：v1 `UserServiceImpl.getMyFriends`（L453–506）
- **v1 做法**：方法末尾执行 `friendMapper.updateById(friendList)`，即一次读好友列表
  会产生 N 条 `UPDATE friend`（N = 好友数）。原注释写着「这个更新虽然非必要更新，
  但是为了保证数据的一致性，还是更新一下」
- **附带问题**：
  - 走 Redis 缓存命中分支时也会 update，即用缓存里可能过期的数据覆盖 MySQL
  - 每次读都会刷新所有好友行的 `update_time`
  - 方法无 `@Transactional`，N 条 UPDATE 不在事务内
  - `UserServiceImpl:1092` 用 `getMyFriends(userId).size()` 取好友数，也会间接触发
- **v2 处置**：`friend.Repository` 拆成 `ListByOwner`（只读）与 `UpdateAll`（只写）
  两个方法，是否调用由 service 决定，仓储层不再隐含副作用。`UpdateAll` 额外包了
  事务，这一点比 v1 强
- **建议**：读路径不再调 `UpdateAll`。好友资料快照的刷新改为用户资料变更时反向推送

### 4.4 好友表缺索引与唯一约束

- **位置**：`friend`、`friend_request` 两张表
- **现状**：`friend` 除主键外无任何索引，而 `owning_id`、`user_id` 是全部查询的
  过滤条件；`friend_request` 的 `receiver_id`、`status` 同样无索引
- **附带问题**：`friend` 没有 `(owning_id, user_id)` 唯一约束，加上 v1 的
  `addFriend` 不查重，重复同意会插入重复好友行
- **v2 处置**：`FindByUserAndOwner` 用 `First` 取第一条。存在重复记录时，
  取到哪条取决于存储顺序，与 v1 `selectOne` 在重复时直接报错的行为不同
- **建议**：补 `(owning_id, user_id)` 唯一索引与 `friend_request(receiver_id, status)`
  联合索引，并先清理存量重复数据

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

### 5.5 `FriendServiceImpl` 存在注释掉的死代码与不可达方法

- **位置**：v1 `FriendServiceImpl` 第 53–94、110–121 行
- **现状**：两个 `sendFriendRequest` 重载被 `/* */` 整段注释，其中各有一处
  `friendRequestMapper.insert`。第二段的逻辑实际已搬到
  `BottleServiceImpl.sendFriendRequest`，即「漂流瓶发好友申请」这一能力现在
  归属 bottle 模块，不在 friend 模块
- **附带问题**：`canCurrentUserOperateBottle`（L98–107）是活的方法定义，
  但唯一调用点在注释块内，属于不可达代码
- **v2 处置**：均未迁移。`friend.Repository` 的方法集只覆盖 16 处生效调用
- **建议**：v1 侧直接删除

### 5.6 好友申请的自动生成消息硬编码中文

- **位置**：v1 `LetterServiceImpl.readLetter`（L1026）
- **v1 做法**：`content` 直接拼中文字面量 `"!我给你写了一封侨批哦,快来加我为好友吧!"`，
  未走 `MessageUtils.message`，与其余消息的国际化方式不一致
- **v2 处置**：`friend` 域只迁移仓储层，该文案属于服务层，未迁移
- **建议**：纳入 i18n 资源

### 5.7 好友申请来源的判空逻辑语义混乱

- **位置**：v1 `FriendServiceImpl.becomeFriend`（L180–195）
- **v1 做法**：
  ```java
  Bottle bottle = new Bottle();
  Letter letter = new Letter();
  if (friendRequest.getBottleId() != null) { bottle = bottleMapper.selectById(...); ... }
  else if (friendRequest.getLetterId() != null) { letter = letterMapper.selectById(...); ... }
  if (bottle == null || letter == null) { throw ...; }
  ```
  两个变量都被初始化成空对象，只有走进对应分支被 `selectById` 覆盖后才可能为 null，
  而判断用的是 `||`。看起来在做校验，实际只能拦住「走了某分支且查不到」这一种情况
- **附带问题**：两个 ID 都为空时 `mineToFriendAddresses` 保持 null，会一路写进
  `addresses` JSON 列
- **v2 处置**：定义了 `friend.ErrSourceNotFound`，但取哪种语义由 service 决定，
  仓储层不涉及
- **建议**：明确为「来源 ID 必须有且仅有一个，且对应记录必须存在」

### 5.8 `addFriend` 反向记录写错了头像与备注

- **位置**：v1 `FriendServiceImpl.addFriend`（L263–264）
- **v1 做法**：建立双向好友关系时，反向那条记录的 `avatar` 与 `remark` 取的是
  `friend.getAvatar()` / `friend.getNickname()`，即**对方**的头像和昵称，
  而不是同一行其他字段所用的 `myInfoOfUser`
- **v2 处置**：`friend` 域只迁移仓储层，`addFriend` 属于服务层，未迁移
- **建议**：迁移服务层时改为 `myInfoOfUser`，并确认是否需要修正存量数据

### 5.9 `updateFriendRemark` 缺少归属校验

- **位置**：v1 `UserServiceImpl.updateFriendRemark`（L844–852）
- **v1 做法**：用 `selectById(friendId)` 取记录，**未校验 `owning_id` 等于当前用户**。
  同文件的 `getFriendAddress`、`setFriendDefaultAddress`、`deleteFriendAddress`
  都按 `(id, owning_id)` 查询，只有这一处漏了
- **影响**：任何登录用户可以修改任意 friend 行的备注
- **v2 处置**：`Repository` 同时提供 `FindByID` 与 `FindByIDAndOwner`，
  并在 `Friend.OwnedBy` 上提供归属判定，但仓储层不强制
- **建议**：服务层改用 `FindByIDAndOwner`，`FindByID` 只保留给确实不需要归属校验的场景

### 5.10 审计字段自动填充会清空当前用户上下文

- **位置**：v1 `MyMetaObjecthandler`（`Qiaopi-common`）
- **v1 做法**：`insertFill` 与 `updateFill` 的末尾都调用 `UserContext.removeUserId()`，
  即一次写操作之后当前线程的 userId 就没了
- **影响**：写操作之后再取 `UserContext.getUserId()` 会得到 null。
  `setFriendDefaultAddress`、`deleteFriendAddress`、`updateFriendRemark`
  三处都在 `updateById` 之后用它拼缓存键，键实际变成 `...:null`，**缓存删不掉**
- **v2 处置**：Go 侧不存在 ThreadLocal，审计字段由 PO 转换显式填入，
  该问题自然消失。但需注意 v1 存量数据里 `update_user` 可能被写成 0
- **建议**：无需处理，属迁移中自动修正项。若要核对存量数据，
  注意 `update_user = 0` 的行不代表系统操作

### 5.11 信件版式常量命名与取值不符

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

### 6.4 `friend` 表的字段名与列名怪癖

- **位置**：v1 `Friend.OwningId`
- **现状**：Java 字段名首字母大写（`OwningId` 而非 `owningId`），
  靠 Lombok 生成 `getOwningId()`、再由 MyBatis-Plus 驼峰转换推导出列名 `owning_id`。
  能跑通但违反 Java 命名惯例
- **v2 处置**：领域模型改名 `OwnerID`，PO 字段名 `OwningID` 并显式声明
  `gorm:"column:owning_id"`，不依赖命名推导
- **建议**：无需额外处理

### 6.5 `friend.remark` 被复用为好友备注名

- **位置**：`friend.remark`
- **现状**：`remark` 是审计基类的通用备注列，在好友场景被复用为「好友备注名」，
  建立关系时默认填对方昵称
- **v2 处置**：`friend.Friend.Remark` 保留该语义并在注释中说明
- **建议**：若后续需要同时保留「备注名」与「通用备注」，需拆列

### 6.6 `questions` 表四个选项列

- **位置**：`questions.option_a` ~ `option_d`
- **现状**：同一件东西的四份拷贝，任何遍历都要写死四次
- **v2 处置**：`po_explore.go` 折叠为 `map[explore.Option]string`
- **建议**：若要支持选项数量可变，需改为子表；当前四选一场景下折叠已足够
