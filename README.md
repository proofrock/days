# Days - Journal Application

A simple, portable web-based journaling application for tracking daily entries with various fields like mood, sleep quality, meals, and location.

## Features

- **Calendar View**: Visual calendar interface showing days with entries
- **Rich Entry Fields**: Track multiple aspects of your day including:
  - Location (with geolocation support and Google Maps integration)
  - Day rating (1-10 stars)
  - General overview
  - Mood rating with description
  - Lunch and dinner notes
  - TV shows watched
  - Sleep quality rating with description
- **Responsive Design**: Works on desktop and mobile devices
- **Single User**: No authentication required
- **Portable**: Single binary with embedded frontend
- **Lightweight**: SQLite database for data storage

## Tech Stack

- **Backend**: Go with SQLite
- **Frontend**: Svelte 5, TypeScript, Bootstrap
- **Build**: Makefile, Docker multi-stage builds

## Prerequisites

- Go 1.21 or higher
- Node.js 20 or higher
- Make (for build automation)
- Docker (optional, for containerized deployment)

## Building

### Using Make

```bash
# Install dependencies
make deps

# Build everything (frontend + backend)
make build

# Build for Linux AMD64
make build VERSION=1.0.0 build-linux-amd64
```

### Manual Build

```bash
# Build frontend
cd frontend
npm install
npm run build
cd ..

# Build backend
go build -ldflags "-X main.Version=1.0.0" -o bin/days .
```

### Docker Build

```bash
# Build Docker image
make docker-build VERSION=1.0.0

# Or directly with Docker
docker build -t days-journal:1.0.0 --build-arg VERSION=1.0.0 .
```

## Running

### Local Binary

```bash
# Using Make
make run

# Or run directly
./bin/days

# Custom configuration
DB_PATH=./mydata/journal.db PORT=3000 ./bin/days
```

### Docker

```bash
# Using Make
make docker-run

# Or directly with Docker
docker run -p 8080:8080 -v $(pwd)/data:/app/data days-journal:1.0.0
```

The application will be available at `http://localhost:8080`

## Configuration

Environment variables:

- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database file path (default: ./data/journal.db)

## Development

```bash
# Run in development mode (no build)
make dev

# Run frontend dev server separately
cd frontend
npm run dev
```

## Database Schema

See `scripts/schema.sql` for the complete database schema.

The database uses a flexible two-table structure:
- `entries`: Header table with date and timestamp
- `details`: Key-value pairs for all entry fields

This design makes it easy to add new fields in the future.

## API Endpoints

- `GET /api/entries/{date}` - Get entry for a specific date
- `POST /api/entries/{date}` - Create or update entry
- `DELETE /api/entries/{date}` - Delete entry
- `GET /api/entries/month/{year}/{month}` - Get all entry dates for a month

## Versioning

This project follows semantic versioning (MAJOR.MINOR.PATCH).

Update the `VERSION` file and rebuild to change the version:

```bash
echo "1.1.0" > VERSION
make build VERSION=$(cat VERSION)
```

## License

MIT
