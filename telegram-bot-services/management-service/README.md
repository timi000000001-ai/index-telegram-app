# Telegram 管理服务

## 功能概述

**Telegram 管理服务**是 Telegram 数据采集和机器人服务项目的重要组件，作为前端与后端服务之间的桥梁，提供 REST API 用于搜索和会话管理。

### 主要功能

- **搜索 API**：提供带分页和过滤功能的搜索接口
- **会话管理**：处理用户会话和认证
- **前后端桥接**：连接 Svelte 前端与后端服务
- **数据访问**：与 Meilisearch 集成，提供高效搜索能力

## 技术架构

管理服务运行在端口 `:8080`，与其他服务协同工作：

```
[Svelte 前端 (静态 HTML)] ↔ [管理服务 (:8080)] ↔ [Meilisearch (搜索)]
                                    ↕
                          [机器人服务 & 采集服务]
```

## 配置说明

### 环境变量

创建 `.env` 文件或设置以下环境变量：

```
MEILISEARCH_KEY=YOUR_MEILISEARCH_KEY  # Meilisearch 认证密钥
```

## 安装与运行

### 直接运行

```bash
# 初始化模块
go mod init management-service && go mod tidy

# 构建服务
go build -o management-service .

# 运行服务
./management-service
```

### Docker 运行

```bash
# 构建 Docker 镜像
docker build -t management-service .

# 运行容器
docker run -p 8080:8080 management-service
```

## API 端点

- **GET /api/search**：搜索接口，支持以下参数：
  - `q`：搜索查询
  - `page`：页码
  - `limit`：每页结果数
  - `filter`：过滤类型（群组、频道、机器人或全部）

## 开发指南

### 添加新功能

- 扩展 API 端点
- 优化搜索算法
- 增强会话管理功能

### 测试

```bash
go test ./...
```

## 与前端集成

管理服务设计为与 Svelte 前端无缝集成：

```bash
# 构建前端
cd frontend
npm install
npm run build

# 通过 web 服务器（如 Nginx）提供 build/index.html
```

## 注意事项

- 确保 Meilisearch 服务正常运行
- 验证 API 端点的安全性和性能
- 监控服务负载和响应时间

/*
 * 文件功能描述：管理服务 README
 * 主要类/接口说明：提供管理服务的安装、配置和使用说明
 * 修改历史记录：
 * @author fcj
 * @date 2025-09-03
 * @version 1.0.0
 * © Telegram Bot Services Team
 */