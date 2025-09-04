# Telegram 采集服务

## 功能概述

**Telegram 采集服务**是 Telegram 数据采集和机器人服务项目的核心组件，负责登录 Telegram 账户，配置目标群组，并采集消息以供索引和搜索。

### 主要功能

- **Telegram 账户登录**：处理电话号码认证流程
- **群组配置**：允许指定目标群组进行消息采集
- **消息采集**：自动采集和解析群组消息
- **数据存储与索引**：将采集的消息存储到 PocketBase 并索引到 Meilisearch

## 技术架构

采集服务运行在端口 `:8082`，与其他服务协同工作：

```
[采集服务 (:8082)] ↔ [机器人服务 (:8081)]
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
    "port": "8082"
  },
  "telegram": {
    "apiID": 2345678,
    "apiHash": "your_telegram_api_hash"
  },
  "storage": {
    "pocketBaseURL": "http://127.0.0.1:8090"
  },
  "search": {
    "meilisearchURL": "http://127.0.0.1:7700",
    "meilisearchToken": "your_meilisearch_token"
  }
}
```

## 安装与运行

### 直接运行

```bash
# 初始化模块
go mod init collection-service && go mod tidy

# 构建服务
go build -o collection-service .

# 运行服务
./collection-service
```

### Docker 运行

```bash
# 构建 Docker 镜像
docker build -t collection-service .

# 运行容器
docker run -p 8082:8082 collection-service
```

## API 端点

- **POST /login**：登录 Telegram 账户
  ```json
  {"phone_number": "+1234567890"}
  ```

- **POST /configure**：配置目标群组
  ```json
  {"chat_ids": [123, 456]}
  ```

- **GET /health**：健康检查

## 使用流程

### 1. 登录 Telegram 账户

```bash
curl -X POST http://localhost:8082/login -H "Content-Type: application/json" -d '{"phone_number": "+1234567890"}'
```

系统会发送验证码到您的 Telegram 应用。收到验证码后，提交验证码完成登录：

```bash
curl -X POST http://localhost:8082/login -H "Content-Type: application/json" -d '{"phone_number": "+1234567890", "code": "12345"}'
```

### 2. 配置目标群组

```bash
curl -X POST http://localhost:8082/configure -H "Content-Type: application/json" -d '{"chat_ids": [123, 456]}'
```

配置完成后，服务将自动开始采集指定群组的消息。

## 开发指南

### 添加新功能

- 实现实时更新（使用 gotd Updates API）
- 增强消息解析能力
- 优化数据存储和索引效率

### 测试

```bash
go test ./...
```

## 注意事项

- 确保提供有效的 Telegram API 凭据
- 验证 PocketBase 和 Meilisearch 服务正常运行
- 遵守 Telegram API 使用条款和限制
- 注意数据隐私和安全问题

/*
 * 文件功能描述：采集服务 README
 * 主要类/接口说明：提供采集服务的安装、配置和使用说明
 * 修改历史记录：
 * @author fcj
 * @date 2025-09-03
 * @version 1.0.0
 * © Telegram Bot Services Team
 */