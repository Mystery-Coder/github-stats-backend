# GitHub Analytics Dashboard

A high-performance, professional analytics dashboard that visualizes GitHub user data with a premium aesthetic. Built with Next.js 15 and Go.

## Tech Stack

### Frontend
- **Framework**: Next.js 15 (App Router)
- **Styling**: Tailwind CSS v4
- **Animations**: Framer Motion
- **Charts**: Recharts
- **Icons**: Lucide React

### Backend
- **Language**: Go (Golang)
- **Framework**: Gin Web Framework
- **API**: GitHub GraphQL API

## Getting Started

### Prerequisites
- Node.js 18+
- Go 1.21+
- GitHub Personal Access Token (PAT)

### 1. Backend Setup
Navigate to the root directory and start the Go server:

```bash
# Set your GitHub Token
$env:GITHUB_TOKEN="your_token_here"

# Run the server
go run main.go
```
The backend will start on `http://localhost:8080`.

### 2. Frontend Setup
Navigate to the frontend directory:

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```
Open [http://localhost:3000](http://localhost:3000) to view the dashboard.

## Production Build

To build the application for production:

**Frontend:**
```bash
npm run build
npm start
```