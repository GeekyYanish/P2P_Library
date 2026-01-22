const Footer = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="footer">
      <div className="container">
        <div className="footer-content">
          {/* Brand */}
          <div className="footer-brand">
            <span className="footer-logo">📚</span>
            <span className="footer-title">The Knowledge Exchange</span>
          </div>
          
          {/* Tagline */}
          <p className="footer-tagline">
            Decentralized P2P Academic Library — Share Knowledge, Build Community
          </p>
          
          {/* Tech Stack */}
          <div className="footer-tech">
            <span>Built with</span>
            <div className="tech-badges">
              <span className="tech-badge">
                <span className="tech-icon">🔵</span> Go
              </span>
              <span className="tech-badge">
                <span className="tech-icon">⚛️</span> React
              </span>
              <span className="tech-badge">
                <span className="tech-icon">⚡</span> Vite
              </span>
            </div>
          </div>
          
          {/* Copyright */}
          <p className="footer-copyright">
            © {currentYear} Knowledge Exchange. Open Source & Decentralized.
          </p>
        </div>
      </div>

      <style>{`
        .footer-content {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 1rem;
        }
        
        .footer-brand {
          display: flex;
          align-items: center;
          gap: 0.75rem;
        }
        
        .footer-logo {
          font-size: 2rem;
          filter: drop-shadow(0 0 10px var(--primary-glow));
        }
        
        .footer-title {
          font-size: 1.25rem;
          font-weight: 700;
          background: var(--gradient-primary);
          -webkit-background-clip: text;
          background-clip: text;
          color: transparent;
        }
        
        .footer-tagline {
          color: var(--text-secondary);
          font-size: 0.95rem;
        }
        
        .footer-tech {
          display: flex;
          align-items: center;
          gap: 1rem;
          margin: 0.5rem 0;
        }
        
        .footer-tech > span {
          color: var(--text-muted);
          font-size: 0.875rem;
        }
        
        .tech-badges {
          display: flex;
          gap: 0.5rem;
        }
        
        .tech-badge {
          display: inline-flex;
          align-items: center;
          gap: 0.35rem;
          padding: 0.35rem 0.75rem;
          background: var(--bg-glass);
          border: 1px solid var(--border);
          border-radius: 9999px;
          font-size: 0.8rem;
          color: var(--text-secondary);
          transition: all 0.2s;
        }
        
        .tech-badge:hover {
          border-color: var(--primary);
          background: rgba(139, 92, 246, 0.1);
          color: var(--text-primary);
        }
        
        .tech-icon {
          font-size: 0.9rem;
        }
        
        .footer-copyright {
          color: var(--text-muted);
          font-size: 0.8rem;
          margin-top: 0.5rem;
        }
        
        @media (max-width: 600px) {
          .footer-tech {
            flex-direction: column;
            gap: 0.5rem;
          }
        }
      `}</style>
    </footer>
  );
};

export default Footer;
