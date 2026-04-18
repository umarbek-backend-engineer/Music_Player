# 🎵 Music Player

An AI-powered music streaming platform built as a microservice system. Upload tracks, stream them in the browser, and get **AI-generated timestamped lyrics** via OpenAI Whisper — all wired together with gRPC, Go, Python, and Docker.

---

## ✨ Features

- **Upload music** — stream audio files directly to storage via chunked gRPC streaming
- **Stream music** — browser-compatible audio streaming with `Range` request support (seek / scrub)
- **AI Lyrics** — automatic lyric transcription using OpenAI Whisper (`medium` model), stored with per-segment timestamps
- **Timestamped display** — lyrics sync to playback position in real time
- **Microservice architecture** — services communicate over gRPC; only the gateway exposes HTTP

---

## 🏗️ Architecture

```
Client (Browser)
    │  HTTP REST
    ▼
┌─────────┐        gRPC        ┌───────────────┐       SQL      ┌──────────┐
│ Gateway │ ──────────────────▶│ music-service │ ─────────────▶│ Postgres │
│  :9090  │                    │    :50051     │               │  :5433   │
│         │        gRPC        └───────────────┘               └──────────┘
│         │ ──────────────────▶┌───────────────┐       SQL           ▲
└─────────┘                    │lyrics-service │ ────────────────────┘
                               │    :50052     │
                               └───────┬───────┘
                                       │  HTTP POST /transcribe
                                       ▼
                               ┌───────────────────────┐
                               │ transcription-service │
                               │    Whisper (medium)   │
                               │        :5001          │
                               └───────────────────────┘
```

### Services at a glance

| Service | Language | Transport | Port (host→container) | Purpose |
|---|---|---|---|---|
| `gateway` | Go / Gin | HTTP REST | `9090 → 8080` | Public-facing API |
| `music-service` | Go | gRPC | `50051 → 50051` | Music metadata + file storage |
| `lyrics-service` | Go | gRPC | `50052 → 50052` | Lyrics storage + transcription orchestration |
| `transcription-service` | Python / Flask | HTTP | `5001 → 5001` | Whisper AI transcription |
| `postgres` | PostgreSQL | SQL | `5433 → 5432` | Persistent storage (2 databases) |

---

## 🚀 Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) + [Docker Compose](https://docs.docker.com/compose/)
- *(Optional)* NVIDIA GPU + Docker GPU support for faster transcription

### Run

```bash
git clone <repo-url>
cd Music_Player
docker-compose up --build
```

This will automatically:
1. Start PostgreSQL
2. Run DB migrations for both `music_db` and `lyrics_db`
3. Start `transcription-service` (health-checked at `/health` before dependants start)
4. Start `music-service`, `lyrics-service`, and `gateway`

**Gateway is available at:** `http://localhost:9090`  
**Frontend is available at:** `http://localhost:3000`  
**Transcription health check:** `http://localhost:5001/health`

If your API runs on a different host/port, open the frontend with an override:

```text
http://localhost:3000/?api=http://localhost:9090
```

> ⚠️ On first run, `transcription-service` may take several minutes to download and cache the Whisper `medium` model.

### No GPU?

Remove or comment out the GPU reservation block in `docker-compose.yml` under `transcription-service`:

```yaml
# deploy:
#   resources:
#     reservations:
#       devices:
#         - driver: nvidia
#           count: 1
#           capabilities: [gpu]
```

---

## 🗄️ Database Migrations

Migrations run automatically on `docker-compose up`. To run them manually:

```bash
# Music DB — up
docker-compose run migrator-music \
  -path=/migrations \
  -database=postgresql://docker:1241@postgres:5432/music_db?sslmode=disable up

# Lyrics DB — up
docker-compose run migrator-lyrics \
  -path=/migrations \
  -database=postgresql://docker:1241@postgres:5432/lyrics_db?sslmode=disable down

# Roll back
docker-compose run migrator-music  -path=/migrations -database=postgresql://docker:1241@postgres:5432/music_db?sslmode=disable down
docker-compose run migrator-lyrics -path=/migrations -database=postgresql://docker:1241@postgres:5432/lyrics_db?sslmode=disable down
```

---

## ⚙️ Environment Variables

### `gateway/.env`

| Variable | Default | Description |
|---|---|---|
| `API_PORT` | `8080` | Gateway listen port (inside container) |
| `API_HOST` | `gateway` | Gateway hostname |
| `GRPC_HOST` | `music-service` | Music service gRPC host |
| `GRPC_MUSIC_SERVICE_PORT` | `50051` | Music service gRPC port |
| `GRPC_LYRICS_SERVICE_HOST` | `lyrics-service` | Lyrics service gRPC host |
| `GRPC_LYRICS_SERVICE_PORT` | `50052` | Lyrics service gRPC port |

### `music-service/.env`

| Variable | Default | Description |
|---|---|---|
| `API_PORT` | `50051` | gRPC server port |
| `DB_HOST` | `postgres` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `docker` | Database user |
| `DB_PASSWORD` | `1241` | Database password |
| `DB_NAME` | `music_db` | Database name |
| `STORAGEPATH` | `/app/storage` | File storage mount path |
| `NETWORK_PROTOCOL` | `tcp` | Network protocol |

### `lyrics-service/.env`

| Variable | Default | Description |
|---|---|---|
| `API_PORT` | `50052` | gRPC server port |
| `DB_HOST` | `postgres` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `docker` | Database user |
| `DB_PASSWORD` | `1241` | Database password |
| `DB_NAME` | `lyrics_db` | Database name |
| `NETWORK_PROTOCOL` | `tcp` | Network protocol |

---

## 📁 Project Structure

```
Music_Player/
├── gateway/                    # Go HTTP gateway (Gin)
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── config/             # Env config loader
│   │   ├── grpc_init/          # gRPC client init (music + lyrics)
│   │   ├── handler/            # HTTP handlers
│   │   ├── modules/            # Request/response models
│   │   └── router/             # Route definitions
│   └── proto/                  # Protobuf definitions + generated code
│
├── music-service/              # Go gRPC music service
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── model/
│   │   ├── repository/         # Postgres queries
│   │   └── service/            # gRPC service implementation
│   ├── migrator/               # DB migration runner
│   └── proto/
│
├── lyrics-service/             # Go gRPC lyrics service
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── model/
│   │   ├── repository/         # Postgres queries
│   │   └── service/            # gRPC + Whisper orchestration
│   ├── migrator/
│   ├── pkg/utils/wisper.go     # Whisper HTTP client
│   └── proto/
│
├── transcription-service/      # Python Flask + Whisper
│   ├── app.py
│   ├── Dockerfile
│   └── requirements.txt
│
├── frontend/                   # Static HTML frontend
│   └── index.html
│
├── postgres/
│   └── init.sql                # DB init (creates music_db + lyrics_db)
│
├── storage/                    # Shared volume for audio files
└── docker-compose.yml
```

---

## 📝 Notes

- **RabbitMQ** variables exist in `.env` files but RabbitMQ is not currently wired into `docker-compose.yml`. This is scaffolding for a future async upload pipeline.
- **Whisper model cache** is persisted in a Docker volume (`whisper_cache`) so it is not re-downloaded on every restart.
- The `transcription-service` forces `language="en"` and uses `fp16=False` (CPU-safe). Switch to `fp16=True` for GPU inference.
- `StreamMusic` uses `http.ServeContent` which handles HTTP `Range` requests, enabling seeking/scrubbing in the browser audio player.