# 迁移文件修复说明 - m_1721986800_create_initial_collections.go

## 问题描述

原始迁移文件使用了 PocketBase v0.22.x 及更早版本的 API，在 v0.23.0+ 版本中以下包和API已废弃：
- `github.com/pocketbase/pocketbase/models`
- `github.com/pocketbase/pocketbase/models/schema`
- `app.Dao()` 方法

## 修复内容

### 1. 导入包调整
**移除：**
```go
"github.com/pocketbase/pocketbase/models"
"github.com/pocketbase/pocketbase/models/schema"
```

### 2. API 方法更新

| 旧方法 | 新方法 | 说明 |
|--------|--------|------|
| `app.Dao().FindCollectionByNameOrId()` | `app.FindCollectionByNameOrId()` | 直接使用 app 实例 |
| `app.Dao().SaveCollection()` | `app.Save()` | 直接使用 app 实例 |
| `&models.Collection{}` | `core.NewBaseCollection()` | 使用核心包工厂方法 |

### 3. 字段定义方式变更

**旧方式：**
```go
&schema.SchemaField{
    Name: "field_name",
    Type: schema.FieldTypeText,
    Options: &schema.TextOptions{Max: types.Pointer(100)},
}
```

**新方式：**
```go
&core.TextField{
    Name: "field_name",
    Max:  100,
}
```

### 4. 支持的字段类型

- `core.RelationField` - 关联字段
- `core.TextField` - 文本字段
- `core.URLField` - URL字段
- `core.DateField` - 日期字段
- `core.NumberField` - 数字字段

## 迁移文件结构

迁移文件现在遵循以下模式：

```go
package migrations

import (
    "github.com/pocketbase/pocketbase/core"
    m "github.com/pocketbase/pocketbase/migrations"
    "github.com/pocketbase/pocketbase/tools/types"
)

func init() {
    m.Register(func(app core.App) error {
        // 升级逻辑
        return nil
    }, func(app core.App) error {
        // 降级逻辑（可选）
        return nil
    })
}
```

## 验证

迁移文件修复后已通过以下验证：
1. `go mod tidy` - 依赖检查通过
2. `go build` - 编译通过
3. 服务启动 - 数据库迁移自动应用

## 注意事项

1. **PocketBase 版本要求：** v0.23.0+
2. **向后兼容性：** 新API不兼容旧版本
3. **字段选项：** 部分字段选项的配置方式有所变化
4. **错误处理：** 保持原有的错误处理逻辑

## 相关资源

- [PocketBase v0.23.0 升级指南](https://pocketbase.io/v023upgrade/go/)
- [PocketBase Go 迁移文档](https://pocketbase.io/docs/go-migrations/)
- [集合操作文档](https://pocketbase.io/docs/go-collections/)

---

*最后更新: 2024年*  
*修复版本: PocketBase v0.29.3*