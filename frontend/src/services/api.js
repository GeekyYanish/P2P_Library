import axios from 'axios';

const API_BASE_URL = '/api';

// Create axios instance
const apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add token to requests if available
apiClient.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// API Service
const api = {
    // Authentication
    register: (email, username, password) =>
        apiClient.post('/auth/register', { email, username, password }),

    login: (email, password) =>
        apiClient.post('/auth/login', { email, password }),

    logout: () =>
        apiClient.post('/auth/logout'),

    getCurrentUser: () =>
        apiClient.get('/auth/me'),

    // Health & Stats
    checkHealth: () =>
        apiClient.get('/health'),

    getStats: () =>
        apiClient.get('/stats'),

    // Peers
    registerPeer: (peerData) =>
        apiClient.post('/peers/register', peerData),

    getPeers: () =>
        apiClient.get('/peers'),

    getOnlinePeers: () =>
        apiClient.get('/peers/online'),

    // Files
    getFiles: () =>
        apiClient.get('/files'),

    searchFiles: (query) =>
        apiClient.get(`/files/search?q=${encodeURIComponent(query)}`),

    uploadFile: (formData) =>
        apiClient.post('/files/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
        }),

    downloadFile: (cid, requesterId) =>
        apiClient.get(`/files/download?cid=${cid}&requester_id=${requesterId}`),

    // Reputation
    getReputation: (peerId) =>
        apiClient.get(`/reputation?peer_id=${peerId}`),

    getReputationHistory: (peerId) =>
        apiClient.get(`/reputation/history?peer_id=${peerId}`),

    getTopContributors: () =>
        apiClient.get('/reputation/top'),

    // Ratings
    rateFile: (raterId, fileCid, score, comment) =>
        apiClient.post('/ratings/file', { rater_id: raterId, file_cid: fileCid, score, comment }),

    ratePeer: (raterId, targetId, score, comment) =>
        apiClient.post('/ratings/peer', { rater_id: raterId, target_id: targetId, score, comment }),

    getRatings: (targetId, type) =>
        apiClient.get(`/ratings?target_id=${targetId}&type=${type}`),
};

export default api;
