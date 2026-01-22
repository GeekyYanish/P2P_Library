import { NavLink } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useState } from 'react';

const Header = () => {
  const { user, logout } = useAuth();
  const [showDropdown, setShowDropdown] = useState(false);

  const getInitials = (name) => {
    if (!name) return 'U';
    return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
  };

  return (
    <header className="header">
      <div className="container">
        {/* Logo */}
        <NavLink to="/home" className="logo">
          <span className="logo-icon">📚</span>
          <h1>Knowledge Exchange</h1>
        </NavLink>

        {/* Navigation */}
        <nav className="nav">
          <NavLink to="/home" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
            <span className="nav-icon">🏠</span>
            <span className="nav-text">Home</span>
          </NavLink>
          <NavLink to="/library" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
            <span className="nav-icon">📚</span>
            <span className="nav-text">Library</span>
          </NavLink>
          <NavLink to="/upload" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
            <span className="nav-icon">📤</span>
            <span className="nav-text">Upload</span>
          </NavLink>
          <NavLink to="/reputation" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
            <span className="nav-icon">⭐</span>
            <span className="nav-text">Reputation</span>
          </NavLink>
          <NavLink to="/peers" className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}>
            <span className="nav-icon">🌐</span>
            <span className="nav-text">Peers</span>
          </NavLink>
        </nav>

        {/* User Menu */}
        <div className="user-menu">
          <div className="connection-status">
            <span className="status-indicator online"></span>
            <span className="status-text">Connected</span>
          </div>
          
          <div 
            className="user-avatar-container"
            onClick={() => setShowDropdown(!showDropdown)}
          >
            <div className="user-avatar">
              {getInitials(user?.username)}
            </div>
            <span className="user-name">{user?.username || 'User'}</span>
            <span className="dropdown-arrow">▼</span>
            
            {showDropdown && (
              <div className="user-dropdown">
                <div className="dropdown-header">
                  <div className="dropdown-avatar">
                    {getInitials(user?.username)}
                  </div>
                  <div className="dropdown-info">
                    <strong>{user?.username}</strong>
                    <span>{user?.email}</span>
                  </div>
                </div>
                <div className="dropdown-divider"></div>
                <NavLink to="/reputation" className="dropdown-item" onClick={() => setShowDropdown(false)}>
                  <span>⭐</span> My Reputation
                </NavLink>
                <NavLink to="/upload" className="dropdown-item" onClick={() => setShowDropdown(false)}>
                  <span>📤</span> Upload Resource
                </NavLink>
                <div className="dropdown-divider"></div>
                <button className="dropdown-item danger" onClick={logout}>
                  <span>🚪</span> Sign Out
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      <style>{`
        .nav-icon {
          font-size: 1rem;
        }
        
        .nav-text {
          font-size: 0.9rem;
        }
        
        .nav-link {
          display: flex;
          align-items: center;
          gap: 0.4rem;
        }
        
        .user-menu {
          display: flex;
          align-items: center;
          gap: 1rem;
        }
        
        .user-avatar-container {
          display: flex;
          align-items: center;
          gap: 0.75rem;
          padding: 0.5rem 1rem;
          background: var(--bg-glass);
          border: 1px solid var(--border);
          border-radius: 9999px;
          cursor: pointer;
          transition: all 0.2s;
          position: relative;
        }
        
        .user-avatar-container:hover {
          background: var(--bg-glass-strong);
          border-color: var(--border-light);
        }
        
        .user-avatar {
          width: 32px;
          height: 32px;
          background: var(--gradient-primary);
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: 600;
          font-size: 0.8rem;
          color: white;
        }
        
        .user-name {
          font-weight: 500;
          font-size: 0.9rem;
        }
        
        .dropdown-arrow {
          font-size: 0.6rem;
          color: var(--text-muted);
          transition: transform 0.2s;
        }
        
        .user-avatar-container:hover .dropdown-arrow {
          transform: rotate(180deg);
        }
        
        .user-dropdown {
          position: absolute;
          top: calc(100% + 0.5rem);
          right: 0;
          width: 260px;
          background: var(--bg-card);
          border: 1px solid var(--border);
          border-radius: var(--radius-lg);
          box-shadow: var(--shadow-lg);
          z-index: 100;
          animation: slideDown 0.2s ease;
          overflow: hidden;
        }
        
        @keyframes slideDown {
          from {
            opacity: 0;
            transform: translateY(-10px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
        
        .dropdown-header {
          display: flex;
          align-items: center;
          gap: 1rem;
          padding: 1rem;
        }
        
        .dropdown-avatar {
          width: 48px;
          height: 48px;
          background: var(--gradient-primary);
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: 600;
          color: white;
        }
        
        .dropdown-info {
          flex: 1;
          overflow: hidden;
        }
        
        .dropdown-info strong {
          display: block;
          margin-bottom: 0.25rem;
        }
        
        .dropdown-info span {
          font-size: 0.875rem;
          color: var(--text-muted);
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
          display: block;
        }
        
        .dropdown-divider {
          height: 1px;
          background: var(--border);
        }
        
        .dropdown-item {
          display: flex;
          align-items: center;
          gap: 0.75rem;
          padding: 0.75rem 1rem;
          font-size: 0.9rem;
          color: var(--text-secondary);
          transition: all 0.2s;
          width: 100%;
          border: none;
          background: none;
          cursor: pointer;
          text-align: left;
        }
        
        .dropdown-item:hover {
          background: var(--bg-glass);
          color: var(--text-primary);
        }
        
        .dropdown-item.danger {
          color: var(--danger);
        }
        
        .dropdown-item.danger:hover {
          background: rgba(239, 68, 68, 0.1);
        }
        
        @media (max-width: 900px) {
          .nav-text {
            display: none;
          }
          
          .nav-link {
            padding: 0.5rem 0.75rem;
          }
          
          .user-name {
            display: none;
          }
          
          .connection-status .status-text {
            display: none;
          }
        }
      `}</style>
    </header>
  );
};

export default Header;
