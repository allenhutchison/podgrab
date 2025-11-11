# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## About Podgrab

Podgrab is a self-hosted podcast manager written in Go that automatically downloads podcast episodes. It uses a web interface with server-side rendered HTML templates and WebSocket communication for real-time updates.

## Build and Run Commands

### Local Development
```bash
# Build the application
go build -o ./app ./main.go

# Run locally (requires .env file with CONFIG, DATA, and CHECK_FREQUENCY variables)
./app

# The server runs on port 8080 by default (configurable via PORT environment variable)
```

### Docker
```bash
# Build Docker image
docker build -t podgrab .

# Run with Docker
docker run -d -p 8080:8080 --name=podgrab -v "/host/path/to/config:/config" -v "/host/path/to/assets:/assets" allenhutchison/podgrab

# Run with Docker Compose
docker-compose up -d
```

### Environment Variables
- `CONFIG`: Path to config directory (contains podgrab.db SQLite database)
- `DATA`: Path to data/assets directory (podcast episode files)
- `CHECK_FREQUENCY`: How often to check for new episodes in minutes (default: 30)
- `PASSWORD`: If set, enables basic auth with username "podgrab"
- `PORT`: Change internal application port (default: 8080)
- `GIN_MODE`: Set to "release" for production

## Architecture Overview

### Core Technology Stack
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **Database**: SQLite with GORM ORM
- **Real-time Updates**: WebSockets via gorilla/websocket
- **Frontend**: Server-side rendered Go templates (in `client/` directory)
- **Scheduling**: gocron for periodic tasks

### Directory Structure
```
main.go                 # Application entry point, routing, and scheduled tasks
controllers/            # HTTP handlers for web routes and API endpoints
  - pages.go           # Page rendering controllers
  - podcast.go         # Podcast CRUD operations
  - websockets.go      # WebSocket communication
db/                    # Database layer
  - db.go             # Database initialization
  - dbfunctions.go    # Core database queries
  - podcast.go        # Podcast-specific queries
  - migrations.go     # Database migrations
  - base.go           # Database models (Podcast, PodcastItem, Tag, etc.)
service/               # Business logic layer
  - podcastService.go # Podcast operations and RSS parsing
  - fileService.go    # File download and management
  - itunesService.go  # iTunes API integration for podcast search
  - gpodderService.go # GPodder API support
model/                 # Data structures for API responses and RSS/OPML parsing
client/                # HTML templates (Go templates with custom funcs)
webassets/             # Static assets (CSS, JS, images)
```

### Application Flow

1. **Startup** (main.go):
   - Initializes SQLite database at `$CONFIG/podgrab.db`
   - Runs migrations to ensure schema is up-to-date
   - Sets up Gin router with optional basic authentication
   - Registers routes for both API endpoints and page rendering
   - Starts WebSocket message handler goroutine
   - Initializes cron jobs for background tasks

2. **Background Jobs** (main.go:220-234):
   - `RefreshEpisodes()`: Checks RSS feeds for new episodes
   - `CheckMissingFiles()`: Verifies downloaded files still exist
   - `DownloadMissingImages()`: Downloads missing podcast artwork
   - `UnlockMissedJobs()`: Cleans up stale job locks
   - `UpdateAllFileSizes()`: Updates file size metadata
   - `CreateBackup()`: Creates database backups (every 2 days)
   - All tasks run at intervals based on `CHECK_FREQUENCY` environment variable

3. **Podcast Workflow**:
   - User adds podcast via RSS URL, OPML import, or iTunes search
   - Service layer fetches and parses RSS feed (service/podcastService.go)
   - Episodes are stored in database as PodcastItem records
   - Background jobs automatically download new episodes to `$DATA` directory
   - WebSocket notifications inform frontend of download progress

4. **Database Schema**:
   - `Podcast`: RSS feed metadata (title, URL, author, image, etc.)
   - `PodcastItem`: Individual episodes with download status tracking
   - `Tag`: Labels/groups for organizing podcasts
   - `Setting`: Application configuration
   - `JobLock`: Prevents duplicate background job execution
   - `Migration`: Tracks applied schema migrations

### Key Patterns and Conventions

**Controllers vs Services**:
- Controllers (controllers/) handle HTTP requests/responses and validation
- Services (service/) contain business logic and external API calls
- Database operations are in db/ layer

**Download Status Tracking**:
- Episodes have a `DownloadStatus` enum: `NotDownloaded`, `Downloading`, `Downloaded`, `Deleted`
- File paths are constructed as: `$DATA/<podcast-title>/<episode-filename>`
- Files are detected by checking filesystem before re-downloading

**Template Functions** (main.go:35-128):
- Custom template helpers for formatting dates, file sizes, durations
- `naturalDate`: Converts timestamps to "2 hours ago" format
- `formatFileSize`: Converts bytes to human-readable sizes
- `downloadedEpisodes`/`downloadingEpisodes`: Count episodes by status

**WebSocket Communication**:
- Used for real-time updates to frontend (download progress, new episodes)
- Single WebSocket endpoint at `/ws`
- Messages handled by `controllers.HandleWebsocketMessages()`

**Authentication**:
- Basic HTTP auth enabled when `PASSWORD` env var is set
- Username is always "podgrab"
- Applied to all routes via Gin middleware

## Common Development Patterns

### Adding a New API Endpoint
1. Define query/request structs in controllers/podcast.go
2. Add handler function following pattern: `func HandlerName(c *gin.Context)`
3. Register route in main.go using `router.GET/POST/etc`
4. Implement business logic in service/ layer
5. Add database queries in db/ layer if needed

### Adding a New Background Job
1. Create function in service/ layer
2. Register in `intiCron()` function in main.go
3. Use `gocron.Every(duration).Minutes().Do(functionName)`
4. Consider using JobLock to prevent concurrent execution

### Database Migrations
- Add migration functions in db/migrations.go
- Migrations are tracked in `Migration` table
- Run automatically on startup via `db.Migrate()`

## Testing

This project does not currently have automated tests. When adding new features, manually test:
- Different podcast RSS feed formats
- Download failures and retries
- File system edge cases (missing directories, permissions)
- WebSocket connectivity
- Authentication when enabled
