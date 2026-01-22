import { useState, useRef } from 'react';
import api from '../services/api';
import { useToast } from '../context/ToastContext';

const Upload = () => {
  const [file, setFile] = useState(null);
  const [formData, setFormData] = useState({
    title: '',
    subject: 'Computer Science',
    description: ''
  });
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [dragging, setDragging] = useState(false);
  const fileInputRef = useRef(null);
  const { showToast } = useToast();

  const subjects = [
    'Computer Science',
    'Mathematics',
    'Physics',
    'Chemistry',
    'Biology',
    'Engineering',
    'Literature',
    'History',
    'Economics',
    'Psychology',
    'Other'
  ];

  const handleDragOver = (e) => {
    e.preventDefault();
    setDragging(true);
  };

  const handleDragLeave = () => {
    setDragging(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setDragging(false);
    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile) {
      setFile(droppedFile);
      if (!formData.title) {
        setFormData({ ...formData, title: droppedFile.name.replace(/\.[^/.]+$/, '') });
      }
    }
  };

  const handleFileSelect = (e) => {
    const selectedFile = e.target.files[0];
    if (selectedFile) {
      setFile(selectedFile);
      if (!formData.title) {
        setFormData({ ...formData, title: selectedFile.name.replace(/\.[^/.]+$/, '') });
      }
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!file) {
      showToast('Please select a file to upload', 'error');
      return;
    }

    setUploading(true);
    setProgress(0);

    // Simulate upload progress
    const progressInterval = setInterval(() => {
      setProgress(prev => {
        if (prev >= 90) {
          clearInterval(progressInterval);
          return 90;
        }
        return prev + 10;
      });
    }, 200);

    try {
      // In a real app, this would upload the file
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      clearInterval(progressInterval);
      setProgress(100);
      
      showToast('File uploaded successfully! +5 reputation', 'success');
      
      // Reset form
      setTimeout(() => {
        setFile(null);
        setFormData({ title: '', subject: 'Computer Science', description: '' });
        setProgress(0);
        setUploading(false);
      }, 1000);
    } catch (error) {
      clearInterval(progressInterval);
      showToast('Upload failed. Please try again.', 'error');
      setUploading(false);
      setProgress(0);
    }
  };

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
            <span>📤</span> Share Resources
          </h1>
          <p className="page-subtitle">
            Upload academic resources to share with the community and earn reputation
          </p>
        </div>

        <div className="upload-container">
          {/* Benefits Sidebar */}
          <div className="upload-info">
            <h3>🎁 Benefits of Sharing</h3>
            <p className="text-secondary">Contribute to the community and get rewarded!</p>
            
            <ul className="upload-benefits">
              <li>
                <span>⭐</span>
                <div>
                  <strong>+5 Reputation</strong>
                  <p>Earn points for every upload</p>
                </div>
              </li>
              <li>
                <span>🚀</span>
                <div>
                  <strong>Priority Downloads</strong>
                  <p>Higher reputation = faster speeds</p>
                </div>
              </li>
              <li>
                <span>🏆</span>
                <div>
                  <strong>Leaderboard Rankings</strong>
                  <p>Top contributors get featured</p>
                </div>
              </li>
              <li>
                <span>♾️</span>
                <div>
                  <strong>Unlimited Access</strong>
                  <p>No download limits for top sharers</p>
                </div>
              </li>
            </ul>

            <div className="upload-tips">
              <h4>📝 Tips for Quality Uploads</h4>
              <ul>
                <li>Use clear, descriptive titles</li>
                <li>Select the correct subject category</li>
                <li>Add helpful descriptions</li>
                <li>Ensure files are complete and readable</li>
              </ul>
            </div>
          </div>

          {/* Upload Form */}
          <div className="upload-form">
            <form onSubmit={handleSubmit}>
              {/* File Drop Zone */}
              <div 
                className={`file-drop-zone ${dragging ? 'dragging' : ''} ${file ? 'has-file' : ''}`}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
                onClick={() => fileInputRef.current?.click()}
              >
                <input
                  type="file"
                  ref={fileInputRef}
                  onChange={handleFileSelect}
                  accept=".pdf,.doc,.docx,.txt,.ppt,.pptx"
                  style={{ display: 'none' }}
                />
                
                {file ? (
                  <div className="file-preview">
                    <span className="file-icon">📄</span>
                    <div className="file-details">
                      <strong>{file.name}</strong>
                      <span>{formatSize(file.size)}</span>
                    </div>
                    <button 
                      type="button" 
                      className="remove-file"
                      onClick={(e) => {
                        e.stopPropagation();
                        setFile(null);
                      }}
                    >
                      ✕
                    </button>
                  </div>
                ) : (
                  <>
                    <span className="drop-icon">📁</span>
                    <p className="drop-text">
                      Drag & drop your file here, or <span className="browse-link">browse</span>
                    </p>
                    <p className="drop-hint">Supports PDF, DOC, DOCX, PPT, TXT</p>
                  </>
                )}
              </div>

              {/* Form Fields */}
              <div className="form-group">
                <label htmlFor="title">Resource Title</label>
                <input
                  type="text"
                  id="title"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  placeholder="e.g., Advanced Machine Learning Notes"
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="subject">Subject Category</label>
                <select
                  id="subject"
                  value={formData.subject}
                  onChange={(e) => setFormData({ ...formData, subject: e.target.value })}
                >
                  {subjects.map(subject => (
                    <option key={subject} value={subject}>{subject}</option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label htmlFor="description">Description (Optional)</label>
                <textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="Add a brief description of the resource..."
                  rows={3}
                />
              </div>

              {/* Progress Bar */}
              {uploading && (
                <div className="upload-progress">
                  <div className="progress-bar">
                    <div 
                      className="progress-fill" 
                      style={{ width: `${progress}%` }}
                    ></div>
                  </div>
                  <span className="progress-text">{progress}%</span>
                </div>
              )}

              {/* Submit Button */}
              <button 
                type="submit" 
                className="btn btn-primary btn-block btn-lg"
                disabled={!file || uploading}
              >
                {uploading ? (
                  progress === 100 ? '✅ Upload Complete!' : '⏳ Uploading...'
                ) : (
                  '🚀 Upload Resource'
                )}
              </button>
            </form>
          </div>
        </div>
      </div>

      <style>{`
        .upload-info {
          background: var(--bg-card);
          border: 1px solid var(--border);
          border-radius: var(--radius-lg);
          padding: 2rem;
        }
        
        .upload-info h3 {
          font-size: 1.25rem;
          margin-bottom: 0.5rem;
          background: var(--gradient-primary);
          -webkit-background-clip: text;
          background-clip: text;
          color: transparent;
        }
        
        .upload-benefits {
          margin: 1.5rem 0;
        }
        
        .upload-benefits li {
          display: flex;
          align-items: flex-start;
          gap: 1rem;
          padding: 1rem 0;
          border-bottom: 1px solid var(--border);
        }
        
        .upload-benefits li:last-child {
          border-bottom: none;
        }
        
        .upload-benefits li span:first-child {
          font-size: 1.5rem;
          width: 40px;
          text-align: center;
        }
        
        .upload-benefits li strong {
          display: block;
          margin-bottom: 0.25rem;
        }
        
        .upload-benefits li p {
          font-size: 0.875rem;
          color: var(--text-muted);
          margin: 0;
        }
        
        .upload-tips {
          margin-top: 2rem;
          padding: 1.5rem;
          background: rgba(139, 92, 246, 0.1);
          border-radius: var(--radius-md);
          border: 1px solid rgba(139, 92, 246, 0.2);
        }
        
        .upload-tips h4 {
          margin-bottom: 1rem;
          font-size: 1rem;
        }
        
        .upload-tips ul {
          list-style: none;
        }
        
        .upload-tips li {
          padding: 0.5rem 0;
          padding-left: 1.5rem;
          position: relative;
          color: var(--text-secondary);
          font-size: 0.9rem;
        }
        
        .upload-tips li::before {
          content: '✓';
          position: absolute;
          left: 0;
          color: var(--success);
        }
        
        .upload-form {
          background: var(--bg-card);
          border: 1px solid var(--border);
          border-radius: var(--radius-lg);
          padding: 2rem;
        }
        
        .file-drop-zone {
          margin-bottom: 1.5rem;
        }
        
        .file-drop-zone.has-file {
          padding: 1.5rem;
        }
        
        .drop-hint {
          font-size: 0.875rem;
          color: var(--text-muted);
          margin-top: 0.5rem;
        }
        
        .file-preview {
          display: flex;
          align-items: center;
          gap: 1rem;
        }
        
        .file-preview .file-icon {
          font-size: 2.5rem;
        }
        
        .file-preview .file-details {
          flex: 1;
        }
        
        .file-preview .file-details strong {
          display: block;
          word-break: break-all;
        }
        
        .file-preview .file-details span {
          font-size: 0.875rem;
          color: var(--text-muted);
        }
        
        .remove-file {
          width: 32px;
          height: 32px;
          border-radius: 50%;
          border: 1px solid var(--border);
          background: var(--bg-glass);
          color: var(--text-muted);
          cursor: pointer;
          transition: all 0.2s;
        }
        
        .remove-file:hover {
          background: var(--danger);
          border-color: var(--danger);
          color: white;
        }
        
        .upload-progress {
          display: flex;
          align-items: center;
          gap: 1rem;
          margin-bottom: 1.5rem;
        }
        
        .progress-bar {
          flex: 1;
          height: 8px;
          background: var(--bg-glass);
          border-radius: 4px;
          overflow: hidden;
        }
        
        .progress-fill {
          height: 100%;
          background: var(--gradient-primary);
          border-radius: 4px;
          transition: width 0.3s ease;
        }
        
        .progress-text {
          font-size: 0.875rem;
          font-weight: 600;
          min-width: 40px;
          text-align: right;
        }
      `}</style>
    </section>
  );
};

export default Upload;
