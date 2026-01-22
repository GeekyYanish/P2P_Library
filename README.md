# 📚 The Knowledge Exchange
### Decentralized P2P Academic Library

The Knowledge Exchange is a full-stack decentralized distributed system designed for peer-to-peer (P2P) academic resource sharing. It uses a **Go (Golang)** backend with microservices architecture and a **vanilla JavaScript** frontend.

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
- **Frontend**: HTML5, CSS3, JavaScript (ES6+)
- **Storage**: Local file system (data/sharedFiles)
- **Networking**: TCP/HTTP for peer communication

---

## 🏃‍♂️ How to Run

### Prerequisites
- Go 1.21 or higher installed
- A modern web browser

### Steps

1. **Initialize the Module** (if not already done)
   ```bash
   cd backend
   go mod tidy
   ```

2. **Run the Backend Server**
   ```bash
   cd backend
   go run cmd/main.go
   ```
   
   *Optional flags:*
   - `-port=8080` (Default API port)
   - `-name="My Peer Name"` (Custom display name)

3. **Access the Application**
   Open your browser and navigate to:
   [http://localhost:3000](http://localhost:3000) (or whatever port you configured)

   *Note: The backend serves the frontend files automatically.*

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
│   └── utils/          # Hashing & Network helpers
├── frontend/
│   ├── css/            # Styles
│   ├── js/             # Application Logic
│   └── index.html      # Main UI
└── data/               # Local storage
```

---

## 🧪 Testing

1. Open http://localhost:3000 in your browser.
2. Enter a username to register as a peer.
3. Go to the "Upload" tab and share a file (creates a dummy file in `data/sharedFiles`).
4. See your reputation increase!
5. Open an Incognito window to simulate a second peer.

---

## 📝 License
MIT License
