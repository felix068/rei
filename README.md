![REI Banner](./banner.png)

<div align="center">

# REI - RSS Reader

**R**SS **E**xtraordinary Amaz**I**ng reader
A fast, self-hosted RSS feed reader for local usage

</div>

## Features

- 📰 **Subscribe to RSS feeds** - Add any RSS/Atom feed
- 🔄 **Automatic sync** - New posts fetched every 15 minutes
- ✅ **Mark as read** - Track what you've already seen
- 🌓 **Dark/Light theme** - Easy on the eyes, day or night
- 📱 **Responsive** - Works on desktop, tablet, and mobile
- 🔍 **Sort & Filter** - By date, feed name, or alphabetically
- 🐳 **Docker-ready** - Easy deployment
- 🚀 **Minimal & fast** - Simple, no bloat

## Quick Start

### Development

```bash
# Start with Docker Compose
sudo docker compose -f docker-compose.dev.yml up --build
```

Access at **http://localhost:5173**

### Production

```bash
# Start production services
sudo docker compose up --build
```

Access at **http://localhost** (port 80)

## Usage

### Adding Feeds

1. Click "Feeds"
2. Enter RSS feed URL (e.g., `https://news.ycombinator.com/rss`)
3. Click "Add"
4. All posts appear immediately, new ones sync every 15 min

### Reading Posts

- Unread posts appear on home page
- Click "Read" to open article and mark as read
- Sort by date, feed name, or title
- Auto-refresh every 10 seconds

## Architecture

```
rei/
├── api/          # Go backend (Echo framework)
├── web/          # React frontend (TypeScript + Vite)
├── sql/          # Database schema & migrations
└── docker-compose*.yml
```

- **Frontend**: React 19 + TypeScript + Vite + Tailwind CSS
- **Backend**: Go 1.24 + Echo + gofeed
- **Database**: PostgreSQL 15
- **Deployment**: Docker + Docker Compose

## Environment Variables

Create `.env` file:

```env
POSTGRES_USER=reiuser
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=rei_db
DATABASE_URL=postgres://reiuser:your_secure_password@postgres:5432/rei_db?sslmode=disable
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/add_feed` | Add new RSS feed |
| GET | `/list_feeds` | List all feeds |
| GET | `/list_posts` | List all posts |
| GET | `/list_unread_posts` | List unread posts |
| PUT | `/posts/:id/read` | Mark post as read |
| DELETE | `/feeds/:id` | Delete feed |

## Database Migration

Upgrading from older version? Run migration:

```bash
# Connect to postgres
sudo docker exec -it <postgres-container> psql -U $POSTGRES_USER -d $POSTGRES_DB

# Run migration
\i /docker-entrypoint-initdb.d/migrate.sql
```

Or rebuild from scratch:

```bash
sudo docker compose down -v
sudo docker compose up --build
```

## Testing

Use a local RSS feed that generates items every 30 seconds:

```bash
python3 test/local_rss_feed.py
# Server starts at http://localhost:8000/feed.xml
```

From Docker, use: `http://host.docker.internal:8000/feed.xml`

## Documentation

- [API Documentation](./api/README.md)
- [Web Documentation](./web/README.md)

## Contributing

Contributions welcome! Open issues or pull requests.

---

<div align="center">
Made with ❤️ for RSS lovers
</div>
