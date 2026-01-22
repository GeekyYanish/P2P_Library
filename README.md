# 📚 The Knowledge Exchange
### Decentralized P2P Academic Library

The Knowledge Exchange is a full-stack decentralized distributed system designed for peer-to-peer (P2P) academic resource sharing. It uses a **Go (Golang)** backend with microservices architecture and a **React** frontend.

---

## 🚀 Features

- **Decentralized Storage**: Files are shared directly between peers (simulated for MVP).
- **Reputation System**: Earn reputation by uploading useful resources.
- **Fair Access**: Leechers (users who only download) are throttled or restricted.
- **Microservices Architecture**: Separate modules for Gateway, Library, Analytics, and P2P logic.
- **Interactive UI**: Modern dashboard to manage files, reputation, and peers.

---

## 🛠️ Tech Stack

- **Backend**: Go 1.21+ (net/http, goroutines, channels)
- **Frontend**: React + Vite
- **Storage**: Local file system (data/sharedFiles)
- **Networking**: TCP/HTTP for peer communication

---

## 🏃‍♂️ How to Run

### Prerequisites
- Go 1.21 or higher installed
- Node.js and npm installed
- A modern web browser

### Steps

1. **Initialize the Backend Module** (if not already done)
   ```bash
   cd backend
   go mod tidy
   ```

2. **Install Frontend Dependencies**
   ```bash
   cd frontend
   npm install
   ```

3. **Run the Backend Server**
   ```bash
   cd backend
   go run cmd/main.go
   ```
   
   *Optional flags:*
   - `-port=8080` (Default API port)
   - `-name="My Peer Name"` (Custom display name)

4. **Run the Frontend Development Server**
   ```bash
   cd frontend
   npm run dev
   ```

5. **Access the Application**
   Open your browser and navigate to:
   [http://localhost:5173](http://localhost:5173)

---

## 📂 Project Structure

```
knowledge-exchange/
├── backend/
│   ├── cmd/            # Entry point (main.go)
│   ├── models/         # Data structures (Student, File, Rating)
│   ├── gateway/        # HTTP Server & Router
│   ├── library/        # Indexing & Transfer logic
│   ├── analytics/      # Reputation & Throttling engine
│   ├── auth/           # Authentication logic
│   ├── storage/        # Data persistence
│   └── utils/          # Hashing & Network helpers
├── frontend/
│   ├── src/
│   │   ├── components/ # Reusable UI components
│   │   ├── pages/      # Page components
│   │   ├── context/    # React contexts
│   │   └── services/   # API service layer
│   ├── public/         # Static assets
│   └── index.html      # Main HTML template
└── data/               # Local storage
```

---

## 🧪 Testing

1. Open http://localhost:5173 in your browser.
2. Register with a username and password.
3. Login to access the main dashboard.
4. Go to the "Upload" tab and share a file.
5. See your reputation increase!
6. Open an Incognito window to simulate a second peer.

---

## 📝 License
MIT License
