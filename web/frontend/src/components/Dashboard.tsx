import React, { useState, useEffect } from 'react';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Box,
  Chip,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
} from '@mui/material';
import { TrendingUp, Today, Psychology, Extension } from '@mui/icons-material';
import { api } from '../api';
import { TodayStats, Suggestion, Plugin } from '../types';

export const Dashboard: React.FC = () => {
  const [todayStats, setTodayStats] = useState<TodayStats | null>(null);
  const [suggestions, setSuggestions] = useState<Suggestion[]>([]);
  const [plugins, setPlugins] = useState<Plugin[]>([]);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [statsRes, suggestionsRes, pluginsRes] = await Promise.all([
        api.getTodayStats(),
        api.getSuggestions(),
        api.getPlugins(),
      ]);
      
      setTodayStats(statsRes.data);
      setSuggestions(suggestionsRes.data);
      setPlugins(pluginsRes.data);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    }
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
        🚀 Mission Control Dashboard
      </Typography>

      <Box display="flex" flexDirection="column" gap={3}>
        {/* Today's Stats */}
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center" mb={2}>
              <Today color="primary" sx={{ mr: 1 }} />
              <Typography variant="h6">Today's Progress</Typography>
            </Box>
            {todayStats && (
              <Box display="flex" gap={4}>
                <Box>
                  <Typography variant="h3" color="primary">
                    {todayStats.sessions}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Sessions completed
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="h3" color="success.main">
                    {Math.round(todayStats.hours * 10) / 10}h
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Focus time
                  </Typography>
                </Box>
              </Box>
            )}
          </CardContent>
        </Card>

        {/* AI Suggestions */}
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center" mb={2}>
              <Psychology color="secondary" sx={{ mr: 1 }} />
              <Typography variant="h6">AI Insights</Typography>
            </Box>
            <List dense>
              {suggestions.slice(0, 3).map((suggestion, index) => (
                <ListItem key={index} sx={{ px: 0 }}>
                  <ListItemIcon>
                    <TrendingUp color="success" />
                  </ListItemIcon>
                  <ListItemText
                    primary={suggestion.message}
                    secondary={`Confidence: ${Math.round(suggestion.confidence * 100)}%`}
                  />
                </ListItem>
              ))}
              {suggestions.length === 0 && (
                <Typography variant="body2" color="text.secondary">
                  Complete a few sessions to get personalized suggestions!
                </Typography>
              )}
            </List>
          </CardContent>
        </Card>

        {/* Active Plugins */}
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center" mb={2}>
              <Extension color="warning" sx={{ mr: 1 }} />
              <Typography variant="h6">Active Plugins</Typography>
            </Box>
            <Box>
              {plugins.filter(p => p.enabled).map((plugin) => (
                <Chip
                  key={plugin.name}
                  label={plugin.name}
                  color="success"
                  variant="outlined"
                  size="small"
                  sx={{ mr: 1, mb: 1 }}
                />
              ))}
              {plugins.filter(p => p.enabled).length === 0 && (
                <Typography variant="body2" color="text.secondary">
                  No plugins enabled
                </Typography>
              )}
            </Box>
          </CardContent>
        </Card>
      </Box>
    </Box>
  );
};