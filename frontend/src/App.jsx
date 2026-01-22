import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import { ToastProvider } from './context/ToastContext'
import Login from './pages/Login'
import Signup from './pages/Signup'
import Home from './pages/Home'
import Library from './pages/Library'
import Upload from './pages/Upload'
import Reputation from './pages/Reputation'
import Peers from './pages/Peers'
import ProtectedRoute from './components/ProtectedRoute'
import Header from './components/Header'
import Footer from './components/Footer'

function App() {
  return (
    <AuthProvider>
      <ToastProvider>
        <div className="app">
          <Routes>
            {/* Public routes */}
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<Signup />} />
            
            {/* Protected routes */}
            <Route
              path="/*"
              element={
                <ProtectedRoute>
                  <Header />
                  <main className="main">
                    <Routes>
                      <Route path="/home" element={<Home />} />
                      <Route path="/library" element={<Library />} />
                      <Route path="/upload" element={<Upload />} />
                      <Route path="/reputation" element={<Reputation />} />
                      <Route path="/peers" element={<Peers />} />
                      <Route path="/" element={<Navigate to="/home" replace />} />
                    </Routes>
                  </main>
                  <Footer />
                </ProtectedRoute>
              }
            />
          </Routes>
        </div>
      </ToastProvider>
    </AuthProvider>
  )
}

export default App
