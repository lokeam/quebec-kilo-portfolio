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
    port: 3000,
    // Uncomment below to enable polling if file watch events are unreliable.
    // watch: { usePolling: true, interval: 500 },
    hmr: {
      overlay: false // Disable the HMR overlay if you feel it's impacting performance.
    }
  },
  optimizeDeps: {
    // Pre-bundle critical dependencies to reduce the number of requests.
    include: ['react', 'react-dom']
    // You can add additional libraries here.
  },
  // Use a dedicated local cache directory (stored inside the container)
  cacheDir: 'node_modules/.vite-cache'
});