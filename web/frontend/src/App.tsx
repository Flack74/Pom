import React, { useState } from 'react';
import { ThemeProvider, CssBaseline, Container, AppBar, Toolbar, Typography, Tabs, Tab, Box } from '@mui/material';
import { galacticTheme } from './theme';
import { Timer } from './components/Timer';
import { Dashboard } from './components/Dashboard';
import { TimerSession } from './types';
import { api } from './api';

function App() {
  const [currentTab, setCurrentTab] = useState(0);
  const [session, setSession] = useState<TimerSession | null>(null);

  const handleStartSession = async (sessionData: Partial<TimerSession>) => {
    try {
      await api.startSession(sessionData);
      // In a real app, you'd get the session back from the API
      setSession({
        id: 'session_1',
        work_time: sessionData.work_time || 25,
        break_time: sessionData.break_time || 5,
        sessions: sessionData.sessions || 4,
        current_session: 1,
        is_running: true,
        is_paused: false,
        is_break: false,
        time_left: (sessionData.work_time || 25) * 60,
        profile: sessionData.profile || 'default',
      });
    } catch (error) {
      console.error('Failed to start session:', error);
    }
  };

  return (
    <ThemeProvider theme={galacticTheme}>
      <CssBaseline />
      <Box sx={{ 
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #0B0F1A 0%, #1A1F2E 50%, #0B0F1A 100%)',
      }}>
        <AppBar position="static" sx={{ background: 'rgba(26, 31, 46, 0.9)', backdropFilter: 'blur(10px)' }}>
          <Toolbar>
            <Typography variant="h6" sx={{ flexGrow: 1, fontWeight: 700 }}>
              🚀 Galactic Pomodoro
            </Typography>
            <Tabs value={currentTab} onChange={(_, v) => setCurrentTab(v)} textColor="inherit">
              <Tab label="Timer" />
              <Tab label="Dashboard" />
            </Tabs>
          </Toolbar>
        </AppBar>
        
        <Container maxWidth="lg" sx={{ py: 4 }}>
          {currentTab === 0 && (
            <Timer session={session} onStart={handleStartSession} />
          )}
          {currentTab === 1 && <Dashboard />}
        </Container>
      </Box>
    </ThemeProvider>
  );
}

export default App;
