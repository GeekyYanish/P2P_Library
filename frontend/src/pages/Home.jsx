import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../services/api';
import { useToast } from '../context/ToastContext';
import { useAuth } from '../context/AuthContext';

const Home = () => {
  const [stats, setStats] = useState({ peers: 0, files: 0, downloads: 0, avgRating: 0 });
  const [loading, setLoading] = useState(true);
  const { showToast } = useToast();
  const { user } = useAuth();

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const response = await api.getStats();
      if (response.data.success) {
        setStats(response.data.data);
      }
    } catch (error) {
      console.log('Stats not available yet');
    } finally {
      setLoading(false);
    }
  };

  const features = [
    {
      icon: '🌐',
      title: 'Decentralized Network',
      description: 'Share resources directly with peers without central servers'
    },
    {
      icon: '⭐',
      title: 'Reputation System',
      description: 'Build your reputation by contributing quality resources'
    },
    {
      icon: '🔒',
      title: 'Secure Sharing',
      description: 'End-to-end encrypted transfers with verified peers'
    },
    {
      icon: '⚡',
      title: 'Lightning Fast',
      description: 'Optimized P2P transfers for maximum speed'
    }
  ];

  return (
    <section className="section">
      <div className="container">
        {/* Hero Section */}
        <div className="hero">
          <div className="hero-particles">
            {[...Array(20)].map((_, i) => (
              <div key={i} className="particle" style={{
                '--delay': `${Math.random() * 5}s`,
                '--x': `${Math.random() * 100}%`,
                '--duration': `${15 + Math.random() * 10}s`
              }}></div>
            ))}
          </div>
          
          <div className="hero-content">
            <div className="hero-badge">
              <span className="badge-dot"></span>
              Welcome back, {user?.username || 'Scholar'}!
            </div>
            
            <h1 className="hero-title">
              <span>Decentralized</span>
              <br />
              <span className="gradient-text">Academic Resource Sharing</span>
            </h1>
            
            <p className="hero-subtitle">
              Share and discover academic resources in a peer-to-peer network.
              Contribute to earn reputation and unlock unlimited downloads.
            </p>
            
            <div className="hero-actions">
              <Link to="/library" className="btn btn-primary btn-lg">
                <span>🔍</span> Browse Library
              </Link>
              <Link to="/upload" className="btn btn-secondary btn-lg">
                <span>📤</span> Share Resources
              </Link>
            </div>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-icon">👥</div>
            <div className="stat-value">{loading ? '—' : stats.peers || 0}</div>
            <div className="stat-label">Online Peers</div>
            <div className="stat-trend positive">
              <span>↑</span> Active now
            </div>
          </div>
          
          <div className="stat-card">
            <div className="stat-icon">📄</div>
            <div className="stat-value">{loading ? '—' : stats.files || 0}</div>
            <div className="stat-label">Shared Files</div>
            <div className="stat-trend">
              <span>📚</span> Resources
            </div>
          </div>
          
          <div className="stat-card">
            <div className="stat-icon">⬇️</div>
            <div className="stat-value">{loading ? '—' : stats.downloads || 0}</div>
            <div className="stat-label">Total Downloads</div>
            <div className="stat-trend positive">
              <span>↑</span> This week
            </div>
          </div>
          
          <div className="stat-card">
            <div className="stat-icon">⭐</div>
            <div className="stat-value">{loading ? '—' : (stats.avgRating || 0).toFixed(1)}</div>
            <div className="stat-label">Avg. Rating</div>
            <div className="stat-trend">
              <span>✨</span> Quality
            </div>
          </div>
        </div>

        {/* Features Section */}
        <div className="features-section">
          <h2 className="section-title">
            <span className="gradient-text">Why Choose Knowledge Exchange?</span>
          </h2>
          
          <div className="features-grid">
            {features.map((feature, index) => (
              <div key={index} className="feature-card" style={{ '--delay': `${index * 0.1}s` }}>
                <div className="feature-icon">{feature.icon}</div>
                <h3 className="feature-title">{feature.title}</h3>
                <p className="feature-desc">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Quick Actions */}
        <div className="quick-actions">
          <h2 className="section-title">Quick Actions</h2>
          
          <div className="actions-grid">
            <Link to="/library" className="action-card">
              <div className="action-icon">📚</div>
              <div className="action-content">
                <h3>Browse Library</h3>
                <p>Discover academic resources shared by peers</p>
              </div>
              <span className="action-arrow">→</span>
            </Link>
            
            <Link to="/upload" className="action-card">
              <div className="action-icon">📤</div>
              <div className="action-content">
                <h3>Upload Resources</h3>
                <p>Share your knowledge and earn reputation</p>
              </div>
              <span className="action-arrow">→</span>
            </Link>
            
            <Link to="/peers" className="action-card">
              <div className="action-icon">🌐</div>
              <div className="action-content">
                <h3>View Network</h3>
                <p>Connect with peers in the network</p>
              </div>
              <span className="action-arrow">→</span>
            </Link>
            
            <Link to="/reputation" className="action-card">
              <div className="action-icon">⭐</div>
              <div className="action-content">
                <h3>Your Reputation</h3>
                <p>Track your contributions and achievements</p>
              </div>
              <span className="action-arrow">→</span>
            </Link>
          </div>
        </div>
      </div>

      <style>{`
        .hero {
          position: relative;
          overflow: hidden;
        }
        
        .hero-particles {
          position: absolute;
          inset: 0;
          overflow: hidden;
          pointer-events: none;
        }
        
        .particle {
          position: absolute;
          width: 4px;
          height: 4px;
          background: var(--primary);
          border-radius: 50%;
          left: var(--x);
          bottom: -10px;
          opacity: 0.5;
          animation: rise var(--duration) linear infinite;
          animation-delay: var(--delay);
        }
        
        @keyframes rise {
          0% {
            transform: translateY(0) scale(1);
            opacity: 0;
          }
          10% {
            opacity: 0.5;
          }
          90% {
            opacity: 0.3;
          }
          100% {
            transform: translateY(-600px) scale(0.5);
            opacity: 0;
          }
        }
        
        .hero-badge {
          display: inline-flex;
          align-items: center;
          gap: 0.5rem;
          padding: 0.5rem 1rem;
          background: rgba(139, 92, 246, 0.1);
          border: 1px solid rgba(139, 92, 246, 0.2);
          border-radius: 9999px;
          font-size: 0.875rem;
          color: var(--primary-light);
          margin-bottom: 1.5rem;
        }
        
        .badge-dot {
          width: 8px;
          height: 8px;
          background: var(--success);
          border-radius: 50%;
          animation: pulse 2s infinite;
        }
        
        .gradient-text {
          background: var(--gradient-primary);
          -webkit-background-clip: text;
          background-clip: text;
          color: transparent;
        }
        
        .stat-trend {
          display: flex;
          align-items: center;
          gap: 0.25rem;
          font-size: 0.75rem;
          color: var(--text-muted);
          margin-top: 0.5rem;
        }
        
        .stat-trend.positive {
          color: var(--success);
        }
        
        .features-section {
          margin-bottom: 4rem;
        }
        
        .section-title {
          font-size: 1.75rem;
          font-weight: 700;
          text-align: center;
          margin-bottom: 2rem;
        }
        
        .features-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
          gap: 1.5rem;
        }
        
        .feature-card {
          background: var(--bg-card);
          border: 1px solid var(--border);
          border-radius: var(--radius-lg);
          padding: 2rem;
          text-align: center;
          transition: all 0.3s ease;
          animation: fadeInUp 0.5s ease forwards;
          animation-delay: var(--delay);
          opacity: 0;
        }
        
        @keyframes fadeInUp {
          from {
            opacity: 0;
            transform: translateY(20px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
        
        .feature-card:hover {
          transform: translateY(-5px);
          border-color: var(--primary);
          box-shadow: 0 20px 40px rgba(139, 92, 246, 0.15);
        }
        
        .feature-icon {
          font-size: 3rem;
          margin-bottom: 1rem;
        }
        
        .feature-title {
          font-size: 1.125rem;
          font-weight: 600;
          margin-bottom: 0.5rem;
        }
        
        .feature-desc {
          color: var(--text-secondary);
          font-size: 0.9rem;
          line-height: 1.6;
        }
        
        .quick-actions {
          margin-bottom: 2rem;
        }
        
        .actions-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
          gap: 1rem;
        }
        
        .action-card {
          display: flex;
          align-items: center;
          gap: 1rem;
          padding: 1.25rem 1.5rem;
          background: var(--bg-card);
          border: 1px solid var(--border);
          border-radius: var(--radius-lg);
          transition: all 0.3s ease;
          cursor: pointer;
        }
        
        .action-card:hover {
          border-color: var(--primary);
          background: var(--bg-card-hover);
          transform: translateX(5px);
        }
        
        .action-card:hover .action-arrow {
          transform: translateX(5px);
          color: var(--primary);
        }
        
        .action-icon {
          font-size: 2rem;
          width: 50px;
          height: 50px;
          display: flex;
          align-items: center;
          justify-content: center;
          background: rgba(139, 92, 246, 0.1);
          border-radius: var(--radius-md);
        }
        
        .action-content {
          flex: 1;
        }
        
        .action-content h3 {
          font-size: 1rem;
          font-weight: 600;
          margin-bottom: 0.25rem;
        }
        
        .action-content p {
          font-size: 0.875rem;
          color: var(--text-muted);
        }
        
        .action-arrow {
          font-size: 1.25rem;
          color: var(--text-muted);
          transition: all 0.3s ease;
        }
      `}</style>
    </section>
  );
};

export default Home;
