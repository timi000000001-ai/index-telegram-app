# 用户表（tele_user）

本表用于存储 Telegram 官方“User/teleUser”信息的快照与扩展字段，仅面向数据建模与业务使用，不涉及认证与授权逻辑。数据来源：Telegram Bot API（User 对象）。

## 表结构（参考 Telegram Bot API 的 User 对象）

| 字段名                     | 类型            | 必填 | 默认值              | 说明 |
|----------------------------|-----------------|------|---------------------|------|
| id                         | bigint          | 是   |                     | 本系统主键ID（自增/雪花/UUID 映射为 bigint 可选，保证全局唯一）。 |
| tg_user_id                 | bigint          | 是   |                     | Telegram 用户ID（全局唯一，建议唯一索引）。 |
| is_bot                     | tinyint(1)      | 是   | 0                   | 是否为机器人账号（true/false）。 |
| first_name                 | varchar(64)     | 是   |                     | 用户的名字（必有）。 |
| last_name                  | varchar(64)     | 否   |                     | 用户的姓氏（可选）。 |
| username                   | varchar(32)     | 否   |                     | Telegram 用户名（可选，可能不存在）。 |
| language_code              | varchar(10)     | 否   |                     | IETF 语言代码（如 "en"、"zh-CN"，可选）。 |
| is_premium                 | tinyint(1)      | 否   | 0                   | 是否为 Telegram Premium 用户（可选）。 |
| can_join_groups            | tinyint(1)      | 否   | NULL                | 能否加入群组（Bot 能力字段，普通用户可留空）。 |
| can_read_all_group_messages| tinyint(1)      | 否   | NULL                | 能否读取所有群消息（Bot 能力字段，普通用户可留空）。 |
| supports_inline_queries    | tinyint(1)      | 否   | NULL                | 是否支持内联查询（Bot 能力字段，普通用户可留空）。 |
| added_to_attachment_menu   | tinyint(1)      | 否   | NULL                | 是否已添加到附件菜单（Bot 能力字段，普通用户可留空）。 |
| bio                        | text            | 否   |                     | 个人简介（可选）。 |
| photo_url                  | text            | 否   |                     | 头像 URL（如有媒体代理/CDN，可存镜像地址）。 |
| raw_json                   | json            | 否   | NULL                | 原始 Telegram User 对象的 JSON 快照（来自 Bot API），用于无损还原与比对，不用于查询。 |
| create_time                | datetime        | 是   | CURRENT_TIMESTAMP   | 创建时间。 |
| update_time                | datetime        | 是   | CURRENT_TIMESTAMP   | 更新时间（更新时自动刷新）。 |

说明：
- 字段对齐 Telegram Bot API 的 User 常见属性；部分“Bot 能力”字段对普通用户可为空，用于统一建模与兼容机器人账号。 
- 不包含登录口令、认证状态等与权限相关的字段（由其他系统/集合承担）。
- raw_json 存储完整的原始对象，有助于新版本字段兼容与审计，但请注意隐私与合规要求。

## 字段说明
- id：系统内部主键，用于高效关联与检索。
- tg_user_id：Telegram 平台侧用户的唯一标识，强唯一索引以避免重复。
- is_bot：标识该记录是否对应机器人账号（true/false）。
- first_name/last_name：来自 Telegram 的用户姓名信息。
- username：Telegram 用户名（可能不存在）；若存在，建议建立普通索引以支持 @username 查询。
- language_code：IETF 语言代码（如 en、zh-CN）。
- is_premium：是否为 Telegram Premium 用户（如同步得到）。
- can_join_groups/can_read_all_group_messages/supports_inline_queries/added_to_attachment_menu：Bot 能力相关字段；普通用户可为 NULL，机器人账号可存储真实能力。
- bio：用户个人简介（如通过扩展接口获取）。
- photo_url：头像 URL（可自行代理/持久化）。
- raw_json：保存从 Bot API 获得的 User 原始 JSON；用于兼容未来字段变化、故障排查与审计对比，不建议参与业务查询。
- create_time/update_time：审计字段，记录创建与最近更新时间。

## 索引与约束
- PK：PRIMARY KEY(id)
- UK_tg_user_id：UNIQUE(tg_user_id)
- 常用查询索引：INDEX(username)、INDEX(is_bot)、INDEX(update_time)

## 同步与一致性建议（Bot API）
- 以 Bot API 返回为准进行增量/全量同步，更新本地快照与 raw_json。
- 允许 username、last_name 等可选字段为空，避免强制约束导致同步失败。
- 如需头像持久化，建议引入媒体存储与定期刷新策略（防止外链失效）。