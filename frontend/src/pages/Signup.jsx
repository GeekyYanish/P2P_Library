import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../services/api';
import { useToast } from '../context/ToastContext';

const Signup = () => {
  const [formData, setFormData] = useState({
    email: '',
    username: '',
    password: '',
    confirmPassword: ''
  });
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { showToast } = useToast();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const getPasswordStrength = () => {
    const { password } = formData;
    if (!password) return { level: 0, text: '', color: '' };
    if (password.length < 6) return { level: 1, text: 'Too short', color: '#ef4444' };
    if (password.length < 8) return { level: 2, text: 'Weak', color: '#f59e0b' };
    if (password.length >= 8 && /[A-Z]/.test(password) && /[0-9]/.test(password)) {
      return { level: 4, text: 'Strong', color: '#10b981' };
    }
    return { level: 3, text: 'Good', color: '#06b6d4' };
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (formData.password !== formData.confirmPassword) {
      showToast('Passwords do not match', 'error');
      return;
    }

    if (formData.password.length < 6) {
      showToast('Password must be at least 6 characters', 'error');
      return;
    }

    setLoading(true);

    try {
      const response = await api.register({
        email: formData.email,
        username: formData.username,
        password: formData.password
      });
      
      if (response.data.success) {
        showToast('Account created! Please login.', 'success');
        navigate('/login');
      } else {
        showToast(response.data.error || 'Registration failed', 'error');
      }
    } catch (error) {
      showToast(error.response?.data?.error || 'Registration failed', 'error');
    }
    
    setLoading(false);
  };

  const strength = getPasswordStrength();

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
              <span className="logo-icon">🎓</span>
              <div className="logo-glow"></div>
            </div>
            <h2>Join the Network</h2>
            <p>Create your Knowledge Exchange account</p>
          </div>

          <form onSubmit={handleSubmit} className="auth-form">
            <div className="form-group">
              <label htmlFor="email">Email Address</label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                required
                placeholder="you@example.com"
                autoComplete="email"
              />
            </div>

            <div className="form-group">
              <label htmlFor="username">Username</label>
              <input
                type="text"
                id="username"
                name="username"
                value={formData.username}
                onChange={handleChange}
                required
                placeholder="Choose a username"
                autoComplete="username"
              />
            </div>

            <div className="form-group">
              <label htmlFor="password">Password</label>
              <input
                type="password"
                id="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
                required
                placeholder="Create a password"
                autoComplete="new-password"
              />
              {formData.password && (
                <div className="password-strength">
                  <div className="strength-bars">
                    {[1, 2, 3, 4].map(i => (
                      <div 
                        key={i} 
                        className={`strength-bar ${strength.level >= i ? 'active' : ''}`}
                        style={{ backgroundColor: strength.level >= i ? strength.color : undefined }}
                      />
                    ))}
                  </div>
                  <span className="strength-text" style={{ color: strength.color }}>
                    {strength.text}
                  </span>
                </div>
              )}
            </div>

            <div className="form-group">
              <label htmlFor="confirmPassword">Confirm Password</label>
              <input
                type="password"
                id="confirmPassword"
                name="confirmPassword"
                value={formData.confirmPassword}
                onChange={handleChange}
                required
                placeholder="Confirm your password"
                autoComplete="new-password"
              />
              {formData.confirmPassword && formData.password !== formData.confirmPassword && (
                <small style={{ color: '#ef4444', marginTop: '0.5rem', display: 'block' }}>
                  Passwords don't match
                </small>
              )}
            </div>

            <button 
              type="submit" 
              className="btn btn-primary btn-block btn-lg" 
              disabled={loading || formData.password !== formData.confirmPassword}
            >
              {loading ? (
                <>
                  <span className="loading-spinner" style={{width: '20px', height: '20px', marginBottom: 0}}></span>
                  Creating account...
                </>
              ) : (
                <>
                  ✨ Create Account
                </>
              )}
            </button>
          </form>

          <div className="auth-footer">
            <p>Already have an account? <Link to="/login">Sign in</Link></p>
          </div>
        </div>

        <div className="auth-features">
          <div className="auth-feature">
            <span>📚</span>
            <span>Share Resources</span>
          </div>
          <div className="auth-feature">
            <span>⭐</span>
            <span>Build Reputation</span>
          </div>
          <div className="auth-feature">
            <span>🤝</span>
            <span>Connect</span>
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
          background: radial-gradient(circle, rgba(6, 182, 212, 0.4) 0%, transparent 70%);
          top: -200px;
          right: -200px;
        }
        
        .orb-2 {
          width: 500px;
          height: 500px;
          background: radial-gradient(circle, rgba(139, 92, 246, 0.3) 0%, transparent 70%);
          bottom: -150px;
          left: -150px;
          animation-delay: -7s;
        }
        
        .orb-3 {
          width: 400px;
          height: 400px;
          background: radial-gradient(circle, rgba(244, 114, 182, 0.25) 0%, transparent 70%);
          top: 40%;
          left: 30%;
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
          background: radial-gradient(circle, rgba(6, 182, 212, 0.4) 0%, transparent 70%);
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
        
        .password-strength {
          display: flex;
          align-items: center;
          gap: 0.75rem;
          margin-top: 0.5rem;
        }
        
        .strength-bars {
          display: flex;
          gap: 4px;
          flex: 1;
        }
        
        .strength-bar {
          height: 4px;
          flex: 1;
          background: rgba(255, 255, 255, 0.1);
          border-radius: 2px;
          transition: all 0.3s ease;
        }
        
        .strength-bar.active {
          background: var(--primary);
        }
        
        .strength-text {
          font-size: 0.75rem;
          font-weight: 500;
          min-width: 60px;
          text-align: right;
        }
      `}</style>
    </div>
  );
};

export default Signup;
