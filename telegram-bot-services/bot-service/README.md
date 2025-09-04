# Telegram 机器人服务

## 功能概述

**Telegram 机器人服务**是 Telegram 数据采集和机器人服务项目的核心组件之一，负责处理用户与 Telegram 机器人的交互，提供搜索功能，并将消息存储到数据库。

### 主要功能

- **多机器人支持**：可同时管理多个 Telegram 机器人，通过 Webhook 接收更新
- **分页搜索**：每页显示 5 条结果，带有"上一页"和"下一页"导航按钮
- **结果过滤**：支持按群组、频道、机器人或所有消息进行过滤
- **命令支持**：内置 `/help`、`/clong`（克隆机器人）、`/sponsor`、`/mini` 等命令
- **数据存储**：将消息存储到 PocketBase 并索引到 Meilisearch 以实现快速搜索

## 技术架构

机器人服务运行在端口 `:8081`，与其他服务协同工作：

```
[机器人服务 (:8081)] ↔ [采集服务 (:8082)]
         ↕
[PocketBase (存储)] ↔ [Meilisearch (搜索)]
```

## 配置说明

服务通过 `configs` 目录下的 JSON 文件进行配置，支持多种运行环境。通过设置 `APP_ENV` 环境变量来选择加载对应的配置文件。

- **development.json**: 开发环境配置 (默认)
- **production.json**: 生产环境配置
- **testing.json**: 测试环境配置

### 环境变量

- `APP_ENV`: 指定运行环境，可选值为 `development`、`production`、`testing`。

例如，要以生产模式运行，请设置：

```bash
export APP_ENV=production
```

### 配置文件示例 (`configs/development.json`)

```json
{
  "server": {
    "port": "8081"
  },
  "storage": {
    "pocketBaseURL": "http://127.0.0.1:8090/",
    "meilisearchURL": "http://127.0.0.1:7700",
    "meilisearchToken": "timigogogo"
  },
  "search": {
    "meilisearchURL": "http://127.0.0.1:7700",
    "meilisearchKey": "timigogogo",
    "managementServiceURL": "http://127.0.0.1:8080"
  }
}
```

### SSL 证书

为 Webhook 支持，需要提供 `cert.pem` 和 `key.pem` 文件（例如通过 Let's Encrypt 获取）：

```bash
certbot certonly --standalone -d your-bot-service.com
```

## 安装与运行

### 直接运行

```bash
# 初始化模块
go mod init bot-service && go mod tidy

# 构建服务
go build -o bot-service .

# 运行服务
./bot-service
```

### Docker 运行

```bash
# 构建 Docker 镜像
docker build -t bot-service .

# 运行容器
docker run -p 8081:8081 -v /path/to/cert.pem:/app/cert.pem -v /path/to/key.pem:/app/key.pem bot-service
```

## API 端点

- **POST /webhook/{token}**：处理 Telegram Webhook 更新

## 开发指南

### 添加新功能

- 扩展机器人命令
- 添加新的搜索过滤器
- 优化搜索结果展示

### 测试

```bash
go test ./...
```

## 注意事项

- 确保 PocketBase 和 Meilisearch 服务正常运行
- 验证 Telegram 机器人令牌的有效性
- 检查 SSL 证书的有效性和路径

/*
 * 文件功能描述：机器人服务 README
 * 主要类/接口说明：提供机器人服务的安装、配置和使用说明
 * 修改历史记录：
 * @author fcj
 * @date 2025-09-03
 * @version 1.0.0
 * © Telegram Bot Services Team
 */