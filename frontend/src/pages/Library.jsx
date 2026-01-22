import { useState, useEffect } from 'react';
import api from '../services/api';
import { useToast } from '../context/ToastContext';

const Library = () => {
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [activeFilter, setActiveFilter] = useState('all');
  const { showToast } = useToast();

  const filters = ['all', 'pdf', 'doc', 'notes', 'research'];

  useEffect(() => {
    loadFiles();
  }, []);

  const loadFiles = async () => {
    try {
      const response = await api.getFiles();
      if (response.data.success) {
        setFiles(response.data.data || []);
      }
    } catch (error) {
      console.log('No files available yet');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      loadFiles();
      return;
    }
    
    setLoading(true);
    try {
      const response = await api.searchFiles(searchQuery);
      if (response.data.success) {
        setFiles(response.data.data || []);
      }
    } catch (error) {
      showToast('Search failed', 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = async (file) => {
    showToast(`Downloading ${file.name}...`, 'success');
    // In a real app, this would trigger the actual download
  };

  const filteredFiles = files.filter(file => {
    if (activeFilter === 'all') return true;
    return file.type?.toLowerCase().includes(activeFilter);
  });

  // Demo files for showcase
  const demoFiles = [
    {
      cid: '1',
      name: 'Advanced Machine Learning Techniques',
      type: 'PDF',
      subject: 'Computer Science',
      size: 2500000,
      rating: 4.8,
      downloads: 342,
      owner_id: 'peer1',
      available: true
    },
    {
      cid: '2',
      name: 'Quantum Computing Fundamentals',
      type: 'PDF',
      subject: 'Physics',
      size: 1800000,
      rating: 4.6,
      downloads: 256,
      owner_id: 'peer2',
      available: true
    },
    {
      cid: '3',
      name: 'Data Structures & Algorithms Notes',
      type: 'DOC',
      subject: 'Computer Science',
      size: 950000,
      rating: 4.9,
      downloads: 523,
      owner_id: 'peer1',
      available: true
    },
    {
      cid: '4',
      name: 'Organic Chemistry Research Paper',
      type: 'PDF',
      subject: 'Chemistry',
      size: 3200000,
      rating: 4.5,
      downloads: 189,
      owner_id: 'peer3',
      available: false
    },
    {
      cid: '5',
      name: 'Neural Networks Deep Dive',
      type: 'PDF',
      subject: 'AI/ML',
      size: 4100000,
      rating: 4.7,
      downloads: 412,
      owner_id: 'peer2',
      available: true
    },
    {
      cid: '6',
      name: 'Linear Algebra Study Guide',
      type: 'NOTES',
      subject: 'Mathematics',
      size: 720000,
      rating: 4.4,
      downloads: 298,
      owner_id: 'peer4',
      available: true
    }
  ];

  const displayFiles = filteredFiles.length > 0 ? filteredFiles : demoFiles;

  const formatSize = (bytes) => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  };

  return (
    <section className="section">
      <div className="container">
        {/* Page Header */}
        <div className="page-header">
          <h1 className="page-title">
            <span>📚</span> Academic Library
          </h1>
          <p className="page-subtitle">
            Discover and download academic resources shared by the community
          </p>
        </div>

        {/* Search & Filters */}
        <div className="library-controls">
          <div className="search-container">
            <div className="search-input">
              <input
                type="text"
                placeholder="Search for resources..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
              />
            </div>
            <button className="btn btn-primary" onClick={handleSearch}>
              🔍 Search
            </button>
          </div>

          <div className="filters">
            {filters.map(filter => (
              <button
                key={filter}
                className={`filter-chip ${activeFilter === filter ? 'active' : ''}`}
                onClick={() => setActiveFilter(filter)}
              >
                {filter.charAt(0).toUpperCase() + filter.slice(1)}
              </button>
            ))}
          </div>
        </div>

        {/* Files Grid */}
        {loading ? (
          <div className="loading">
            <div className="loading-spinner"></div>
            <p>Loading resources...</p>
          </div>
        ) : (
          <>
            <div className="files-header">
              <span className="files-count">{displayFiles.length} resources found</span>
            </div>
            
            <div className="files-grid">
              {displayFiles.map((file, index) => (
                <div 
                  key={file.cid || index} 
                  className={`file-card ${!file.available ? 'unavailable' : ''}`}
                  style={{ '--delay': `${index * 0.05}s` }}
                >
                  <div className="file-header">
                    <span className="file-type">
                      {file.type === 'PDF' ? '📄' : file.type === 'DOC' ? '📝' : '📋'}
                      {file.type}
                    </span>
                    <div className="file-rating">
                      ⭐ {file.rating?.toFixed(1) || '—'}
                    </div>
                  </div>
                  
                  <h3 className="file-title">{file.name}</h3>
                  
                  <div className="file-meta">
                    <span>📁 {file.subject || 'General'}</span>
                    <span>💾 {formatSize(file.size || 0)}</span>
                    <span>⬇️ {file.downloads || 0}</span>
                    <span className={`status ${file.available ? 'available' : 'offline'}`}>
                      {file.available ? '🟢 Available' : '🔴 Offline'}
                    </span>
                  </div>
                  
                  <div className="file-actions">
                    <button 
                      className="btn btn-primary btn-sm"
                      onClick={() => handleDownload(file)}
                      disabled={!file.available}
                    >
                      ⬇️ Download
                    </button>
                    <button className="btn btn-secondary btn-sm">
                      ⭐ Rate
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </>
        )}
      </div>

      <style>{`
        .library-controls {
          margin-bottom: 2rem;
        }
        
        .files-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 1.5rem;
          padding-bottom: 1rem;
          border-bottom: 1px solid var(--border);
        }
        
        .files-count {
          color: var(--text-secondary);
          font-size: 0.9rem;
        }
        
        .file-card {
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
        
        .file-card.unavailable {
          opacity: 0.6;
        }
        
        .file-card.unavailable .file-title {
          color: var(--text-muted);
        }
        
        .file-meta .status {
          font-size: 0.75rem;
        }
        
        .file-meta .status.available {
          color: var(--success);
        }
        
        .file-meta .status.offline {
          color: var(--danger);
        }
      `}</style>
    </section>
  );
};

export default Library;
