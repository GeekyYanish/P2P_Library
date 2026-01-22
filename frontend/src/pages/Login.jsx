import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useToast } from '../context/ToastContext';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuth();
  const { showToast } = useToast();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    const result = await login(email, password);
    
    if (result.success) {
      showToast('Welcome back! Login successful.', 'success');
      navigate('/home');
    } else {
      showToast(result.error || 'Login failed', 'error');
    }
    
    setLoading(false);
  };

  return (
    <div className="auth-page">
      {/* Animated background orbs */}
      <div className="auth-bg-orbs">
        <div className="orb orb-1"></div>
        <div className="orb orb-2"></div>
        <div className="orb orb-3"></div>
      </div>
      
      <div className="auth-container">
        <div className="auth-card">
          <div className="auth-header">
            <div className="auth-logo">
              <span className="logo-icon">📚</span>
              <div className="logo-glow"></div>
            </div>
            <h2>Welcome Back</h2>
            <p>Sign in to The Knowledge Exchange</p>
          </div>

          <form onSubmit={handleSubmit} className="auth-form">
            <div className="form-group">
              <label htmlFor="email">Email Address</label>
              <input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                placeholder="you@example.com"
                autoComplete="email"
              />
            </div>

            <div className="form-group">
              <label htmlFor="password">Password</label>
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                placeholder="••••••••"
                autoComplete="current-password"
              />
            </div>

            <button 
              type="submit" 
              className="btn btn-primary btn-block btn-lg" 
              disabled={loading}
            >
              {loading ? (
                <>
                  <span className="loading-spinner" style={{width: '20px', height: '20px', marginBottom: 0}}></span>
                  Signing in...
                </>
              ) : (
                <>
                  🚀 Sign In
                </>
              )}
            </button>
          </form>

          <div className="auth-footer">
            <p>Don't have an account? <Link to="/signup">Create one</Link></p>
          </div>

          <div className="auth-demo">
            <p><strong>🔐 Demo Credentials</strong></p>
            <p style={{marginTop: '0.5rem'}}>
              <code>admin@knowledge-exchange.com</code> / <code>admin123</code>
            </p>
          </div>
        </div>

        <div className="auth-features">
          <div className="auth-feature">
            <span>🌐</span>
            <span>Decentralized</span>
          </div>
          <div className="auth-feature">
            <span>🔒</span>
            <span>Secure</span>
          </div>
          <div className="auth-feature">
            <span>⚡</span>
            <span>Fast</span>
          </div>
        </div>
      </div>

      <style>{`
        .auth-bg-orbs {
          position: fixed;
          inset: 0;
          overflow: hidden;
          pointer-events: none;
        }
        
        .orb {
          position: absolute;
          border-radius: 50%;
          filter: blur(80px);
          opacity: 0.5;
          animation: float-orb 20s ease-in-out infinite;
        }
        
        .orb-1 {
          width: 600px;
          height: 600px;
          background: radial-gradient(circle, rgba(139, 92, 246, 0.4) 0%, transparent 70%);
          top: -200px;
          left: -200px;
        }
        
        .orb-2 {
          width: 500px;
          height: 500px;
          background: radial-gradient(circle, rgba(6, 182, 212, 0.3) 0%, transparent 70%);
          bottom: -150px;
          right: -150px;
          animation-delay: -7s;
        }
        
        .orb-3 {
          width: 400px;
          height: 400px;
          background: radial-gradient(circle, rgba(244, 114, 182, 0.25) 0%, transparent 70%);
          top: 50%;
          right: 20%;
          animation-delay: -14s;
        }
        
        @keyframes float-orb {
          0%, 100% { transform: translate(0, 0) scale(1); }
          25% { transform: translate(50px, -30px) scale(1.05); }
          50% { transform: translate(-20px, 50px) scale(0.95); }
          75% { transform: translate(30px, 20px) scale(1.02); }
        }
        
        .auth-logo {
          position: relative;
          display: inline-block;
        }
        
        .auth-logo .logo-icon {
          font-size: 5rem;
          display: block;
          animation: bounce-gentle 3s ease-in-out infinite;
        }
        
        .logo-glow {
          position: absolute;
          inset: -20px;
          background: radial-gradient(circle, rgba(139, 92, 246, 0.4) 0%, transparent 70%);
          filter: blur(20px);
          z-index: -1;
          animation: pulse-glow 2s ease-in-out infinite;
        }
        
        @keyframes bounce-gentle {
          0%, 100% { transform: translateY(0); }
          50% { transform: translateY(-10px); }
        }
        
        @keyframes pulse-glow {
          0%, 100% { opacity: 0.5; transform: scale(1); }
          50% { opacity: 0.8; transform: scale(1.1); }
        }
        
        .auth-features {
          display: flex;
          justify-content: center;
          gap: 2rem;
          margin-top: 2rem;
        }
        
        .auth-feature {
          display: flex;
          align-items: center;
          gap: 0.5rem;
          color: var(--text-muted);
          font-size: 0.875rem;
          padding: 0.5rem 1rem;
          background: rgba(255, 255, 255, 0.03);
          border-radius: 9999px;
          border: 1px solid rgba(255, 255, 255, 0.05);
        }
        
        .auth-demo code {
          background: rgba(139, 92, 246, 0.2);
          padding: 0.15rem 0.4rem;
          border-radius: 4px;
          font-size: 0.8rem;
          color: var(--primary-light);
        }
      `}</style>
    </div>
  );
};

export default Login;
