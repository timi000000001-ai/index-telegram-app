# 机器人信息表（bot_info）

用于记录每个机器人实例的基础信息、运行周期与统计数据。该设计与 PocketBase 迁移中的 `bot_info` 集合字段保持一致，可直接映射。

## 表结构

| 字段名               | 类型         | 必填 | 默认值            | 说明 |
|----------------------|--------------|------|-------------------|------|
| id                   | bigint       | 是   |                   | 主键ID（全局唯一标识）。 |
| user                 | bigint       | 是   |                   | 所属用户ID（外键，关联 users.id）。在迁移中为 Relation 字段。 |
| bot_name             | varchar(255) | 是   |                   | 机器人名称。 |
| bot_token            | varchar(255) | 是   |                   | 机器人访问令牌（仅密文安全存储，不应明文展示/日志）。 |
| bot_status           | varchar(32)  | 否   | offline           | 机器人状态（如 offline/online/paused 等，可扩展）。 |
| start_time           | datetime     | 否   |                   | 启动时间。 |
| end_time             | datetime     | 否   |                   | 停止时间。 |
| running_days         | int          | 否   | 0                 | 运行天数（冗余统计）。 |
| total_messages       | bigint       | 否   | 0                 | 累计消息数。 |
| total_users          | bigint       | 否   | 0                 | 覆盖用户数。 |
| total_groups         | bigint       | 否   | 0                 | 覆盖群组数。 |
| avg_messages_per_day | decimal(10,2)| 否   | 0                 | 日均消息数。 |
| create_time          | datetime     | 是   | CURRENT_TIMESTAMP | 创建时间。 |
| update_time          | datetime     | 是   | CURRENT_TIMESTAMP | 更新时间。 |

说明：迁移脚本中该集合字段包括 user(关联)、bot_name、bot_token、bot_status、start_time、end_time、running_days、total_messages、total_users、total_groups、avg_messages_per_day。

## 字段说明
- id：主键。
- user：所属用户，外键关联用户表。
- bot_name：机器人名称，业务侧可唯一约束（可选）。
- bot_token：敏感凭证，需加密或安全存储（KMS/密钥管理），严禁打印日志。
- bot_status：运行状态枚举，可与运行器/调度器联动。
- start_time/end_time：运行区间；若在线运行可仅记录 start_time。
- running_days：辅助统计字段，也可由 start/end 计算得出。
- total_messages/total_users/total_groups：运行期累计统计。
- avg_messages_per_day：日均消息数，可由总量/天数计算得出。
- create_time/update_time：审计字段。

## 索引与约束
- PK：PRIMARY KEY(id)
- FK_user：FOREIGN KEY(user) REFERENCES users(id)
- 常用查询索引：INDEX(user, bot_status)、INDEX(start_time)、INDEX(update_time)

## 安全建议
- bot_token：
  - 存储：加密或使用安全后端（如 Vault/KMS），数据库中可存密文或引用。
  - 使用：在调用外部 API 时解密，避免在应用日志、监控中泄露。
  - 访问：仅限拥有最小权限的服务访问，记录访问审计。