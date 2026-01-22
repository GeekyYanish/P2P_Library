import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import api from '../services/api';

const Reputation = () => {
  const { user } = useAuth();
  const [loading, setLoading] = useState(true);
  const [reputationData, setReputationData] = useState(null);

  useEffect(() => {
    loadReputation();
  }, []);

  const loadReputation = async () => {
    try {
      const response = await api.getReputation();
      if (response.data.success) {
        setReputationData(response.data.data);
      }
    } catch (error) {
      console.log('Using default reputation data');
    } finally {
      setLoading(false);
    }
  };

  // Default/demo data
  const reputation = reputationData || {
    score: user?.reputation || 8.5,
    uploads: 12,
    downloads: 45,
    ratings_given: 23,
    ratings_received: 38,
    rank: 'Scholar',
    percentile: 85
  };

  const achievements = [
    { icon: '🚀', title: 'First Upload', description: 'Uploaded your first resource', unlocked: true },
    { icon: '⭐', title: 'Quality Contributor', description: 'Received 10+ positive ratings', unlocked: true },
    { icon: '🎯', title: 'Helpful Peer', description: 'Rated 20+ resources', unlocked: true },
    { icon: '🏆', title: 'Top Scholar', description: 'Reached top 10% reputation', unlocked: reputation.percentile >= 90 },
    { icon: '💎', title: 'Diamond Sharer', description: 'Uploaded 50+ resources', unlocked: false },
    { icon: '👑', title: 'Legendary', description: 'Achieved maximum reputation', unlocked: false }
  ];

  const activities = [
    { type: 'upload', text: 'Uploaded "Machine Learning Notes"', time: '2 hours ago', points: '+5' },
    { type: 'rating', text: 'Received 5-star rating', time: '5 hours ago', points: '+2' },
    { type: 'download', text: 'Downloaded "Quantum Physics Guide"', time: '1 day ago', points: '-1' },
    { type: 'rating', text: 'Rated "Data Structures" 4 stars', time: '2 days ago', points: '+1' }
  ];

  const getRankColor = (rank) => {
    const colors = {
      'Novice': '#71717a',
      'Contributor': '#06b6d4',
      'Scholar': '#8b5cf6',
      'Expert': '#f59e0b',
      'Master': '#ef4444',
      'Legend': '#f472b6'
    };
    return colors[rank] || colors['Scholar'];
  };

  return (
    <section className="section">
      <div className="container">
        {/* Page Header */}
        <div className="page-header">
          <h1 className="page-title">
            <span>⭐</span> Reputation Dashboard
          </h1>
          <p className="page-subtitle">
            Track your contributions and achievements in the network
          </p>
        </div>

        {/* Main Reputation Card */}
        <div className="reputation-card">
          <div className="reputation-score">
            <div className="score-circle">
              <div className="score-ring" style={{ '--progress': `${(reputation.score / 10) * 100}%` }}></div>
              <div className="score-inner">
                <span className="score-value">{reputation.score.toFixed(1)}</span>
                <span className="score-max">/ 10</span>
              </div>
            </div>
            <div className="rank-badge" style={{ '--rank-color': getRankColor(reputation.rank) }}>
              {reputation.rank}
            </div>
            <p className="percentile-text">Top {100 - reputation.percentile}% of users</p>
          </div>

          <div className="reputation-details">
            <div className="detail-item">
              <span className="detail-icon">📤</span>
              <div className="detail-info">
                <span className="detail-value">{reputation.uploads}</span>
                <span className="detail-label">Uploads</span>
              </div>
            </div>
            <div className="detail-item">
              <span className="detail-icon">⬇️</span>
              <div className="detail-info">
                <span className="detail-value">{reputation.downloads}</span>
                <span className="detail-label">Downloads</span>
              </div>
            </div>
            <div className="detail-item">
              <span className="detail-icon">⭐</span>
              <div className="detail-info">
                <span className="detail-value">{reputation.ratings_given}</span>
                <span className="detail-label">Ratings Given</span>
              </div>
            </div>
            <div className="detail-item">
              <span className="detail-icon">🌟</span>
              <div className="detail-info">
                <span className="detail-value">{reputation.ratings_received}</span>
                <span className="detail-label">Ratings Received</span>
              </div>
            </div>
          </div>
        </div>

        {/* Achievements Section */}
        <div className="achievements-section">
          <h2 className="section-title">🏆 Achievements</h2>
          <div className="achievements-grid">
            {achievements.map((achievement, index) => (
              <div 
                key={index} 
                className={`achievement-card ${achievement.unlocked ? 'unlocked' : 'locked'}`}
              >
                <div className="achievement-icon">{achievement.icon}</div>
                <div className="achievement-info">
                  <h4>{achievement.title}</h4>
                  <p>{achievement.description}</p>
                </div>
                {achievement.unlocked && <span className="achievement-check">✓</span>}
              </div>
            ))}
          </div>
        </div>

        {/* Recent Activity */}
        <div className="activity-section">
          <h2 className="section-title">📊 Recent Activity</h2>
          <div className="activity-list">
            {activities.map((activity, index) => (
              <div key={index} className="activity-item">
                <div className="activity-icon">
                  {activity.type === 'upload' ? '📤' : activity.type === 'rating' ? '⭐' : '⬇️'}
                </div>
                <div className="activity-content">
                  <p className="activity-text">{activity.text}</p>
                  <span className="activity-time">{activity.time}</span>
                </div>
                <span className={`activity-points ${activity.points.startsWith('+') ? 'positive' : 'negative'}`}>
                  {activity.points}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Reputation Info */}
        <div className="reputation-info">
          <h2 className="section-title">💡 How to Increase Reputation</h2>
          <div className="info-grid">
            <div className="info-card">
              <span className="info-icon">📤</span>
              <h4>Upload Resources</h4>
              <p>+5 points per upload</p>
            </div>
            <div className="info-card">
              <span className="info-icon">⭐</span>
              <h4>Rate Resources</h4>
              <p>+1 point per rating</p>
            </div>
            <div className="info-card">
              <span className="info-icon">🌟</span>
              <h4>Get Good Ratings</h4>
              <p>+2 points for 5-star</p>
            </div>
            <div className="info-card">
              <span className="info-icon">🤝</span>
              <h4>Stay Active</h4>
              <p>Daily login bonus</p>
            </div>
          </div>
        </div>
      </div>

      <style>{`
        .reputation-card {
          position: relative;
        }
        
        .score-circle {
          position: relative;
          width: 160px;
          height: 160px;
        }
        
        .score-ring {
          position: absolute;
          inset: 0;
          border-radius: 50%;
          background: conic-gradient(
            var(--primary) 0% var(--progress),
            rgba(255, 255, 255, 0.1) var(--progress) 100%
          );
          mask: radial-gradient(transparent 55%, black 55%);
          -webkit-mask: radial-gradient(transparent 55%, black 55%);
        }
        
        .score-inner {
          position: absolute;
          inset: 0;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          background: var(--bg-card);
          border-radius: 50%;
          margin: 12px;
        }
        
        .rank-badge {
          margin-top: 1rem;
          padding: 0.5rem 1.5rem;
          background: rgba(139, 92, 246, 0.2);
          border: 1px solid var(--rank-color);
          border-radius: 9999px;
          font-weight: 600;
          font-size: 0.9rem;
          color: var(--rank-color);
        }
        
        .percentile-text {
          margin-top: 0.5rem;
          font-size: 0.875rem;
          color: var(--text-muted);
        }
        
        .detail-item {
          display: flex;
          align-items: center;
          gap: 1rem;
          padding: 1.25rem;
          background: var(--bg-glass);
          border-radius: var(--radius-md);
          border: 1px solid var(--border);
        }
        
        .detail-icon {
          font-size: 1.5rem;
        }
        
        .detail-info {
          display: flex;
          flex-direction: column;
        }
        
        .detail-value {
          font-size: 1.5rem;
          font-weight: 700;
        }
        
        .detail-label {
          font-size: 0.875rem;
          color: var(--text-muted);
        }
        
        .achievements-section, .activity-section, .reputation-info {
          margin-top: 3rem;
        }
        
        .section-title {
          font-size: 1.5rem;
          font-weight: 700;
          margin-bottom: 1.5rem;
        }
        
        .achievements-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
          gap: 1rem;
        }
        
        .achievement-card {
          display: flex;
          align-items: center;
          gap: 1rem;
          padding: 1.25rem;
          background: var(--bg-card);
          border: 1px solid var(--border);
          border-radius: var(--radius-md);
          transition: all 0.3s ease;
        }
        
        .achievement-card.unlocked {
          border-color: var(--primary);
          background: rgba(139, 92, 246, 0.05);
        }
        
        .achievement-card.locked {
          opacity: 0.5;
        }
        
        .achievement-card.locked .achievement-icon {
          filter: grayscale(1);
        }
        
        .achievement-icon {
          font-size: 2rem;
        }
        
        .achievement-info h4 {
          font-weight: 600;
          margin-bottom: 0.25rem;
        }
        
        .achievement-info p {
          font-size: 0.875rem;
          color: var(--text-muted);
        }
        
        .achievement-check {
          margin-left: auto;
          width: 24px;
          height: 24px;
          background: var(--success);
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 0.75rem;
          color: white;
        }
        
        .activity-list {
          background: var(--bg-card);
          border: 1px solid var(--border);
          border-radius: var(--radius-lg);
          overflow: hidden;
        }
        
        .activity-item {
          display: flex;
          align-items: center;
          gap: 1rem;
          padding: 1.25rem 1.5rem;
          border-bottom: 1px solid var(--border);
          transition: background 0.2s;
        }
        
        .activity-item:last-child {
          border-bottom: none;
        }
        
        .activity-item:hover {
          background: var(--bg-glass);
        }
        
        .activity-icon {
          font-size: 1.5rem;
        }
        
        .activity-content {
          flex: 1;
        }
        
        .activity-text {
          margin-bottom: 0.25rem;
        }
        
        .activity-time {
          font-size: 0.875rem;
          color: var(--text-muted);
        }
        
        .activity-points {
          font-weight: 700;
          padding: 0.25rem 0.75rem;
          border-radius: var(--radius-sm);
        }
        
        .activity-points.positive {
          background: rgba(16, 185, 129, 0.1);
          color: var(--success);
        }
        
        .activity-points.negative {
          background: rgba(239, 68, 68, 0.1);
          color: var(--danger);
        }
        
        .info-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
          gap: 1rem;
        }
        
        .info-card {
          text-align: center;
          padding: 1.5rem;
          background: var(--bg-card);
          border: 1px solid var(--border);
          border-radius: var(--radius-md);
          transition: all 0.3s ease;
        }
        
        .info-card:hover {
          border-color: var(--primary);
          transform: translateY(-3px);
        }
        
        .info-card .info-icon {
          font-size: 2rem;
          display: block;
          margin-bottom: 0.75rem;
        }
        
        .info-card h4 {
          font-weight: 600;
          margin-bottom: 0.25rem;
        }
        
        .info-card p {
          font-size: 0.875rem;
          color: var(--text-muted);
        }
      `}</style>
    </section>
  );
};

export default Reputation;
