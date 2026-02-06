# Docker Compose Configuration Documentation

This document provides a detailed explanation of the Docker Compose setup for the `smanzary_site` project. The project consists of four interconnected services that form a full-stack media management application.

## 🚀 Overview

The architecture follows a microserviced approach where each component is isolated in its own container:

1.  **Database (`postgres`)**: Persistent storage for all application data.
2.  **Backend (`backend`)**: Go-based API handling logic, database interactions, and authentication.
3.  **Frontend (`frontend`)**: React-based Single Page Application (SPA) served via Nginx.
4.  **Thumbnailer (`thumbnailer`)**: Background worker for generating media thumbnails.

---

## 🛠 Services

### 1. `postgres`
-   **Image**: `postgres:16-alpine` (Lightweight PostgreSQL 16).
-   **Purpose**: Stores users, media metadata, and application state.
-   **Key Configs**:
    -   **Environment**: Uses `POSTGRES_USER`, `POSTGRES_PASSWORD`, and `POSTGRES_DB` from the `.env` file.
    -   **Persistence**: Mounts `./smanzy_data/postgres` to `/var/lib/postgresql/data`.
    -   **Port Mapping**: `127.0.0.1:5432:5432`.
    -   **Healthcheck**: Uses `pg_isready` to ensure the database is ready before other services start.

### 2. `backend`
-   **Build**: Located in `./smanzy_backend`.
-   **Purpose**: The core Go server (`server`).
-   **Key Configs**:
    -   **Command**: `./server -migrate` (automatically runs database migrations on startup).
    -   **Environment**:
        -   `DB_DSN`: Dynamically constructed using database credentials.
        -   `JWT_SECRET`: Used for authentication tokens.
        -   `UPLOAD_DIR`: Pointed to `/app/uploads` inside the container.
    -   **Volume**: Mounts local `./smanzy_data/uploads` to `/app/uploads`.
    -   **Dependencies**: Depends on `postgres` being healthy.
    -   **Port Mapping**: `127.0.0.1:8080:8080`.

### 3. `frontend`
-   **Build**: Located in `./smanzy_react_spa`.
-   **Purpose**: Serves the React frontend using Nginx.
-   **Key Configs**:
    -   **Dependencies**: Depends on `backend` being healthy.
    -   **Port Mapping**: `127.0.0.1:3000:80` (access the app via `http://localhost:3000`).
    -   **Healthcheck**: Checks if `http://localhost:80/` (internal nginx port) is responding.

### 4. `thumbnailer`
-   **Build**: Located in `./smanzy_thumbgen`.
-   **Purpose**: A Go-based worker that watches the `uploads` directory and generates thumbnails for images and videos.
-   **Key Configs**:
    -   **Command**: `./thumbgen` (runs in watcher mode).
    -   **Volume**: Shares the same `./smanzy_data/uploads` volume as the backend to process files as they are uploaded.
    -   **Dependencies**: Depends on `backend` being healthy (to ensure the environment is ready).

---

## 🌐 Networking

All services communicate over a private bridge network named `smanzy_network`.

-   Services can reach each other using their container names as hostnames:
    -   Backend connects to database at `postgres:5432`.
    -   Frontend communicates with backend at `backend:8080` (internally).

---

## 💾 Data Persistence (Volumes)

All persistent data is stored in the `./smanzy_data` directory to ensure data survives container restarts and updates:

-   `./smanzy_data/postgres`: Database files.
-   `./smanzy_data/uploads`: Original media files and generated thumbnails.

---

## 🔑 Environment Variables (`.env`)

The system relies on a `.env` file for sensitive data and configuration. Ensure these variables are set:

| Variable | Description |
| :--- | :--- |
| `POSTGRES_USER` | Database username |
| `POSTGRES_PASSWORD` | Database password |
| `POSTGRES_DB` | Primary database name |
| `JWT_SECRET` | Secret key for signing authentication tokens |
| `YOUTUBE_API_KEY` | (Optional) Key for YouTube integrations |
| `ENV` | Set to `development` or `production` |

---

## 🚦 Startup Logic

The configuration uses Docker's `healthcheck` and `depends_on` (with conditions) to ensure a stable startup sequence:

1.  **Postgres** starts first.
2.  **Backend** waits for Postgres to be healthy, then runs migrations and starts the API.
3.  **Frontend** and **Thumbnailer** wait for the Backend API to be healthy before starting their processes.

## 🏃 How to Run

To start the entire stack:
```bash
docker-compose up -d
```

To stop:
```bash
docker-compose down
```

### Production vs Development

Currently, `docker-compose.yml` and `docker-compose.prod.yml` are identical in this project. In a typical workflow:
- `docker-compose.yml`: Used for local development (mapped to `127.0.0.1`).
- `docker-compose.prod.yml`: Can be used for production-specific overrides (e.g., different restart policies, environment-specific images, or non-localhost port bindings).

To run specifically with the production file:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

To view logs:
```bash
docker-compose logs -f
```
