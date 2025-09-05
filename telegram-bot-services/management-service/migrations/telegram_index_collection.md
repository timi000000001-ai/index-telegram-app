# Telegram 索引表结构设计

## 表名：telegram_index

### 基础信息字段
| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | string | 主键 系统主键 | PRIMARY KEY |
| chat_id | string | 主键，使用 chatid/teleid | PRIMARY KEY |
| type | string | 类型：group/channel/bot | NOT NULL |
| title | string | 名称/标题 | NOT NULL |
| username | string | 用户名（@username） |  |
| description | text | 描述信息 |  |
| first_name | string | 名字（个人聊天） |  |
| last_name | string | 姓氏（个人聊天） |  |
| is_verified | boolean | 是否认证 | DEFAULT false |
| is_restricted | boolean | 是否受限 | DEFAULT false |
| is_scam | boolean | 是否诈骗 | DEFAULT false |
| is_fake | boolean | 是否虚假 | DEFAULT false |
| language_code | string | 语言代码 |  |

### 统计信息字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| members_count | integer | 成员数量 |
| active_members | integer | 活跃成员数 |
| online_count | integer | 在线成员数 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 最后更新时间 |
| last_activity | datetime | 最后活跃时间 |
| message_count | integer | 消息总数 |
| avg_message_per_day | float | 日均消息数 |
| growth_rate | float | 增长率 |

### 内容特征字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| tags | []string | 主题标签数组 |
| keywords | []string | 关键词数组 |
| content_types | []string | 内容类型数组 |
| topics | []string | 讨论主题数组 |
| categories | []string | 分类标签数组 |
| language | string | 主要语言 |
| country | string | 主要国家/地区 |
| content_quality | integer | 内容质量评分(1-5) |

### 管理信息字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| admin_list | []string | 管理员ID列表 |
| creator_id | string | 创建者ID |
| permissions | JSON | 权限设置 |
| invite_link | string | 邀请链接 |
| linked_chat_id | string | 关联聊天ID |
| slow_mode_delay | integer | 慢模式延迟秒数 |
| has_protected_content | boolean | 是否有保护内容 |
| can_set_sticker_set | boolean | 是否可以设置贴纸包 |

### 评分字段（多维度评分体系）
| 字段名 | 类型 | 说明 | 权重 |
|--------|------|------|------|
| overall_score | float | 综合评分 | 1.0 |
| popularity_score | float | 流行度评分 | 0.3 |
| activity_score | float | 活跃度评分 | 0.25 |
| content_score | float | 内容质量评分 | 0.2 |
| relevance_score | float | 相关性评分 | 0.15 |
| authority_score | float | 权威性评分 | 0.1 |
| last_calculated | datetime | 最后评分计算时间 |  |

### Bot 特定字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| can_join_groups | boolean | 是否可以加入群组 |
| can_read_all_group_messages | boolean | 是否可以读取所有群消息 |
| supports_inline_queries | boolean | 是否支持内联查询 |
| bot_type | string | 机器人类型 |
| commands | JSON | 支持的命令列表 |

### 扩展字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| metadata | JSON | 元数据扩展字段 |
| custom_fields | JSON | 自定义字段 |
| indexed_at | datetime | 索引时间 |
| data_source | string | 数据来源 |
| version | integer | 数据结构版本 |

## 索引设计

### 主索引
- PRIMARY KEY (id)

### 查询索引
- INDEX idx_type (type)
- INDEX idx_score (overall_score DESC)
- INDEX idx_members (members_count DESC) 
- INDEX idx_activity (last_activity DESC)
- INDEX idx_created (created_at DESC)
- INDEX idx_username (username)
- INDEX idx_tags (tags)
- INDEX idx_keywords (keywords)

### 复合索引
- INDEX idx_type_score (type, overall_score DESC)
- INDEX idx_type_activity (type, last_activity DESC)
- INDEX idx_type_members (type, members_count DESC)

## 示例记录结构

```json
{
  "id": "-100123456789",
  "type": "group",
  "title": "技术交流群",
  "username": "techgroup",
  "description": "专注于编程和技术讨论的Telegram群组",
  "members_count": 1500,
  "active_members": 300,
  "online_count": 50,
  "created_at": "2023-01-15T10:00:00Z",
  "last_activity": "2024-01-20T15:30:00Z",
  "message_count": 50000,
  "tags": ["programming", "technology", "development"],
  "keywords": ["python", "javascript", "ai", "web开发"],
  "content_types": ["discussion", "news", "qna"],
  "admin_list": ["123456789", "987654321"],
  "permissions": {
    "can_send_messages": true,
    "can_send_media": true,
    "can_send_polls": true,
    "can_change_info": false
  },
  "overall_score": 8.5,
  "popularity_score": 9.0,
  "activity_score": 7.8,
  "content_score": 8.2,
  "relevance_score": 8.7,
  "authority_score": 7.5,
  "indexed_at": "2024-01-20T16:00:00Z"
}
```

## 查询优化建议

1. **排序查询**：使用复合索引 (type + score) 进行分页排序
2. **标签搜索**：使用数组字段索引进行标签过滤
3. **全文搜索**：对 title、description、keywords 建立全文索引
4. **范围查询**：对数值型字段（members_count、scores）使用范围索引
5. **时间查询**：对时间字段使用时间范围索引

## 扩展性考虑

1. **JSON字段**：metadata 和 custom_fields 用于存储非结构化数据
2. **版本控制**：version 字段支持数据结构升级
3. **插件系统**：预留扩展字段支持未来功能添加
4. **多语言支持**：language 和 language_code 字段支持国际化

这个表结构设计支持高效的查询、排序和扩展，能够满足Telegram信息索引的所有需求。