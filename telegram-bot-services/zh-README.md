# Telegram 数据采集和机器人服务

欢迎使用 Telegram 数据采集和机器人服务项目！这是一个基于 Go 开发的模块化应用，旨在与 Telegram 交互，采集群组/频道消息，并提供可搜索的机器人界面。该项目包含三个主要服务：机器人服务、管理服务和采集服务，与 PocketBase、Meilisearch 以及 Svelte 前端集成。

**最后更新：** 2025年9月3日星期三 13:38 WIB  
**许可：** MIT（详情见 LICENSE 文件）

## 目录

- [概述](#概述)
- [功能](#功能)
- [架构](#架构)
- [先决条件](#先决条件)
- [安装](#安装)
- [配置](#配置)
- [使用](#使用)
- [开发](#开发)
- [API 端点](#api-端点)
- [法律条款与风险规避](#法律条款与风险规避)
- [贡献](#贡献)
- [许可](#许可)

## 概述

该项目支持以下功能：

- 提供 Telegram 机器人，支持分页和过滤查询（机器人服务）。
- 管理服务提供 REST API 用于搜索和会话管理。
- 采集服务负责登录 Telegram 账户，配置目标群组，并采集消息以供索引。
- 集成 Svelte 前端，提供用户友好的界面。

这些服务协同工作，利用 PocketBase 存储数据，Meilisearch 提供快速搜索功能。

## 功能

### 机器人服务

- 支持多个 Telegram 机器人，集成 Webhook。
- 提供分页搜索结果（每页 5 条），包含"上一页"和"下一页"按钮。
- 支持按群组、频道、机器人或所有消息过滤搜索结果。
- 命令：/help、/clong（克隆机器人）、/sponsor、/mini。
- 将消息存储到 PocketBase 并索引到 Meilisearch。

[机器人服务详细文档](./bot-service/README.md)

### 管理服务

- 提供 REST API 支持带分页和过滤的搜索。
- 作为前端与后端服务的桥梁。

[管理服务详细文档](./management-service/README.md)

### 采集服务

- 处理 Telegram 账户登录，使用电话号码认证。
- 允许配置目标群组以进行消息采集。
- 主动采集并解析群组消息，存储以供搜索。

[采集服务详细文档](./collection-service/README.md)

## 架构

```
[Svelte 前端 (静态 HTML)] ↔ [管理服务 (:8080)]
                                    ↕ (REST API)
[机器人服务 (:8081)] ↔ [采集服务 (:8082)]
      ↕                           ↕
[PocketBase (存储)]       [Meilisearch (搜索)]
      ↕                           ↕
[S3/MinIO (媒体)]
```

- **机器人服务：** 运行在 :8081，处理机器人交互和消息存储。
- **管理服务：** 运行在 :8080，提供搜索 API。
- **采集服务：** 运行在 :8082，管理 Telegram 数据采集。

## 先决条件

- **Go：** 版本 1.21 或更高。
- **Docker：** 用于容器化部署（可选）。
- **Telegram API 凭据：** 从 my.telegram.org 获取 api_id 和 api_hash。
- **PocketBase：** 运行实例，包含 messages 集合。
- **Meilisearch：** 运行实例，包含 messages 索引。
- **SvelteKit：** 用于构建静态前端（可选）。

## 安装

### 克隆仓库

```bash
git clone https://github.com/your-repo/telegram-bot-services.git
cd telegram-bot-services
```

### 初始化模块

导航到每个服务目录，初始化 Go 模块：

```bash
cd bot-service && go mod init bot-service && go mod tidy
cd ../management-service && go mod init management-service && go mod tidy
cd ../collection-service && go mod init collection-service && go mod tidy
cd ..
```

### 构建服务

构建每个服务：

```bash
go build -o bot-service ./bot-service
go build -o management-service ./management-service
go build -o collection-service ./collection-service
```

### Docker 安装（可选）

构建 Docker 镜像：

```bash
docker build -t bot-service ./bot-service
docker build -t management-service ./management-service
docker build -t collection-service ./collection-service
```

## 配置

### 环境变量

在每个服务目录中创建 .env 文件或设置环境变量：

#### 机器人服务

```
BOT_TOKENS=YOUR_BOT_TOKEN_1,YOUR_BOT_TOKEN_2
POCKETBASE_TOKEN=YOUR_POCKETBASE_TOKEN
MEILISEARCH_KEY=YOUR_MEILISEARCH_KEY
AUTH_TOKEN=YOUR_AUTH_TOKEN
```

#### 管理服务

```
MEILISEARCH_KEY=YOUR_MEILISEARCH_KEY
```

#### 采集服务

```
API_ID=123456
API_HASH=your_api_hash
POCKETBASE_URL=http://your-pocketbase-url
POCKETBASE_TOKEN=YOUR_POCKETBASE_TOKEN
MEILISEARCH_URL=http://your-meilisearch-url
MEILISEARCH_KEY=YOUR_MEILISEARCH_KEY
```

### SSL 证书（机器人服务）

为 Webhook 支持，提供 cert.pem 和 key.pem（例如通过 Let's Encrypt）：

```bash
certbot certonly --standalone -d your-bot-service.com
```

## 使用

### 运行服务

#### 机器人服务

```bash
./bot-service
```

或使用 Docker：

```bash
docker run -p 8081:8081 -v /path/to/cert.pem:/app/cert.pem -v /path/to/key.pem:/app/key.pem bot-service
```

#### 管理服务

```bash
./management-service
```

或使用 Docker：

```bash
docker run -p 8080:8080 management-service
```

#### 采集服务

```bash
./collection-service
```

或使用 Docker：

```bash
docker run -p 8082:8082 collection-service
```

### 与机器人交互

- 使用命令如 /help、/clong、/sponsor、/mini。
- 输入查询（≤10 个字符）以查看分页结果并应用过滤。

### 配置采集

#### 登录

```bash
curl -X POST http://localhost:8082/login -H "Content-Type: application/json" -d '{"phone_number": "+1234567890"}'
```

输入 Telegram 应用发送的验证码。

#### 配置群组

```bash
curl -X POST http://localhost:8082/configure -H "Content-Type: application/json" -d '{"chat_ids": [123, 456]}'
```

服务将开始从指定 chat_ids 采集消息。

### 前端集成

构建 Svelte 前端并提供静态 HTML：

```bash
cd frontend
npm install
npm run build
```

通过 web 服务器（例如 Nginx）提供 build/index.html。

## 开发

### 依赖

- 使用 go mod tidy 更新 go.mod 文件。
- 根据需要添加新依赖。

### 测试

为每个服务运行测试：

```bash
go test ./...
```

### 添加功能

- **机器人服务：** 扩展命令或添加新过滤器。
- **采集服务：** 使用 gotd Updates API 实现实时更新。
- **管理服务：** 添加更多 API 端点。

## API 端点

### 机器人服务

- **POST /webhook/{token}：** 处理 Telegram Webhook 更新。

### 管理服务

- **GET /api/search：** 搜索，参数包括 q、page、limit、filter。

### 采集服务

- **POST /login：** 登录，数据格式 {"phone_number": "..."}。
- **POST /configure：** 配置，数据格式 {"chat_ids": [...]}。
- **GET /health：** 健康检查。

## 法律条款与风险规避

- **使用范围限制：** 本项目不得在中国大陆地区使用。Telegram 在中国大陆被政府限制访问，且相关数据采集和处理可能违反当地法律法规。使用本项目需确保遵守您所在国家或地区的法律和规定。
- **免责声明：** 开发者不对因违反当地法律、使用不当或数据隐私问题导致的任何后果负责。用户应自行评估法律风险，并在必要时咨询法律专业人士。
- **建议：** 如果您位于中国大陆，请勿下载、安装或运行本项目。如需类似功能，请寻找符合当地法规的替代方案。

## 贡献

1. Fork 该仓库。
2. 创建功能分支（git checkout -b feature/awesome-feature）。
3. 提交更改（git commit -m "添加 awesome 功能"）。
4. 推送分支（git push origin feature/awesome-feature）。
5. 打开 Pull Request。

## 许可

该项目采用 MIT 许可。详情见 LICENSE 文件。

/*
 * 文件功能描述：Telegram 数据采集和机器人服务主 README
 * 主要类/接口说明：机器人服务、管理服务和采集服务概述
 * 修改历史记录：
 * @author fcj
 * @date 2025-09-03
 * @version 1.1.0
 * © Telegram Bot Services Team
 */

