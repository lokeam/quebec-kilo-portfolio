// @ts-check
import { defineConfig } from 'astro/config';
import tailwind from '@astrojs/tailwind';

// https://astro.build/config
export default defineConfig({
  integrations: [
    tailwind(),
  ],
  // Disable build caching
  build: {
    inlineStylesheets: 'auto',
  },
  // Disable Vite caching and force reload
  vite: {
    // Force reload on file changes
    server: {
      hmr: {
        overlay: true,
      },
    },
    // Disable caching for immediate development feedback
    optimizeDeps: {
      force: true,
    },
  },
});
