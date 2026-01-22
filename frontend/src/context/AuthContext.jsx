import { createContext, useContext, useState, useEffect } from 'react';
import api from '../services/api';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Load user on mount if token exists
  useEffect(() => {
    if (token) {
      loadUser();
    } else {
      setIsLoading(false);
    }
  }, [token]);

  const loadUser = async () => {
    try {
      const response = await api.getCurrentUser();
      if (response.data.success) {
        setUser(response.data.data);
        setIsAuthenticated(true);
      }
    } catch (error) {
      console.error('Failed to load user:', error);
      logout();
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (email, password) => {
    try {
      const response = await api.login(email, password);
      if (response.data.success) {
        const { user: userData, token: authToken } = response.data.data;
        setUser(userData);
        setToken(authToken);
        setIsAuthenticated(true);
        localStorage.setItem('token', authToken);
        return { success: true };
      }
    } catch (error) {
      const message = error.response?.data?.error || 'Login failed';
      return { success: false, error: message };
    }
  };

  const signup = async (email, username, password) => {
    try {
      const response = await api.register(email, username, password);
      if (response.data.success) {
        return { success: true, message: response.data.message };
      }
    } catch (error) {
      const message = error.response?.data?.error || 'Registration failed';
      return { success: false, error: message };
    }
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    setIsAuthenticated(false);
    localStorage.removeItem('token');
    api.logout().catch(() => {});
  };

  const isAdmin = () => {
    return user?.role === 'admin';
  };

  const value = {
    user,
    token,
    isAuthenticated,
    isLoading,
    isAdmin,
    login,
    signup,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
