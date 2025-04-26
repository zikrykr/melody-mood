# 🎧 MelodyMood

MelodyMood is a stateless, Golang-powered backend service that generates personalized music recommendations using OpenAI and Spotify APIs. It also uses Redis for caching and tracking usage limits.

---

## 🚀 Features

- 🎯 Personalized music recommendations based on:
  - Personality
  - Favorite music genre
  - Typical listening occasion
- 🤖 AI-generated playlist suggestions (via OpenAI)
- 🎵 Fetch real Spotify tracks and album cover art
- ⚡ Caching recommendations and playlists per session (with Redis)
- 🐳 Full Dockerized setup (no manual server installation)
- 🔥 Session-based rate-limiting (e.g., 20 requests/hour)
- 📦 Stateless design — no database needed

---

## 📊 System Architecture

```mermaid
flowchart TD
  A[User Input: Personality, Genre, Occasion]
  B[MelodyMood Backend - Golang]
  C[OpenAI API - LLM]
  D[Spotify API]
  R[Redis]

  A -->|POST /recommendations| B

  subgraph External Services
    C
    D
  end

  B -->|Send Prompt| C
  C -->|Return Recommendations| B
  B -->|Fetch Track Info| D
  B -->|Store Session Data| R
  B -->|Display Songs| A

  A -->|Select 1-5 Songs| B
  B -->|POST /playlists| C
  B -->|Create Playlist| D
  B -->|Track Usage Counts| R
```

---

## Deployment Diagram

![Deployment Diagram](deploy-diagram.png)

## 🐳 Running Locally with Docker

### 1. Clone the repository

```bash
git clone https://github.com/your-username/melodymood.git
cd melodymood
```

### 2. Create a `.env` file

```env
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
OPENAI_API_KEY=your_openai_api_key
REDIRECT_URI=http://localhost:8080/callback

APP_ENV=dev
APP_PORT=8080
APP_NAME=melody-mood-service

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
```

### 3. Docker Compose Up

```bash
docker-compose up --build
```

This will start:

- MelodyMood Backend
- Redis Server

You can access the backend at:

```
http://localhost:8080
```

---

## 📦 Redis Usage

| Key Format                                  | Purpose                            | TTL     |
| :------------------------------------------ | :--------------------------------- | :------ |
| `session:{session_id}:recommendation_count` | Count of user recommendations      | 1 hour  |
| `session:{session_id}:recommendations`      | Cached recommendations per session | 2 hours |
| `session:{session_id}:playlist`             | Cached playlist per session        | 2 hours |

✅ **`is_regenerate: false`** → Try to serve from Redis cache first  
✅ **`is_regenerate: true`** → Force re-generate from OpenAI and overwrite cache

---

## 📡 API Endpoints

### `POST /api/v1/recommendations`

**Request Body**:

```json
{
  "personality": "calm and introverted",
  "genre": "lo-fi hip hop",
  "occasion": "working late at night",
  "is_regenerate": false
}
```

**Headers**:

- `X-Session-ID: your-unique-session-id`

**Behavior**:

- If `is_regenerate: false`, tries Redis cache first.
- If no cache found, generates new recommendations.
- If `is_regenerate: true`, always generates new and updates cache.

**Response Headers**:

- `X-RateLimit-Remaining: number` — number of allowed requests remaining this hour.

---

### `POST /api/v1/playlists`

**Request Body**:

```json
{
  "picked_songs": [
    {
      "song_name": "Blinding Lights",
      "song_artist": "The Weeknd"
    }
  ],
  "genre": "pop",
  "is_regenerate": false
}
```

**Headers**:

- `X-Session-ID: your-unique-session-id`

**Behavior**:

- If `is_regenerate: false`, tries Redis cache first.
- If no cache found, generates a new playlist.
- If `is_regenerate: true`, always generates a new playlist and updates cache.

---

### `POST /api/v1/playlists/spotify`

**Request Body**:

```json
{
  "name": "My Custom Playlist"
}
```

**Headers**:

- `X-Session-ID: your-unique-session-id`

**Behavior**:

- Creates a Spotify playlist for the user using the generated playlist tracks.

---

## 📄 License

MIT License. See `LICENSE` file.

---

## ❤️ Credits

Built with love using:

- [Go](https://golang.org/)
- [Spotify Web API](https://developer.spotify.com/documentation/web-api/)
- [OpenAI](https://platform.openai.com/)

---

> 🎶 _MelodyMood – Music that feels you._
