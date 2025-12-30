# GitHub Stats Backend

A lightweight Go backend service that fetches GitHub user statistics using the GitHub GraphQL API. Includes a Next.js frontend dashboard.

## Features

- **User Stats** — Fetch comprehensive profile data including repositories, contributions, and activity metrics
- **Pinned Repositories** — Retrieve a user's pinned/showcase repositories
- **Private Contributions** — Access private contribution counts (for authenticated user)
- **Language Statistics** — Get language breakdown across repositories

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Health check |
| `GET` | `/api/stats/username` | Get user profile and contribution stats |
| `GET` | `/api/pinned/username` | Get user's pinned repositories |

## Setup

### Prerequisites

- Go 1.21+
- GitHub Personal Access Token with `read:user` and `repo` scopes

### Installation

1. Clone the repository
   ```bash
   git clone https://github.com/Mystery-Coder/github-stats-backend.git
   cd github-stats-backend
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Create a `.env` file
   ```env
   GITHUB_TOKEN=your_github_token_here
   PORT=8080
   ```

4. Run the server
   ```bash
   go run main.go
   ```

## Frontend

The frontend is a Next.js dashboard located in the `frontend/` directory.

```bash
cd frontend
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) to view the dashboard.

## Usage

You can access the API via browser or command line:

**Browser:**
```
http://localhost:8080/api/stats/Mystery-Coder
http://localhost:8080/api/pinned/Mystery-Coder
```

**cURL:**
```bash
# Get user stats
curl http://localhost:8080/api/stats/Mystery-Coder

# Get pinned repos
curl http://localhost:8080/api/pinned/Mystery-Coder
```