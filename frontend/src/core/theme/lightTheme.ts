import { ThemeOptions } from '@mui/material';
import { commonTheme } from '@/core/theme/commonTheme';

export const lightTheme: ThemeOptions = {
  ...commonTheme,
  palette: {
    mode: 'light',
    primary: {
      main: '#2563eb',
      light: '#60a5fa',
      dark: '#1e40af',
    },
    secondary: {
      main: '#4f46e5',
      light: '#818cf8',
      dark: '#3730a3',
    },
    error: {
      main: '#dc2626',
      light: '#ef4444',
      dark: '#991b1b',
    },
    warning: {
      main: '#d97706',
      light: '#f59e0b',
      dark: '#92400e',
    },
    info: {
      main: '#0284c7',
      light: '#0ea5e9',
      dark: '#075985',
    },
    success: {
      main: '#16a34a',
      light: '#22c55e',
      dark: '#15803d',
    },
    background: {
      default: '#ffffff',
      paper: '#f8fafc',
    },
    text: {
      primary: '#0f172a',
      secondary: '#475569',
    },
  },
};
