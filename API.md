# API Documentation

Base URL: `http://localhost:9090`

All responses are JSON unless noted. CORS is open (`*`) for all origins.

---

## Music Endpoints

### `GET /`

Health check — confirms the gateway is running.

**Response `200`**
```json
{ "status": "gateway is running" }
```

---

### `POST /music` — Upload Music

Upload an audio file. The file is streamed to `music-service` in 32 KB chunks over gRPC client-streaming and saved to the shared storage volume.

**Request**

`Content-Type: multipart/form-data`

| Field | Type | Required | Description |
|---|---|---|---|
| `file` | file | ✅ | Audio file to upload. Must be a valid audio MIME type and within the size limit. |

**Response `200`**
```json
{ "status": "uploaded successfully" }
```

**Error Responses**

| Status | Reason |
|---|---|
| `400` | No file in request, or file fails validation (wrong MIME type / too large) |
| `500` | Failed to open file, create gRPC stream, read chunks, or receive final response |

**Example**
```bash
curl -X POST http://localhost:9090/music \
  -F "file=@/path/to/song.mp3"
```

---

### `GET /music` — List All Music

Returns metadata for every track stored in `music_db`.

**Response `200`**
```json
[
  { "id": "uuid-1", "filename": "song_name.mp3" },
  { "id": "uuid-2", "filename": "another_song.mp3" }
]
```

**Error Responses**

| Status | Reason |
|---|---|
| `500` | gRPC call to `music-service` failed |

**Example**
```bash
curl http://localhost:9090/music
```

---

### `GET /music/:id` — Stream Music

Streams the audio file for a given track ID. The response is a full audio byte stream and supports **HTTP Range requests**, enabling seeking/scrubbing in browser `<audio>` players.

**Path Parameters**

| Parameter | Type | Description |
|---|---|---|
| `id` | string (UUID) | ID returned from `GET /music` |

**Response Headers**

```
Content-Type: audio/mpeg
Accept-Ranges: bytes
Content-Disposition: inline
```

**Response `200`**

Binary audio stream (or partial `206` content if `Range` header is provided by the client).

**Error Responses**

| Status | Reason |
|---|---|
| `500` | gRPC stream could not be opened or a chunk failed to receive |

**Example**
```bash
# Direct download
curl -o song.mp3 http://localhost:9090/music/<id>

# With range (e.g. seek to byte 500000)
curl -H "Range: bytes=500000-" http://localhost:9090/music/<id>
```

---

## Lyrics Endpoints

### `POST /lyrics` — Add / Transcribe Lyrics

Triggers AI transcription for a track. Internally:

1. Checks if lyrics for this music already exist — returns immediately if so (idempotent).
2. Pulls the full audio from `music-service` via gRPC streaming.
3. Forwards the audio bytes to `transcription-service` (Whisper `medium` model).
4. Saves the resulting timestamped segments and detected language to `lyrics_db`.

> ⚠️ This is a **slow operation** — transcription can take 10–120 seconds depending on track length and hardware. It should be called once per track.

**Request**

`Content-Type: application/json`

| Field | Type | Required | Description |
|---|---|---|---|
| `music_id` | string (UUID) | ✅ | ID of the track to transcribe |
| `text` | string | ✅ | Human-readable title / label for the track (used for duplicate detection) |

**Request Body Example**
```json
{
  "music_id": "a1b2c3d4-...",
  "text": "I Will Survive"
}
```

**Response `200`**
```json
{
  "success": true,
  "music_id": "a1b2c3d4-...",
  "message": "Lyrics Added successfully"
}
```

**Error Responses**

| Status | Reason |
|---|---|
| `400` | Missing or invalid JSON body (`music_id` or `text` not provided) |
| `502` | gRPC call to `lyrics-service` failed (downstream error, e.g. Whisper unavailable) |

**Example**
```bash
curl -X POST http://localhost:9090/lyrics \
  -H "Content-Type: application/json" \
  -d '{ "music_id": "a1b2c3d4-...", "text": "I Will Survive" }'
```

---

### `GET /lyrics/:id` — Get Lyrics

Retrieves stored timestamped lyrics for a track. Each segment has a start and end time (in seconds) that can be used to sync lyrics to audio playback.

**Path Parameters**

| Parameter | Type | Description |
|---|---|---|
| `id` | string (UUID) | Music ID (same as used in `POST /lyrics`) |

**Response `200`**
```json
{
  "lyrics": [
    { "start": 0.0,  "end": 3.5,  "text": "At first I was afraid" },
    { "start": 3.5,  "end": 7.2,  "text": "I was petrified" },
    { "start": 7.2,  "end": 12.0, "text": "Kept thinking I could never live without you by my side" }
  ],
  "language": "en"
}
```

**Response Fields**

| Field | Type | Description |
|---|---|---|
| `lyrics` | array | Ordered list of lyric segments |
| `lyrics[].start` | float | Segment start time in seconds |
| `lyrics[].end` | float | Segment end time in seconds |
| `lyrics[].text` | string | Transcribed text for this segment |
| `language` | string | ISO 639-1 language code detected by Whisper |

**Error Responses**

| Status | Reason |
|---|---|
| `500` | gRPC call to `lyrics-service` failed or lyrics not found |

**Example**
```bash
curl http://localhost:9090/lyrics/<music-id>
```

---

## Internal Service Interfaces (gRPC)

These are not exposed publicly — documented here for development reference.

### MusicService (`:50051`)

Defined in `gateway/proto/music.proto`

| RPC | Request | Response | Type | Description |
|---|---|---|---|---|
| `UploadMusic` | stream `UploadMusicChunks` | `UploadMusicResponse` | Client streaming | Upload file in 32 KB chunks; filename sent in first chunk only |
| `ListMusic` | `Empty` | `ListResponse` | Unary | Returns all stored tracks |
| `StreamMusic` | `StreamRequest { id }` | stream `MusicChunk` | Server streaming | Streams audio bytes for a given track ID |

### LyricsService (`:50052`)

Defined in `gateway/proto/lyrics.proto`

| RPC | Request | Response | Type | Description |
|---|---|---|---|---|
| `AddLyrics` | `AddLyricsRequest { music_id, text }` | `Empty` | Unary | Triggers transcription and stores results (idempotent) |
| `GetLyrics` | `GetLyricsRequest { music_id }` | `LyricsResponse` | Unary | Returns timestamped lyrics from DB |

### TranscriptionService (`:5001`)

Python/Flask HTTP service (internal only).

| Endpoint | Method | Description |
|---|---|---|
| `POST /transcribe` | `multipart/form-data` with `file` field | Runs Whisper transcription; returns `{ lyrics, language, duration }` |
| `GET /health` | — | Health check used by Docker Compose; returns `{ status: "ok" }` |
