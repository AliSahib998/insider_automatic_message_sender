# Insider Message Service

A Golang-based microservice that automatically sends scheduled messages via webhook, persists state in PostgreSQL, and caches metadata in Redis.

## 🚀 Features

- Sends 2 unsent messages every 2 minutes
- Truncates message content to 160 characters
- Caches `messageId + timestamp` in Redis
- Swagger API documentation
- Dockerized with PostgreSQL and Redis

## 🧰 Tech Stack

- Go 1.23+
- Gin Web Framework
- PostgreSQL
- Redis
- Docker & Docker Compose
- Swagger (Swaggo)

## 🛠 Getting Started

### Prerequisites

- Docker & Docker Compose installed

### 🚀 Run with Docker

```bash
docker compose up -d

```

## 📘 API Documentation (Swagger)

Once the app is running, access:

👉 `http://localhost:3015/api/v1/messages/swagger/index.html#/`

## 📡 API Endpoints

| Method | Path                       | Description                           |
|--------|----------------------------|---------------------------------------|
| `POST` | `/api/v1/messages/start`   | Start the automatic sender            |
| `POST` | `/api/v1/messages/stop`    | Stop the automatic sender             |
| `POST` | `/api/v1/messages/`        | Create new message will be sent later |
| `GET`  | `/api/v1/messages/sent`    | List sent messages                    |
