# REI Web Frontend Documentation

React-based web interface for the REI RSS reader.

## Tech Stack

- **React** 19.1
- **TypeScript** 5.8
- **Vite** 6.3 - Build tool
- **Tailwind CSS** 4.1 - Styling
- **React Router** 7.6 - Routing

## Project Structure

```
web/
├── src/
│   ├── layouts/
│   │   └── main-layout.tsx    # Main app layout with nav
│   ├── pages/
│   │   └── Feeds.tsx           # Feeds management page
│   ├── utils/
│   │   └── api.ts              # API client
│   ├── Home.tsx                # Unread posts page
│   ├── main.tsx                # App entry point
│   └── index.css               # Global styles
├── public/
│   └── logo/                   # Logo assets
├── package.json
└── vite.config.ts
```

## Components

### MainLayout (`layouts/main-layout.tsx`)

Main application shell:
- **Navigation bar** with logo and page links
- **Theme toggle** (dark/light mode)
- **Unread counter** (updates every 10 seconds)
- **Responsive** design for mobile

Features:
- Persists theme choice in `localStorage`
- Switches logo variant based on theme
- Auto-refreshes unread count

### Home (`Home.tsx`)

Unread posts page:
- Displays all unread posts
- **Sort options**: Date, Feed name, Alphabetical
- **Feed badge** showing source of each post
- **Mark as read** button (opens link + marks read)
- Auto-refreshes every 10 seconds
- Fully responsive

### Feeds (`pages/Feeds.tsx`)

Feed management page:
- **Add feed** form with URL input
- **Feed list** with details
- **Delete button** for each feed
- Error display for invalid feeds
- Responsive layout

## API Client (`utils/api.ts`)

TypeScript API wrapper with interfaces:

```typescript
interface Feed {
  id: string;
  name: string;
  description: string;
  link: string;
  rssLink: string;
  createdAt: string;
  updatedAt: string;
}

interface Post {
  id: string;
  name: string;
  description: string;
  link: string;
  feedId: string;
  feedName: string;
  isRead: boolean;
  createdAt: string;
  updatedAt: string;
}
```

Functions:
- `addFeed(url)` - Add RSS feed
- `listFeeds()` - Get all feeds
- `listPosts()` - Get all posts
- `listUnreadPosts()` - Get unread posts
- `markPostAsRead(id)` - Mark post as read
- `deleteFeed(id)` - Delete feed

## Routing

```typescript
/           → Home (unread posts)
/feeds      → Feeds management
```

## Styling

### Theme System

Two themes with localStorage persistence:

**Light theme:**
- Background: `#ffffff`
- Text: `#000000`
- Cards: `#ffffff`
- Accent: `rgb(237, 108, 61)` (orange from logo)

**Dark theme:**
- Background: `#1a1a1a` (dark gray, like GitHub)
- Text: `#ffffff`
- Cards: `#2d2d2d`
- Borders: `#4d4d4d`
- Accent: `rgb(237, 108, 61)`

### Responsive Breakpoints

Tailwind defaults:
- `sm`: 640px
- `md`: 768px
- `lg`: 1024px

Mobile-first approach with `sm:` and `md:` prefixes.

## Features

### Sort & Filter
- Date (newest first) - default
- Feed name (A-Z)
- Title (A-Z)

### Auto-refresh
- Unread posts: Every 10 seconds
- Unread counter: Every 10 seconds
- Theme sync: Every 100ms (between tabs)

### Responsive Design
- Mobile: Single column, smaller text
- Tablet: Optimized spacing
- Desktop: Full width (max 1280px)

## Development

### Install:
```bash
npm install
```

### Run dev server:
```bash
npm run dev
# Opens http://localhost:5173
```

### Build:
```bash
npm run build
# Output to dist/
```

### Lint:
```bash
npm run lint
```

## Configuration

### API Base URL
Auto-detected based on mode:
- Development: `http://localhost:1323`
- Production: `/api` (proxied by nginx)

### Vite Proxy
Dev server proxies `/api` → `http://localhost:1323`

## Docker

Development:
```dockerfile
FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
```

Production:
```dockerfile
# Multi-stage build
# 1. Build static files
# 2. Serve with nginx + API proxy
```

## Browser Support

Modern browsers with ES2020+ support:
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Notes

- Uses inline styles for theme to override browser extensions (DarkReader)
- Logo switches automatically with theme
- All state is ephemeral (no localStorage except theme)
- Post descriptions are truncated in UI
