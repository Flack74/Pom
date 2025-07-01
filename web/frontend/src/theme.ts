import { createTheme } from '@mui/material/styles';

export const galacticTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#18FFFF', // Neon cyan
    },
    secondary: {
      main: '#FF4081', // Vibrant pink
    },
    success: {
      main: '#00E676', // Emerald green
    },
    warning: {
      main: '#FFD600', // Solar yellow
    },
    background: {
      default: '#0B0F1A', // Deep space navy
      paper: '#1A1F2E',
    },
    text: {
      primary: '#FFFFFF', // Crisp white
      secondary: '#B0BEC5', // Soft steel gray
    },
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontWeight: 700,
      background: 'linear-gradient(45deg, #18FFFF, #FF4081)',
      WebkitBackgroundClip: 'text',
      WebkitTextFillColor: 'transparent',
    },
    h2: {
      fontWeight: 600,
      color: '#18FFFF',
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          textTransform: 'none',
          fontWeight: 600,
          boxShadow: '0 4px 20px rgba(24, 255, 255, 0.3)',
          '&:hover': {
            boxShadow: '0 6px 25px rgba(24, 255, 255, 0.4)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          background: 'linear-gradient(135deg, #1A1F2E 0%, #0B0F1A 100%)',
          border: '1px solid rgba(24, 255, 255, 0.1)',
          borderRadius: 16,
          backdropFilter: 'blur(10px)',
        },
      },
    },
  },
});