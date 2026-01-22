# 📚 The Knowledge Exchange
### Decentralized P2P Academic Library

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18.x-61DAFB?logo=react)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)

The Knowledge Exchange is a full-stack **decentralized peer-to-peer (P2P) academic resource sharing platform** designed to facilitate knowledge distribution among students and researchers. Built with a **Go (Golang)** microservices backend and a modern **React + Vite** frontend, it implements a reputation-based system to encourage fair resource sharing and discourage leeching behavior.

---

## 📖 Table of Contents

- [Overview](#-overview)
- [Key Features](#-key-features)
- [Architecture](#-architecture)
- [Tech Stack](#️-tech-stack)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [API Documentation](#-api-documentation)
- [Features in Detail](#-features-in-detail)
- [Configuration](#️-configuration)
- [Testing](#-testing)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [Roadmap](#-roadmap)
- [License](#-license)

---

## 🌟 Overview

Traditional academic resource sharing often relies on centralized platforms that can be restrictive, expensive, or subject to censorship. The Knowledge Exchange addresses these challenges by creating a **peer-to-peer ecosystem** where:

- **Resources are shared directly** between users without central authority
- **Contributors are rewarded** with reputation points for sharing quality content
- **Free riders are discouraged** through intelligent throttling mechanisms
- **Privacy is maintained** through decentralized architecture
- **Accessibility is enhanced** through modern, intuitive UI/UX

### Use Cases

- **University Students**: Share textbooks, research papers, and study materials
- **Research Communities**: Distribute academic papers and datasets
- **Study Groups**: Collaborate on course materials and resources
- **Open Education**: Support free access to educational content

---

## 🚀 Key Features

### Core Functionality

- ✅ **Peer-to-Peer File Sharing**: Direct file transfers between peers with no central storage
- ✅ **User Authentication**: Secure registration and login system with JWT tokens
- ✅ **File Upload & Download**: Seamless file sharing with progress tracking
- ✅ **Search & Discovery**: Find resources across the network by title, subject, or author
- ✅ **Reputation System**: Earn points for uploading; lose points for only downloading
- ✅ **Fair Access Control**: Throttle or restrict users with negative reputation
- ✅ **Resource Ratings**: Community-driven quality assessment (1-5 stars)
- ✅ **Real-time Notifications**: Toast notifications for user actions and events

### Advanced Features

- 🔒 **Secure Authentication**: Password hashing with bcrypt, JWT-based sessions
- 📊 **Analytics Dashboard**: Track uploads, downloads, and reputation over time
- 🌐 **RESTful API**: Clean, well-documented API for all operations
- 🎨 **Modern UI**: Responsive React interface with dark mode aesthetics
- 🔍 **Advanced Search**: Filter and sort resources by multiple criteria
- 📈 **Reputation Leaderboard**: View top contributors in the community

---

## 🏗 Architecture

The Knowledge Exchange follows a **microservices architecture** with clear separation of concerns:

### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend (React)                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │   Home   │  │  Upload  │  │  Library │  │   Login  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│         │              │              │              │       │
│         └──────────────┴──────────────┴──────────────┘       │
│                          │ (API Service)                     │
└──────────────────────────┼───────────────────────────────────┘
                           │ HTTP/REST (Port 5173 → 8080)
┌──────────────────────────┼───────────────────────────────────┐
│                          ▼                                    │
│              ┌────────────────────┐                          │
│              │   Gateway Service  │  (HTTP Router)           │
│              └────────────────────┘                          │
│                │         │         │                          │
│       ┌────────┼─────────┼─────────┼────────┐                │
│       ▼        ▼         ▼         ▼        ▼                │
│  ┌────────┐ ┌──────┐ ┌────────┐ ┌──────┐ ┌──────┐          │
│  │  Auth  │ │Library│ │Analytics│ │Storage│ │P2P  │          │
│  │Service │ │Service│ │ Service │ │Service│ │Layer│          │
│  └────────┘ └──────┘ └────────┘ └──────┘ └──────┘          │
│                                                               │
│          Backend (Go) - Microservices Architecture           │
└───────────────────────────────────────────────────────────────┘
                           │
                           ▼
             ┌──────────────────────────┐
             │   Local File System      │
             │   (data/sharedFiles)     │
             └──────────────────────────┘
```

### Backend Microservices

| Service | Responsibility | Key Components |
|---------|---------------|----------------|
| **Gateway** | HTTP routing, request handling, CORS | `router.go`, `handlers.go`, `auth_handlers.go` |
| **Auth** | User authentication, JWT tokens, sessions | `auth.go`, `middleware.go` |
| **Library** | File indexing, search, upload/download logic | `indexer.go`, `transfer.go` |
| **Analytics** | Reputation calculation, throttling engine | `reputation.go`, `throttle.go` |
| **Storage** | Data persistence for users, files, ratings | `user_store.go`, `file_store.go` |
| **P2P** | Peer discovery, network communication | `discovery.go`, `network.go` |

### Data Flow

1. **User Authentication**: Frontend → Gateway → Auth Service → Storage
2. **File Upload**: Frontend → Gateway → Library Service → File System + Storage
3. **File Download**: Frontend → Gateway → Library Service → File System
4. **Reputation Update**: Library Service → Analytics Service → Storage

---

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.21+
- **HTTP Server**: Native `net/http` with custom routing
- **Concurrency**: Goroutines and channels for async operations
- **Authentication**: JWT tokens, bcrypt password hashing
- **Storage**: JSON-based file storage (upgradeable to SQL/NoSQL)

### Frontend
- **Framework**: React 18.x
- **Build Tool**: Vite 5.x
- **State Management**: React Context API (AuthContext, ToastContext)
- **Routing**: React Router DOM v6
- **HTTP Client**: Axios
- **Styling**: Modern CSS with CSS variables and animations

### Development Tools
- **Version Control**: Git & GitHub
- **Package Management**: Go modules, npm
- **API Testing**: Can use Postman, curl, or browser DevTools

---

## 📂 Project Structure

```
knowledge-exchange/
├── backend/
│   ├── cmd/
│   │   └── main.go                 # Entry point, server initialization
│   ├── models/
│   │   ├── student.go              # User/Student data model
│   │   ├── file.go                 # Shared file metadata
│   │   └── rating.go               # File rating model
│   ├── gateway/
│   │   ├── router.go               # HTTP route definitions
│   │   ├── handlers.go             # General API handlers
│   │   ├── auth_handlers.go        # Auth endpoints (login, register)
│   │   └── discovery.go            # Peer discovery logic
│   ├── library/
│   │   ├── indexer.go              # File indexing and search
│   │   ├── transfer.go             # Upload/download logic
│   │   └── catalog.go              # File catalog management
│   ├── analytics/
│   │   ├── reputation.go           # Reputation calculation engine
│   │   └── throttle.go             # Download throttling logic
│   ├── auth/
│   │   ├── auth.go                 # JWT generation, password hashing
│   │   └── middleware.go           # Auth middleware for protected routes
│   ├── storage/
│   │   ├── user_store.go           # User data persistence
│   │   ├── file_store.go           # File metadata persistence
│   │   └── rating_store.go         # Rating data persistence
│   ├── utils/
│   │   ├── hash.go                 # File hashing utilities
│   │   └── network.go              # Network helper functions
│   ├── go.mod                      # Go module dependencies
│   └── go.sum                      # Dependency checksums
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Header.jsx          # Navigation header
│   │   │   ├── Footer.jsx          # Footer component
│   │   │   ├── FileCard.jsx        # File display card
│   │   │   ├── Toast.jsx           # Toast notification
│   │   │   ├── SearchBar.jsx       # Search input component
│   │   │   └── ProtectedRoute.jsx  # Route protection wrapper
│   │   ├── pages/
│   │   │   ├── Home.jsx            # Dashboard/home page
│   │   │   ├── Login.jsx           # Login page
│   │   │   ├── Register.jsx        # Registration page
│   │   │   ├── Upload.jsx          # File upload page
│   │   │   ├── Library.jsx         # Browse files page
│   │   │   └── Reputation.jsx      # User stats/leaderboard
│   │   ├── context/
│   │   │   ├── AuthContext.jsx     # Authentication state
│   │   │   └── ToastContext.jsx    # Toast notification state
│   │   ├── services/
│   │   │   └── api.js              # Axios API client
│   │   ├── App.jsx                 # Main app component with routing
│   │   ├── main.jsx                # React entry point
│   │   └── App.css                 # Global styles
│   ├── public/
│   │   └── vite.svg                # Favicon
│   ├── index.html                  # HTML template
│   ├── vite.config.js              # Vite configuration (proxy setup)
│   ├── package.json                # npm dependencies
│   └── package-lock.json           # Dependency lock file
│
├── data/
│   └── sharedFiles/                # Uploaded files storage
│
├── .gitignore                      # Git ignore rules
└── README.md                       # This file
```

---

## 🚀 Getting Started

### Prerequisites

Ensure you have the following installed on your system:

- **Go**: Version 1.21 or higher ([Download](https://go.dev/dl/))
- **Node.js**: Version 18.x or higher ([Download](https://nodejs.org/))
- **npm**: Comes with Node.js (or use yarn/pnpm)
- **Git**: For version control ([Download](https://git-scm.com/))
- **Modern Browser**: Chrome, Firefox, Edge, or Safari

### Installation Steps

#### 1. Clone the Repository

```bash
git clone https://github.com/GeekyYanish/P2P_Library.git
cd P2P_Library
```

#### 2. Backend Setup

Navigate to the backend directory and initialize Go modules:

```bash
cd backend
go mod tidy
```

This will download all required Go dependencies.

#### 3. Frontend Setup

Navigate to the frontend directory and install npm packages:

```bash
cd ../frontend
npm install
```

This will install React, Vite, React Router, Axios, and other dependencies.

#### 4. Create Data Directory

Ensure the data directory exists for file storage:

```bash
# From project root
mkdir -p data/sharedFiles
```

### Running the Application

You need to run both backend and frontend servers:

#### Terminal 1: Start Backend Server

```bash
cd backend
go run cmd/main.go
```

**Optional flags:**
- `-port=8080` - Specify API port (default: 8080)
- `-name="My Peer"` - Set peer display name

Expected output:
```
🚀 Starting Knowledge Exchange API Server...
📡 Server listening on :8080
🔗 API Base URL: http://localhost:8080
```

#### Terminal 2: Start Frontend Dev Server

```bash
cd frontend
npm run dev
```

Expected output:
```
VITE v5.x.x  ready in 500 ms

➜  Local:   http://localhost:5173/
➜  Network: use --host to expose
```

#### 5. Access the Application

Open your browser and navigate to:

🌐 **http://localhost:5173**

You should see the Knowledge Exchange login/register page.

---

## 📡 API Documentation

### Base URL
```
http://localhost:8080
```

### Authentication Endpoints

#### Register New User
```http
POST /api/register
Content-Type: application/json

{
  "username": "john_doe",
  "password": "secure_password123",
  "name": "John Doe"
}

Response (201):
{
  "message": "User registered successfully",
  "user": {
    "id": "user-uuid",
    "username": "john_doe",
    "name": "John Doe"
  }
}
```

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "secure_password123"
}

Response (200):
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "user-uuid",
    "username": "john_doe",
    "name": "John Doe",
    "reputation": 100
  }
}
```

### File Management Endpoints

#### Upload File
```http
POST /api/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

FormData:
  - file: <binary>
  - title: "Introduction to Algorithms"
  - subject: "Computer Science"
  - description: "Classical algorithms textbook"

Response (200):
{
  "message": "File uploaded successfully",
  "file": {
    "id": "file-uuid",
    "title": "Introduction to Algorithms",
    "filename": "intro-algorithms.pdf",
    "uploader": "john_doe",
    "uploadTime": "2026-01-22T19:15:00Z"
  }
}
```

#### Get All Files
```http
GET /api/files
Authorization: Bearer <token>

Response (200):
{
  "files": [
    {
      "id": "file-uuid",
      "title": "Introduction to Algorithms",
      "filename": "intro-algorithms.pdf",
      "subject": "Computer Science",
      "uploader": "john_doe",
      "uploadTime": "2026-01-22T19:15:00Z",
      "downloads": 5,
      "rating": 4.5
    }
  ]
}
```

#### Download File
```http
GET /api/download/:fileId
Authorization: Bearer <token>

Response (200):
Content-Type: application/octet-stream
Content-Disposition: attachment; filename="intro-algorithms.pdf"

<file binary data>
```

#### Search Files
```http
GET /api/search?q=algorithms&subject=Computer%20Science
Authorization: Bearer <token>

Response (200):
{
  "results": [
    {
      "id": "file-uuid",
      "title": "Introduction to Algorithms",
      "subject": "Computer Science",
      "relevance": 0.95
    }
  ]
}
```

### User & Analytics Endpoints

#### Get User Profile
```http
GET /api/user/profile
Authorization: Bearer <token>

Response (200):
{
  "user": {
    "id": "user-uuid",
    "username": "john_doe",
    "name": "John Doe",
    "reputation": 150,
    "uploadsCount": 10,
    "downloadsCount": 25,
    "joinedAt": "2026-01-15T10:00:00Z"
  }
}
```

#### Get Reputation Leaderboard
```http
GET /api/analytics/leaderboard?limit=10
Authorization: Bearer <token>

Response (200):
{
  "leaderboard": [
    {
      "rank": 1,
      "username": "top_contributor",
      "reputation": 500,
      "uploads": 50
    }
  ]
}
```

#### Rate a File
```http
POST /api/rate
Authorization: Bearer <token>
Content-Type: application/json

{
  "fileId": "file-uuid",
  "rating": 5,
  "comment": "Excellent resource!"
}

Response (200):
{
  "message": "Rating submitted successfully",
  "averageRating": 4.7
}
```

### Error Responses

All endpoints may return error responses:

```json
{
  "error": "Error message description",
  "code": "ERROR_CODE"
}
```

Common HTTP status codes:
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing/invalid token)
- `403` - Forbidden (insufficient reputation)
- `404` - Not Found (resource doesn't exist)
- `500` - Internal Server Error

---

## 🎯 Features in Detail

### 1. Reputation System

The reputation system is the core mechanism for promoting fair sharing:

**Earning Reputation:**
- Upload a new file: `+10 points`
- File gets downloaded: `+2 points per download`
- File gets 5-star rating: `+5 points`

**Losing Reputation:**
- Download a file: `-1 point`
- Receive 1-star rating on your upload: `-3 points`

**Reputation Effects:**
- **100+ points**: Full access, no throttling
- **50-99 points**: Normal access
- **0-49 points**: Download speed throttled
- **Negative points**: Download blocked until contribution

### 2. Fair Access Control

Prevent leeching through intelligent throttling:

- **Download/Upload Ratio Tracking**: System monitors user behavior
- **Adaptive Throttling**: Low contributors get slower downloads
- **Grace Period**: New users get initial reputation to get started
- **Appeals**: Users can regain access by uploading quality content

### 3. Search & Discovery

Advanced search capabilities:

- **Full-text search**: Search in titles, descriptions, subjects
- **Filter by subject**: Computer Science, Mathematics, Physics, etc.
- **Sort options**: Relevance, date, rating, downloads
- **Tag-based discovery**: Browse by topic tags

### 4. Security Features

- **Password Security**: Bcrypt hashing with salt
- **JWT Authentication**: Stateless token-based auth
- **Session Management**: Token expiration and renewal
- **CORS Protection**: Configured for frontend origin
- **Input Validation**: Server-side validation for all inputs

---

## ⚙️ Configuration

### Backend Configuration

Edit `backend/cmd/main.go` to configure:

```go
const (
    DEFAULT_PORT = "8080"
    JWT_SECRET = "your-secret-key-here" // Change in production!
    TOKEN_EXPIRY = 24 * time.Hour
    MAX_FILE_SIZE = 100 * 1024 * 1024 // 100 MB
)
```

### Frontend Configuration

Edit `frontend/vite.config.js` for API proxy:

```javascript
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

### Environment Variables (Optional)

Create `.env` files for configuration:

**Backend (.env):**
```env
PORT=8080
JWT_SECRET=your-secret-key
DATA_PATH=./data
```

**Frontend (.env):**
```env
VITE_API_URL=http://localhost:8080
```

---

## 🧪 Testing

### Manual Testing

1. **Register a new user**
   - Navigate to registration page
   - Create account with username/password
   - Verify successful registration

2. **Login**
   - Use credentials to log in
   - Verify JWT token is stored
   - Check redirect to dashboard

3. **Upload a file**
   - Go to Upload page
   - Fill in file details
   - Upload a test file
   - Verify reputation increases

4. **Browse library**
   - Navigate to Library page
   - Search for uploaded file
   - Verify file appears in results

5. **Download a file**
   - Click download on a file
   - Verify file downloads
   - Check reputation decreased

6. **Test with multiple users**
   - Open incognito window
   - Register second user
   - Test peer interaction

### API Testing with curl

**Register:**
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123","name":"Test User"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123"}'
```

**Get Files (with token):**
```bash
curl http://localhost:8080/api/files \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 🔧 Troubleshooting

### Common Issues

#### Port Already in Use

**Problem:** `Error: listen tcp :8080: bind: address already in use`

**Solution:**
```bash
# Find process using port 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill the process or use different port
go run cmd/main.go -port=8081
```

#### CORS Errors

**Problem:** `Access-Control-Allow-Origin error in browser console`

**Solution:** Ensure backend CORS headers are set correctly in `gateway/router.go`:
```go
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
w.Header().Set("Access-Control-Allow-Credentials", "true")
```

#### Frontend Can't Reach Backend

**Problem:** `Network Error` or `ERR_CONNECTION_REFUSED`

**Solution:**
1. Verify backend is running: `curl http://localhost:8080/api/files`
2. Check proxy in `vite.config.js`
3. Verify API base URL in `frontend/src/services/api.js`

#### File Upload Fails

**Problem:** `413 Payload Too Large` or upload timeout

**Solution:**
1. Check `MAX_FILE_SIZE` in backend config
2. Ensure `data/sharedFiles` directory exists and is writable
3. Verify disk space available

#### JWT Token Expired

**Problem:** `401 Unauthorized` after some time

**Solution:**
- Token expires after 24 hours by default
- Log out and log back in to get new token
- Implement token refresh mechanism (future enhancement)

---

## 🤝 Contributing

We welcome contributions to the Knowledge Exchange! Here's how you can help:

### Ways to Contribute

1. **Report Bugs**: Open an issue with detailed reproduction steps
2. **Suggest Features**: Share your ideas in GitHub Issues
3. **Submit Pull Requests**: Fix bugs or add features
4. **Improve Documentation**: Help make docs clearer
5. **Share Feedback**: User experience improvements

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Test thoroughly
5. Commit: `git commit -m 'Add amazing feature'`
6. Push: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Code Style

- **Go**: Follow `gofmt` standards, use `golint`
- **JavaScript/React**: Use ESLint, Prettier for formatting
- **Commits**: Use conventional commits format

---

## 🗺 Roadmap

### Version 1.0 (MVP) ✅
- [x] User authentication system
- [x] File upload/download
- [x] Basic reputation system
- [x] Search functionality
- [x] React frontend

### Version 1.5 (Planned)
- [ ] Real P2P networking with libp2p
- [ ] Distributed hash table (DHT) for file discovery
- [ ] WebRTC for direct peer connections
- [ ] File encryption for privacy
- [ ] Mobile-responsive UI improvements

### Version 2.0 (Future)
- [ ] Mobile apps (iOS/Android)
- [ ] Blockchain-based reputation (immutable)
- [ ] Smart contracts for resource trading
- [ ] Video/audio streaming support
- [ ] Real-time chat between peers
- [ ] Decentralized identity (DID)

---

## 📄 License

This project is licensed under the **MIT License** - see below for details:

```
MIT License

Copyright (c) 2026 The Knowledge Exchange Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 📞 Contact & Support

- **GitHub**: [https://github.com/GeekyYanish/P2P_Library](https://github.com/GeekyYanish/P2P_Library)
- **Issues**: [Report a bug or request a feature](https://github.com/GeekyYanish/P2P_Library/issues)

---

## 🙏 Acknowledgments

- Inspired by BitTorrent, IPFS, and academic resource sharing needs
- Built with amazing open-source technologies
- Thanks to all contributors and the open-source community

---

<div align="center">

**⭐ If you find this project useful, please consider giving it a star! ⭐**

Made with ❤️ by the Knowledge Exchange Team

</div>
