# Telegram 服务 API 文档

## 概述

本目录包含 Telegram 数据采集和机器人服务项目的 API 接口文档，以 Postman 集合格式提供，方便前端开发和测试人员使用。

## 文件说明

- **bot-service-postman.json**：机器人服务 API 集合
- **management-service-postman.json**：管理服务 API 集合
- **collection-service-postman.json**：采集服务 API 集合
- **telegram-services-environment.json**：环境变量配置

## 使用方法

### 导入到 Postman

1. 打开 Postman 应用
2. 点击左上角的 "Import" 按钮
3. 选择 "File" 选项卡，然后上传本目录中的 JSON 文件
4. 导入环境配置文件 `telegram-services-environment.json`
5. 从右上角的环境下拉菜单中选择 "Telegram 服务环境配置"

### 配置环境变量

根据您的实际部署情况，修改以下环境变量：

- **bot_service_url**：机器人服务的URL（默认：http://localhost:8081）
- **management_service_url**：管理服务的URL（默认：http://localhost:8080）
- **collection_service_url**：采集服务的URL（默认：http://localhost:8082）
- **bot_token**：您的 Telegram 机器人令牌

### 使用 API

#### 机器人服务

- **POST /webhook/{token}**：处理 Telegram Webhook 更新

#### 管理服务

- **GET /api/search**：搜索接口，支持分页和过滤

#### 采集服务

- **POST /login**：登录 Telegram 账户（发送验证码或提交验证码）
- **POST /configure**：配置目标群组
- **GET /health**：健康检查

## 注意事项

- 确保服务已正常运行并可访问
- 使用前请先配置正确的环境变量
- 对于需要认证的接口，请确保提供有效的令牌

## 示例流程

### 采集服务使用流程

1. 发送登录请求，提供电话号码
2. 接收验证码后，提交验证码完成登录
3. 配置目标群组进行消息采集
4. 使用健康检查接口监控服务状态

### 管理服务使用流程

1. 使用搜索接口查询消息
2. 通过调整查询参数进行分页和过滤

/*
 * 文件功能描述：Telegram 服务 API 文档说明
 * 主要内容：Postman 集合使用指南和 API 接口说明
 * 修改历史记录：
 * @author fcj
 * @date 2023-07-18
 * @version 1.0.0
 * © Telegram Bot Services Team
 */