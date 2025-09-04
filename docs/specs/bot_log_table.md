# 机器人操作日志表（bot_log / operation_logs）

用于记录机器人相关的关键操作与审计信息，便于问题追踪与安全合规。与迁移脚本中的 `operation_logs` 集合字段映射如下。

## 表结构

| 字段名            | 类型         | 必填 | 默认值            | 说明 |
|-------------------|--------------|------|-------------------|------|
| id                | bigint       | 是   |                   | 主键ID。 |
| user              | bigint       | 是   |                   | 操作用户ID（外键，关联 users.id）。迁移中为 Relation 字段。 |
| bot_id            | varchar(255) | 否   |                   | 机器人ID/标识（字符串存储，视运行环境而定）。 |
| operation_type    | varchar(64)  | 是   |                   | 操作类型（如 create/update/delete/start/stop/pause 等）。 |
| operation_time    | datetime     | 是   | CURRENT_TIMESTAMP | 操作时间。 |
| operation_details | text         | 否   |                   | 操作详情（结构化/半结构化文本；可考虑 JSON）。 |
| create_time       | datetime     | 是   | CURRENT_TIMESTAMP | 创建时间（与 operation_time 一致或作为审计创建时间）。 |

说明：迁移脚本原始字段包括 user、bot_id、operation_type、operation_time、operation_details；此处增加 `create_time` 以对齐三表的审计字段一致性（如在关系型数据库方案中实现）。

## 字段说明
- id：主键。
- user：执行该操作的用户ID。
- bot_id：被操作的机器人标识（可存外部 ID）。
- operation_type：操作类别枚举，建议配套白名单校验并在业务层定义常量。
- operation_time：操作发生的时间。
- operation_details：详细内容（如 JSON 字符串，记录变更前后差异、上下文信息、IP 等）。
- create_time：记录创建时间，便于统一检索排序。

## 索引与约束
- PK：PRIMARY KEY(id)
- FK_user：FOREIGN KEY(user) REFERENCES users(id)
- 常用查询索引：INDEX(user, operation_type, operation_time)

## 审计与合规建议
- 日志不可随意删除；需要保留策略（例如 180 天或更长）。
- 内容敏感信息打码或脱敏存储（如 token、密钥）。
- 结合请求上下文记录 request_id、trace_id，便于分布式追踪。