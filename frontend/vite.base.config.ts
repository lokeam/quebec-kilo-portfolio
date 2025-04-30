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
    host: '127.0.0.1', // Using 127.0.0.1 instead of 0.0.0.0 for better reliability
    port: 3000,
    strictPort: true, // Added this to fail if port is in use
    hmr: {
      protocol: 'ws',
      host: '127.0.0.1',
      port: 3000,
      overlay: false,
      clientPort: 3000, // Ensure consistent port for HMR
      timeout: 30000 // Increase timeout to 30 seconds
    },
    watch: {
      usePolling: false, // Set to true if you're having file watching issues
      interval: 100
    },
    proxy: {
      '/api': {
        target: 'http://localhost:80',
        changeOrigin: true,
        secure: false,
        headers: {
          'Host': 'api.localhost'
        },
        configure: (proxy) => {
          proxy.on('error', (err) => {
            console.log('proxy error', err);
          });
        }
      }
    }
  },
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      'react-router-dom',
      '@tanstack/react-query',
      '@auth0/auth0-react'
    ],
    force: true // Force dependency optimization
  },
  build: {
    sourcemap: true, // Add sourcemaps for better debugging
    chunkSizeWarningLimit: 1000 // Increase the limit for larger chunks
  },
  cacheDir: 'node_modules/.vite-cache'
});
