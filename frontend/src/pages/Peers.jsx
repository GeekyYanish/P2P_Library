import { useState, useEffect } from 'react';
import api from '../services/api';

const Peers = () => {
  const [peers, setPeers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPeers();
  }, []);

  const loadPeers = async () => {
    try {
      const response = await api.getPeers();
      if (response.data.success) {
        setPeers(response.data.data || []);
      }
    } catch (error) {
      console.log('Using demo peers');
    } finally {
      setLoading(false);
    }
  };

  // Demo peers for showcase
  const demoPeers = [
    { id: '1', name: 'Alice Chen', reputation: 9.2, is_online: true, uploads: 45, downloads: 120, files_shared: 32 },
    { id: '2', name: 'Bob Smith', reputation: 8.7, is_online: true, uploads: 38, downloads: 95, files_shared: 24 },
    { id: '3', name: 'Carol Williams', reputation: 8.4, is_online: false, uploads: 29, downloads: 67, files_shared: 18 },
    { id: '4', name: 'David Lee', reputation: 9.5, is_online: true, uploads: 67, downloads: 145, files_shared: 48 },
    { id: '5', name: 'Emma Davis', reputation: 7.8, is_online: true, uploads: 22, downloads: 54, files_shared: 15 },
    { id: '6', name: 'Frank Miller', reputation: 8.1, is_online: false, uploads: 31, downloads: 78, files_shared: 21 },
    { id: '7', name: 'Grace Taylor', reputation: 9.0, is_online: true, uploads: 52, downloads: 110, files_shared: 36 },
    { id: '8', name: 'Henry Wilson', reputation: 7.5, is_online: true, uploads: 18, downloads: 42, files_shared: 12 }
  ];

  const displayPeers = peers.length > 0 ? peers : demoPeers;
  const onlinePeers = displayPeers.filter(p => p.is_online);
  const offlinePeers = displayPeers.filter(p => !p.is_online);

  const getInitials = (name) => {
    return name.split(' ').map(n => n[0]).join('').toUpperCase();
  };

  const getAvatarColor = (name) => {
    const colors = [
      'linear-gradient(135deg, #8b5cf6, #6366f1)',
      'linear-gradient(135deg, #06b6d4, #0891b2)',
      'linear-gradient(135deg, #f472b6, #ec4899)',
      'linear-gradient(135deg, #10b981, #059669)',
      'linear-gradient(135deg, #f59e0b, #d97706)'
    ];
    return colors[name.length % colors.length];
  };

  return (
    <section className="section">
      <div className="container">
        {/* Page Header */}
        <div className="page-header">
          <h1 className="page-title">
            <span>🌐</span> Network Peers
          </h1>
          <p className="page-subtitle">
            Connect with other users in the P2P network
          </p>
        </div>

        {/* Network Stats */}
        <div className="network-stats">
          <div className="network-stat">
            <div className="stat-circle online">
              <span>{onlinePeers.length}</span>
            </div>
            <div className="stat-info">
              <span className="stat-label">Online Now</span>
              <div className="pulse-indicator"></div>
            </div>
          </div>
          <div className="network-stat">
            <div className="stat-circle total">
              <span>{displayPeers.length}</span>
            </div>
            <div className="stat-info">
              <span className="stat-label">Total Peers</span>
            </div>
          </div>
          <div className="network-stat">
            <div className="stat-circle files">
              <span>{displayPeers.reduce((acc, p) => acc + (p.files_shared || 0), 0)}</span>
            </div>
            <div className="stat-info">
              <span className="stat-label">Shared Files</span>
            </div>
          </div>
        </div>

        {/* Network Visualization */}
        <div className="network-visual">
          <div className="network-center">
            <span>🌐</span>
            <p>Network Hub</p>
          </div>
          <div className="network-nodes">
            {onlinePeers.slice(0, 6).map((peer, index) => (
              <div 
                key={peer.id} 
                className="network-node"
                style={{ 
                  '--angle': `${(index * 60) - 60}deg`,
                  '--delay': `${index * 0.1}s`
                }}
              >
                <div 
                  className="node-avatar"
                  style={{ background: getAvatarColor(peer.name) }}
                >
                  {getInitials(peer.name)}
                </div>
              </div>
            ))}
          </div>
          <div className="connection-lines">
            {onlinePeers.slice(0, 6).map((_, index) => (
              <div 
                key={index} 
                className="connection-line"
                style={{ '--angle': `${(index * 60) - 60}deg` }}
              ></div>
            ))}
          </div>
        </div>

        {/* Online Peers */}
        <div className="peers-section">
          <h2 className="section-title">
            <span className="status-dot online"></span>
            Online Peers ({onlinePeers.length})
          </h2>
          <div className="peers-grid">
            {onlinePeers.map((peer, index) => (
              <div 
                key={peer.id} 
                className="peer-card"
                style={{ '--delay': `${index * 0.05}s` }}
              >
                <div className="peer-header">
                  <div 
                    className="peer-avatar"
                    style={{ background: getAvatarColor(peer.name) }}
                  >
                    {getInitials(peer.name)}
                  </div>
                  <div className="peer-info">
                    <h4>{peer.name}</h4>
                    <div className="peer-status">
                      <span className="status-indicator online"></span>
                      <span>Online</span>
                    </div>
                  </div>
                  <div className="peer-reputation">
                    <span className="rep-value">{peer.reputation.toFixed(1)}</span>
                    <span className="rep-label">Rep</span>
                  </div>
                </div>
                <div className="peer-stats">
                  <div className="peer-stat">
                    <span className="peer-stat-value">{peer.uploads}</span>
                    <span className="peer-stat-label">Uploads</span>
                  </div>
                  <div className="peer-stat">
                    <span className="peer-stat-value">{peer.downloads}</span>
                    <span className="peer-stat-label">Downloads</span>
                  </div>
                  <div className="peer-stat">
                    <span className="peer-stat-value">{peer.files_shared}</span>
                    <span className="peer-stat-label">Files</span>
                  </div>
                </div>
                <div className="peer-actions">
                  <button className="btn btn-primary btn-sm">View Files</button>
                  <button className="btn btn-secondary btn-sm">Connect</button>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Offline Peers */}
        {offlinePeers.length > 0 && (
          <div className="peers-section">
            <h2 className="section-title">
              <span className="status-dot offline"></span>
              Offline Peers ({offlinePeers.length})
            </h2>
            <div className="peers-grid">
              {offlinePeers.map((peer, index) => (
                <div 
                  key={peer.id} 
                  className="peer-card offline"
                  style={{ '--delay': `${index * 0.05}s` }}
                >
                  <div className="peer-header">
                    <div 
                      className="peer-avatar"
                      style={{ background: getAvatarColor(peer.name), opacity: 0.5 }}
                    >
                      {getInitials(peer.name)}
                    </div>
                    <div className="peer-info">
                      <h4>{peer.name}</h4>
                      <div className="peer-status">
                        <span className="status-indicator offline"></span>
                        <span>Offline</span>
                      </div>
                    </div>
                    <div className="peer-reputation">
                      <span className="rep-value">{peer.reputation.toFixed(1)}</span>
                      <span className="rep-label">Rep</span>
                    </div>
                  </div>
                  <div className="peer-stats">
                    <div className="peer-stat">
                      <span className="peer-stat-value">{peer.uploads}</span>
                      <span className="peer-stat-label">Uploads</span>
                    </div>
                    <div className="peer-stat">
                      <span className="peer-stat-value">{peer.downloads}</span>
                      <span className="peer-stat-label">Downloads</span>
                    </div>
                    <div className="peer-stat">
                      <span className="peer-stat-value">{peer.files_shared}</span>
                      <span className="peer-stat-label">Files</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      <style>{`
        .network-stats {
          display: flex;
          justify-content: center;
          gap: 3rem;
          margin-bottom: 3rem;
        }
        
        .network-stat {
          display: flex;
          align-items: center;
          gap: 1rem;
        }
        
        .stat-circle {
          width: 60px;
          height: 60px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 1.5rem;
          font-weight: 700;
        }
        
        .stat-circle.online {
          background: rgba(16, 185, 129, 0.15);
          border: 2px solid var(--success);
          color: var(--success);
        }
        
        .stat-circle.total {
          background: rgba(139, 92, 246, 0.15);
          border: 2px solid var(--primary);
          color: var(--primary);
        }
        
        .stat-circle.files {
          background: rgba(6, 182, 212, 0.15);
          border: 2px solid var(--secondary);
          color: var(--secondary);
        }
        
        .stat-info {
          display: flex;
          flex-direction: column;
          gap: 0.25rem;
        }
        
        .stat-info .stat-label {
          color: var(--text-secondary);
          font-size: 0.9rem;
        }
        
        .pulse-indicator {
          width: 8px;
          height: 8px;
          background: var(--success);
          border-radius: 50%;
          animation: pulse 2s infinite;
        }
        
        .network-visual {
          position: relative;
          width: 300px;
          height: 300px;
          margin: 0 auto 3rem;
        }
        
        .network-center {
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          width: 80px;
          height: 80px;
          background: var(--bg-card);
          border: 2px solid var(--primary);
          border-radius: 50%;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          z-index: 2;
          box-shadow: 0 0 30px var(--primary-glow);
        }
        
        .network-center span {
          font-size: 2rem;
        }
        
        .network-center p {
          font-size: 0.7rem;
          color: var(--text-muted);
        }
        
        .network-nodes {
          position: absolute;
          inset: 0;
        }
        
        .network-node {
          position: absolute;
          top: 50%;
          left: 50%;
          transform: rotate(var(--angle)) translateY(-120px) rotate(calc(-1 * var(--angle)));
          animation: fadeIn 0.5s ease forwards;
          animation-delay: var(--delay);
          opacity: 0;
        }
        
        .node-avatar {
          width: 48px;
          height: 48px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: 600;
          font-size: 0.9rem;
          color: white;
          box-shadow: 0 0 20px rgba(0, 0, 0, 0.3);
        }
        
        .connection-lines {
          position: absolute;
          inset: 0;
          z-index: 1;
        }
        
        .connection-line {
          position: absolute;
          top: 50%;
          left: 50%;
          width: 80px;
          height: 2px;
          background: linear-gradient(90deg, var(--primary), transparent);
          transform-origin: left center;
          transform: rotate(var(--angle));
          opacity: 0.3;
        }
        
        .peers-section {
          margin-bottom: 3rem;
        }
        
        .section-title {
          display: flex;
          align-items: center;
          gap: 0.75rem;
          font-size: 1.25rem;
          font-weight: 600;
          margin-bottom: 1.5rem;
        }
        
        .status-dot {
          width: 10px;
          height: 10px;
          border-radius: 50%;
        }
        
        .status-dot.online {
          background: var(--success);
          box-shadow: 0 0 10px var(--success);
        }
        
        .status-dot.offline {
          background: var(--text-muted);
        }
        
        .peer-card {
          animation: slideUp 0.4s ease forwards;
          animation-delay: var(--delay);
          opacity: 0;
        }
        
        @keyframes slideUp {
          from {
            opacity: 0;
            transform: translateY(20px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
        
        .peer-card.offline {
          opacity: 0.6;
        }
        
        .peer-header {
          display: flex;
          align-items: center;
          gap: 1rem;
          margin-bottom: 1rem;
        }
        
        .peer-avatar {
          width: 48px;
          height: 48px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: 600;
          color: white;
        }
        
        .peer-info {
          flex: 1;
        }
        
        .peer-info h4 {
          font-weight: 600;
          margin-bottom: 0.25rem;
        }
        
        .peer-status {
          display: flex;
          align-items: center;
          gap: 0.5rem;
          font-size: 0.875rem;
          color: var(--text-muted);
        }
        
        .peer-reputation {
          text-align: center;
          padding: 0.5rem 1rem;
          background: rgba(139, 92, 246, 0.1);
          border-radius: var(--radius-sm);
        }
        
        .rep-value {
          display: block;
          font-size: 1.25rem;
          font-weight: 700;
          color: var(--primary);
        }
        
        .rep-label {
          font-size: 0.75rem;
          color: var(--text-muted);
        }
        
        .peer-stats {
          display: grid;
          grid-template-columns: repeat(3, 1fr);
          gap: 0.5rem;
          padding: 1rem 0;
          border-top: 1px solid var(--border);
          border-bottom: 1px solid var(--border);
          margin-bottom: 1rem;
        }
        
        .peer-stat {
          text-align: center;
        }
        
        .peer-stat-value {
          display: block;
          font-weight: 700;
        }
        
        .peer-stat-label {
          font-size: 0.75rem;
          color: var(--text-muted);
        }
        
        .peer-actions {
          display: flex;
          gap: 0.5rem;
        }
        
        .peer-actions .btn {
          flex: 1;
        }
        
        @media (max-width: 768px) {
          .network-stats {
            flex-direction: column;
            align-items: center;
            gap: 1.5rem;
          }
          
          .network-visual {
            transform: scale(0.8);
            margin-bottom: 1rem;
          }
        }
      `}</style>
    </section>
  );
};

export default Peers;
