import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: [
      { find: '@', replacement: resolve(__dirname, './src') },
      { find: '@test', replacement: resolve(__dirname, './test') }
    ]
  },
  server: {
    host: '0.0.0.0',
    port: 5000, // Use a completely different port
    hmr: {
      overlay: false
    },
    proxy: {
      '/api': {
        target: 'http://api.localhost',
        changeOrigin: true,
        secure: false
      }
    }
  }
});