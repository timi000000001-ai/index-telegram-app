# Telegram Data Collector and Bot Service

Welcome to the Telegram Data Collector and Bot Service project! This is a modular Go-based application designed to interact with Telegram, collect group/channel messages, and provide a searchable bot interface. The project consists of three main services: Bot Service, Management Service, and Collection Service, integrated with PocketBase, Meilisearch, and a Svelte frontend.

**Last Updated:** Wednesday, September 03, 2025, 01:35 PM WIB  
**License:** MIT (see LICENSE for details)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Development](#development)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Overview

This project enables:

- A Telegram bot to handle user queries with pagination and filtering (Bot Service).
- A management service to provide REST APIs for search and session management.
- A collection service to log into Telegram accounts, configure target groups, and collect messages for indexing.
- Integration with a Svelte frontend for a user-friendly interface.

The services are designed to work together, leveraging PocketBase for storage and Meilisearch for fast search capabilities.

## Features

### Bot Service

- Supports multiple Telegram bots with Webhook integration.
- Provides paginated search results (5 per page) with "Previous" and "Next" buttons.
- Filters search results by Group, Channel, Bot, or All Messages.
- Commands: /help, /clong (clone bot), /sponsor, /mini.
- Stores messages in PocketBase and indexes them in Meilisearch.

[Detailed Bot Service Documentation](./bot-service/README.md)

### Management Service

- Exposes a REST API for search with pagination and filtering.
- Acts as a bridge between the frontend and backend services.

[Detailed Management Service Documentation](./management-service/README.md)

### Collection Service

- Handles Telegram account login with phone number authentication.
- Allows configuration of target groups for message collection.
- Actively collects and parses group messages, storing them for search.

[Detailed Collection Service Documentation](./collection-service/README.md)

## Architecture

```
[Svelte Frontend (Static HTML)] ↔ [Management Service (:8080)]
                                    ↕ (REST API)
[Bot Service (:8081)] ↔ [Collection Service (:8082)]
      ↕                           ↕
[PocketBase (Storage)]       [Meilisearch (Search)]
      ↕                           ↕
[S3/MinIO (Media)]
```

- **Bot Service:** Runs on :8081, handles bot interactions and message storage.
- **Management Service:** Runs on :8080, provides search APIs.
- **Collection Service:** Runs on :8082, manages Telegram data collection.

## Prerequisites

- **Go:** Version 1.21 or higher.
- **Docker:** For containerized deployment (optional).
- **Telegram API Credentials:** Obtain api_id and api_hash from my.telegram.org.
- **PocketBase:** Running instance with a messages collection.
- **Meilisearch:** Running instance with a messages index.
- **SvelteKit:** For building the static frontend (optional).

## Installation

### Clone the Repository

```bash
git clone https://github.com/your-repo/telegram-bot-services.git
cd telegram-bot-services
```

### Initialize Modules

Navigate to each service directory and initialize the Go module:

```bash
cd bot-service && go mod init bot-service && go mod tidy
cd ../management-service && go mod init management-service && go mod tidy
cd ../collection-service && go mod init collection-service && go mod tidy
cd ..
```

### Build the Services

Build each service:

```bash
go build -o bot-service ./bot-service
go build -o management-service ./management-service
go build -o collection-service ./collection-service
```

### Docker Setup (Optional)

Build Docker images:

```bash
docker build -t bot-service ./bot-service
docker build -t management-service ./management-service
docker build -t collection-service ./collection-service
```

## Configuration

### Environment Variables

Create a .env file in each service directory or set environment variables:

#### Bot Service

```
BOT_TOKENS=YOUR_BOT_TOKEN_1,YOUR_BOT_TOKEN_2
POCKETBASE_TOKEN=YOUR_POCKETBASE_TOKEN
MEILISEARCH_KEY=YOUR_MEILISEARCH_KEY
AUTH_TOKEN=YOUR_AUTH_TOKEN
```

#### Management Service

```
MEILISEARCH_KEY=YOUR_MEILISEARCH_KEY
```

#### Collection Service

```
API_ID=123456
API_HASH=your_api_hash
POCKETBASE_URL=http://your-pocketbase-url
POCKETBASE_TOKEN=YOUR_POCKETBASE_TOKEN
MEILISEARCH_URL=http://your-meilisearch-url
MEILISEARCH_KEY=YOUR_MEILISEARCH_KEY
```

### SSL Certificates (Bot Service)

For Webhook support, provide cert.pem and key.pem (e.g., via Let's Encrypt):

```bash
certbot certonly --standalone -d your-bot-service.com
```

## Usage

### Running the Services

#### Bot Service

```bash
./bot-service
```

Or with Docker:

```bash
docker run -p 8081:8081 -v /path/to/cert.pem:/app/cert.pem -v /path/to/key.pem:/app/key.pem bot-service
```

#### Management Service

```bash
./management-service
```

Or with Docker:

```bash
docker run -p 8080:8080 management-service
```

#### Collection Service

```bash
./collection-service
```

Or with Docker:

```bash
docker run -p 8082:8082 collection-service
```

### Interacting with the Bot

- Use commands like /help, /clong, /sponsor, /mini.
- Search with queries (≤10 characters) to see paginated results with filters.

### Configuring Collection

#### Login

```bash
curl -X POST http://localhost:8082/login -H "Content-Type: application/json" -d '{"phone_number": "+1234567890"}'
```

Enter the code sent to your Telegram app.

#### Configure Groups

```bash
curl -X POST http://localhost:8082/configure -H "Content-Type: application/json" -d '{"chat_ids": [123, 456]}'
```

The service will start collecting messages from the specified chat_ids.

### Frontend Integration

Build the Svelte frontend and serve the static HTML:

```bash
cd frontend
npm install
npm run build
```

Serve build/index.html via a web server (e.g., Nginx).

## Development

### Dependencies

- Update go.mod files with go mod tidy.
- Add new dependencies as needed.

### Testing

Run tests for each service:

```bash
go test ./...
```

### Adding Features

- **Bot Service:** Extend commands or add new filters.
- **Collection Service:** Implement real-time updates with gotd Updates API.
- **Management Service:** Add more API endpoints.

## API Endpoints

### Bot Service

- **POST /webhook/{token}:** Handles Telegram Webhook updates.

### Management Service

- **GET /api/search:** Search with q, page, limit, filter parameters.

### Collection Service

- **POST /login:** Login with {"phone_number": "..."}.
- **POST /configure:** Configure with {"chat_ids": [...]}.
- **GET /health:** Health check.

## Contributing

1. Fork the repository.
2. Create a feature branch (git checkout -b feature/awesome-feature).
3. Commit changes (git commit -m "Add awesome feature").
4. Push to the branch (git push origin feature/awesome-feature).
5. Open a Pull Request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

/*
 * File Description: Main README for Telegram Data Collector and Bot Service
 * Main Components: Overview of Bot Service, Management Service, and Collection Service
 * Modification History:
 * @author fcj
 * @date 2025-09-03
 * @version 1.1.0
 * © Telegram Bot Services Team
 */

