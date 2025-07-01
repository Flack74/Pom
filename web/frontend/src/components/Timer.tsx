import React, { useState } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Button,
  LinearProgress,
  Box,
  Chip,
  Stack,
} from '@mui/material';
import { PlayArrow, Pause, Stop } from '@mui/icons-material';
import { TimerSession } from '../types';

interface TimerProps {
  session: TimerSession | null;
  onStart: (session: Partial<TimerSession>) => void;
}

export const Timer: React.FC<TimerProps> = ({ session, onStart }) => {
  const [workTime, setWorkTime] = useState(25);
  const [breakTime, setBreakTime] = useState(5);
  const [sessions, setSessions] = useState(4);

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const getProgress = () => {
    if (!session) return 0;
    const totalTime = session.is_break ? session.break_time * 60 : session.work_time * 60;
    return ((totalTime - session.time_left) / totalTime) * 100;
  };

  return (
    <Card sx={{ maxWidth: 500, mx: 'auto', mt: 4 }}>
      <CardContent sx={{ p: 4 }}>
        <Typography variant="h4" align="center" gutterBottom>
          🍅 Pomodoro Timer
        </Typography>

        {session ? (
          <Box>
            <Stack direction="row" spacing={1} justifyContent="center" mb={2}>
              <Chip
                label={session.is_break ? '☕ Break Time' : '🎯 Focus Time'}
                color={session.is_break ? 'warning' : 'primary'}
              />
              <Chip
                label={`Session ${session.current_session}/${session.sessions}`}
                variant="outlined"
              />
            </Stack>

            <Typography variant="h2" align="center" sx={{ mb: 2, fontFamily: 'monospace' }}>
              {formatTime(session.time_left)}
            </Typography>

            <LinearProgress
              variant="determinate"
              value={getProgress()}
              sx={{
                height: 8,
                borderRadius: 4,
                mb: 3,
                '& .MuiLinearProgress-bar': {
                  background: session.is_break
                    ? 'linear-gradient(90deg, #FFD600, #FF4081)'
                    : 'linear-gradient(90deg, #18FFFF, #00E676)',
                },
              }}
            />

            <Stack direction="row" spacing={2} justifyContent="center">
              <Button variant="contained" startIcon={<Pause />} color="warning">
                Pause
              </Button>
              <Button variant="outlined" startIcon={<Stop />} color="error">
                Stop
              </Button>
            </Stack>
          </Box>
        ) : (
          <Box>
            <Stack spacing={2} mb={3}>
              <Box>
                <Typography variant="body2" color="text.secondary">
                  Work Time: {workTime} minutes
                </Typography>
                <input
                  type="range"
                  min="5"
                  max="60"
                  value={workTime}
                  onChange={(e) => setWorkTime(Number(e.target.value))}
                  style={{ width: '100%', accentColor: '#18FFFF' }}
                />
              </Box>
              <Box>
                <Typography variant="body2" color="text.secondary">
                  Break Time: {breakTime} minutes
                </Typography>
                <input
                  type="range"
                  min="1"
                  max="30"
                  value={breakTime}
                  onChange={(e) => setBreakTime(Number(e.target.value))}
                  style={{ width: '100%', accentColor: '#FF4081' }}
                />
              </Box>
              <Box>
                <Typography variant="body2" color="text.secondary">
                  Sessions: {sessions}
                </Typography>
                <input
                  type="range"
                  min="1"
                  max="8"
                  value={sessions}
                  onChange={(e) => setSessions(Number(e.target.value))}
                  style={{ width: '100%', accentColor: '#00E676' }}
                />
              </Box>
            </Stack>

            <Button
              variant="contained"
              size="large"
              fullWidth
              startIcon={<PlayArrow />}
              onClick={() => onStart({ work_time: workTime, break_time: breakTime, sessions })}
              sx={{ py: 2 }}
            >
              Start Focus Session
            </Button>
          </Box>
        )}
      </CardContent>
    </Card>
  );
};