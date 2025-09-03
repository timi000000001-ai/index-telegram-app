# Telegram 搜索平台

本项目是一个综合性的 Telegram 搜索平台，由用于数据收集的后端服务和用于用户交互的 Web 应用程序组成。

## 搜索机器人 https://t.me/TeleSearchVipbot
## 搜索大群 https://t.me/SoSo00000000001
## 搜索更新频道 https://t.me/SoSo00000000002
## 搜索消息监听 https://t.me/SoSo00000000003
## Telegram 联系方式：https://t.me/simi001001 

## 项目结构

该项目分为两个主要部分：

- `telegram-bot-services/`：一个用 Go 编写的后端服务集合。这些服务负责与 Telegram API 交互，从群组和频道收集消息，并提供搜索 API。
- `index-telegram-app/`：一个基于 Svelte 的 Web 应用程序，为搜索后端服务收集的内容提供了一个用户友好的界面。

## 功能

- **数据收集**：后端服务可以登录 Telegram 帐户，加入指定的群组和频道，并收集消息。
- **搜索 API**：提供强大的搜索 API 来查询收集到的消息，并支持过滤和分页。
- **Web 界面**：现代化且响应迅速的 Web 界面，允许用户轻松搜索内容、查看结果和管理设置。
- **机器人集成**：该平台包括一个 Telegram 机器人，可用于直接在 Telegram 中搜索内容。

## 架构

后端由三个主要服务组成：

- **机器人服务**：处理用户与 Telegram 机器人的交互，提供搜索功能并存储消息。
- **管理服务**：公开用于搜索和会话管理的 REST API，充当前端和后端服务之间的桥梁。
- **收集服务**：管理 Telegram 帐户登录，配置目标组并收集用于索引的消息。

这些服务旨在协同工作，利用 PocketBase 进行存储，利用 Meilisearch 实现快速搜索功能。

```
[Svelte 前端 (静态 HTML)] ↔ [管理服务 (:8080)]
                                    ↕ (REST API)
[机器人服务 (:8081)] ↔ [收集服务 (:8082)]
      ↕                           ↕
[PocketBase (存储)]       [Meilisearch (搜索)]
      ↕                           ↕
[S3/MinIO (媒体)]
```

## 入门指南

### 先决条件

- Go 1.21+
- Docker (可选)
- Telegram API 凭证
- PocketBase
- Meilisearch
- SvelteKit (可选)

### 安装

1.  **克隆存储库：**
    ```bash
    git clone https://github.com/your-repo/telegram-search-platform.git
    cd telegram-search-platform
    ```
2.  **初始化 Go 模块：**
    ```bash
    cd telegram-bot-services/bot-service && go mod init bot-service && go mod tidy
    cd ../management-service && go mod init management-service && go mod tidy
    cd ../collection-service && go mod init collection-service && go mod tidy
    cd ../..
    ```
3.  **构建服务：**
    ```bash
    go build -o telegram-bot-services/bot-service/bot-service telegram-bot-services/bot-service
    go build -o telegram-bot-services/management-service/management-service telegram-bot-services/management-service
    go build -o telegram-bot-services/collection-service/collection-service telegram-bot-services/collection-service
    ```
4.  **安装前端依赖项：**
    ```bash
    cd index-telegram-app
    npm install
    ```

### 配置

在每个服务目录（`bot-service`、`management-service`、`collection-service`）中创建一个 `.env` 文件，并提供必要的环境变量。有关更多详细信息，请参阅 `telegram-bot-services` 目录中的 `en-README.md`。

## 用法

### 运行服务

-   **机器人服务**：`./telegram-bot-services/bot-service/bot-service`
-   **管理服务**：`./telegram-bot-services/management-service/management-service`
-   **收集服务**：`./telegram-bot-services/collection-service/collection-service`

### 运行前端

```bash
cd index-telegram-app
npm run dev
```

## API 文档

API 文档以 Postman 集合的形式在 `telegram-bot-services/api-docs` 目录中提供。

## 法律声明

**使用限制**：本项目不适用于中国大陆。Telegram 在中国大陆受到政府的访问限制，本项目的数据收集和处理活动可能违反当地法律法规。

**免责声明**：本项目开发人员对因使用不当、违反当地法律或数据隐私问题而导致的任何后果概不负责。用户应自行评估法律风险，并在必要时咨询法律专业人士。

**建议**：如果您位于中国大陆，请不要下载、安装或运行本项目。请寻找符合当地法规的替代方案。

## 赞助

如果您觉得这个项目对您有帮助，请考虑赞助我们。

- **TRX & USDT (TRC20):** `TD5JGaR7cY5ZxDnZNgmCSv66axR9DhrcYz`

<img width="512" height="530" alt="image" src="https://github.com/user-attachments/assets/08a5cf87-e174-4bf5-ae2d-1ed676c7b90e" />



## 许可证

本项目根据 MIT 许可证授权。有关详细信息，请参阅 `LICENSE` 文件。