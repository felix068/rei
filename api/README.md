# REI API Documentation

Backend REST API for the REI RSS reader, built with Go and Echo framework.

## Tech Stack

- **Go** 1.24
- **Echo** v4 - Web framework
- **gofeed** - RSS/Atom parser
- **sqlx** - SQL toolkit
- **PostgreSQL** - Database

## Project Structure

```
api/
├── cmd/
│   └── server.go       # Main entry point
├── services/
│   └── feeds.go        # Feed & post handlers
├── reader.go           # RSS parser wrapper
├── sync.go             # Sync utilities
├── go.mod              # Dependencies
└── Dockerfile.dev      # Development container
```

## Core Components

### 1. Server (`cmd/server.go`)

Main application entry point:
- Initializes Echo server
- Sets up CORS middleware
- Connects to PostgreSQL
- Registers routes
- Starts background sync goroutine (every 15 min)

### 2. Feed Service (`services/feeds.go`)

Handles all feed and post operations:

**Types:**
- `Feed` - RSS feed metadata
- `Post` - Individual RSS item
- `Handler` - Service with DB and RSS reader

**Endpoints:**

| Method | Route | Function | Description |
|--------|-------|----------|-------------|
| POST | `/add_feed` | `addFeed` | Add new RSS feed, fetch all existing posts |
| GET | `/list_feeds` | `listFeeds` | List all feeds (newest first) |
| GET | `/list_posts` | `listPosts` | List all posts with feed name |
| GET | `/list_unread_posts` | `listUnreadPosts` | List unread posts with feed name |
| PUT | `/posts/:id/read` | `markPostAsRead` | Mark specific post as read |
| DELETE | `/feeds/:id` | `deleteFeed` | Delete feed and all its posts |

**Internal Functions:**
- `SyncFeed(feedId, rssLink)` - Fetch new posts from one feed
- `SyncAllFeeds()` - Sync all feeds in database

### 3. RSS Reader (`reader.go`)

Wrapper around gofeed parser:
- `NewParser()` - Create new reader instance
- `ReadFeed(url)` - Parse RSS/Atom feed from URL

### 4. Sync (`sync.go`)

Simple sync utility (currently just calls ReadFeed).

## Database Schema

### feeds table
```sql
id          UUID PRIMARY KEY
name        TEXT NOT NULL
description TEXT
link        TEXT UNIQUE
rss_link    TEXT UNIQUE
created_at  TIMESTAMP NOT NULL
updated_at  TIMESTAMP NOT NULL
```

### posts table
```sql
id          UUID PRIMARY KEY
name        TEXT NOT NULL
description TEXT
link        TEXT UNIQUE
feed_id     UUID NOT NULL
is_read     BOOLEAN NOT NULL DEFAULT false
created_at  TIMESTAMP NOT NULL
updated_at  TIMESTAMP NOT NULL
```

## Key Features

### Automatic Sync
Background goroutine syncs all feeds every 15 minutes:
1. Fetches all feeds from database
2. For each feed, fetches RSS and checks for new posts
3. Inserts only posts that don't exist (by unique link)
4. Updates feed's `updated_at` timestamp

### Duplicate Prevention
- Feeds: Checked by `rss_link` before insert
- Posts: Unique constraint on `link` column

### Error Handling
- Returns appropriate HTTP status codes
- Logs errors to console
- Continues processing on individual post failures

## Environment Variables

Required:
```bash
DATABASE_URL=postgres://user:pass@host:port/dbname?sslmode=disable
```

## Development

### Run with live reload (wgo):
```bash
cd api
wgo run ./cmd/server.go
```

### Build:
```bash
go build -o server ./cmd/server.go
```

### Docker:
```bash
docker build -f Dockerfile.dev -t rei-api .
docker run -p 1323:1323 --env-file ../.env rei-api
```

## API Usage Examples

### Add feed:
```bash
curl -X POST http://localhost:1323/add_feed \
  -H "Content-Type: application/json" \
  -d '{"link": "https://news.ycombinator.com/rss"}'
```

### List unread posts:
```bash
curl http://localhost:1323/list_unread_posts
```

### Mark post as read:
```bash
curl -X PUT http://localhost:1323/posts/{post-id}/read
```

### Delete feed:
```bash
curl -X DELETE http://localhost:1323/feeds/{feed-id}
```

## Notes

- CORS enabled for all origins (`*`)
- Sync waits 30 seconds on startup before first run
- Posts include `feed_name` via SQL JOIN
- Feed deletion cascades to posts (manual DELETE in code)
