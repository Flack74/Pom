import axios from 'axios';
import { Profile, TimerSession, Suggestion, TodayStats, Plugin } from './types';

const API_BASE = '/api';

export const api = {
  // Profiles
  getProfiles: () => axios.get<{profiles: Profile[], current: string}>(`${API_BASE}/profiles`),
  
  // Session
  startSession: (session: Partial<TimerSession>) => 
    axios.post<{status: string}>(`${API_BASE}/session/start`, session),
  
  // Insights
  getSuggestions: () => axios.get<Suggestion[]>(`${API_BASE}/insights/suggestions`),
  getTodayStats: () => axios.get<TodayStats>(`${API_BASE}/insights/today`),
  
  // Plugins
  getPlugins: () => axios.get<Plugin[]>(`${API_BASE}/plugins`),
  
  // Privacy
  getPrivacyStatus: () => axios.get<{privacy_mode: boolean, cloud_sync: boolean}>(`${API_BASE}/privacy/status`),
};