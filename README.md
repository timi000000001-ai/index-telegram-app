# Telegram Search Platform

This project is a comprehensive Telegram search platform, consisting of a backend service for data collection and a web application for user interaction.

## Project Structure

The project is divided into two main parts:

- `telegram-bot-services/`: A collection of backend services written in Go. These services are responsible for interacting with the Telegram API, collecting messages from groups and channels, and providing a search API.
- `index-telegram-app/`: A Svelte-based web application that provides a user-friendly interface for searching the content collected by the backend services.

## Features

- **Data Collection**: The backend services can log into a Telegram account, join specified groups and channels, and collect messages.
- **Search API**: Provides a powerful search API to query the collected messages, with support for filtering and pagination.
- **Web Interface**: A modern and responsive web interface that allows users to easily search for content, view results, and manage settings.
- **Bot Integration**: The platform includes a Telegram bot that can be used to search for content directly within Telegram.

## Architecture

The backend is composed of three main services:

- **Bot Service**: Handles user interactions with the Telegram bot, provides search functionality, and stores messages.
- **Management Service**: Exposes a REST API for search and session management, acting as a bridge between the frontend and backend services.
- **Collection Service**: Manages Telegram account login, configures target groups, and collects messages for indexing.

The services are designed to work together, leveraging PocketBase for storage and Meilisearch for fast search capabilities.

```
[Svelte Frontend (Static HTML)] ↔ [Management Service (:8080)]
                                    ↕ (REST API)
[Bot Service (:8081)] ↔ [Collection Service (:8082)]
      ↕                           ↕
[PocketBase (Storage)]       [Meilisearch (Search)]
      ↕                           ↕
[S3/MinIO (Media)]
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker (optional)
- Telegram API Credentials
- PocketBase
- Meilisearch
- SvelteKit (optional)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-repo/telegram-search-platform.git
    cd telegram-search-platform
    ```
2.  **Initialize Go modules:**
    ```bash
    cd telegram-bot-services/bot-service && go mod init bot-service && go mod tidy
    cd ../management-service && go mod init management-service && go mod tidy
    cd ../collection-service && go mod init collection-service && go mod tidy
    cd ../..
    ```
3.  **Build the services:**
    ```bash
    go build -o telegram-bot-services/bot-service/bot-service telegram-bot-services/bot-service
    go build -o telegram-bot-services/management-service/management-service telegram-bot-services/management-service
    go build -o telegram-bot-services/collection-service/collection-service telegram-bot-services/collection-service
    ```
4.  **Install frontend dependencies:**
    ```bash
    cd index-telegram-app
    npm install
    ```

### Configuration

Create a `.env` file in each service directory (`bot-service`, `management-service`, `collection-service`) and provide the necessary environment variables. Refer to the `en-README.md` in the `telegram-bot-services` directory for more details.

## Usage

### Running the Services

-   **Bot Service**: `./telegram-bot-services/bot-service/bot-service`
-   **Management Service**: `./telegram-bot-services/management-service/management-service`
-   **Collection Service**: `./telegram-bot-services/collection-service/collection-service`

### Running the Frontend

```bash
cd index-telegram-app
npm run dev
```

## API Documentation

API documentation is available in the `telegram-bot-services/api-docs` directory in the form of Postman collections.

## Legal Notice

**Usage Restriction**: This project is not intended for use in mainland China. Access to Telegram is restricted by the government in mainland China, and the data collection and processing activities of this project may violate local laws and regulations.

**Disclaimer**: The developers of this project are not responsible for any consequences resulting from improper use, violation of local laws, or data privacy issues. Users should assess the legal risks themselves and consult with legal professionals if necessary.

**Recommendation**: If you are located in mainland China, please do not download, install, or run this project. Please look for alternative solutions that comply with local regulations.

## Sponsorship

If you find this project helpful, please consider sponsoring us.

- **TRX & USDT (TRC20):** `TD5JGaR7cY5ZxDnZNgmCSv66axR9DhrcYz`

<img width="512" height="530" alt="image" src="https://github.com/user-attachments/assets/08a5cf87-e174-4bf5-ae2d-1ed676c7b90e" />

Contact on Telegram: https://t.me/simi001001 https://t.me/SoSo00000000001
<img width="939" height="1280" alt="image" src="https://github.com/user-attachments/assets/d7771644-519c-4caf-a669-db026fe25b72" />

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
