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
    // Leave any development folders out from build
    build: {
      rollupOptions: {
        external: (id) => {
          // Exclude any path containing dev patterns
          const devPatterns = [
            /-dev\//,
            /dev-\//,
            /experimental\//,
            /staging\//,
            /test-\//,
            /temp-\//,
            /draft-\//,
            /wip-\//,
          ];
          return devPatterns.some(pattern => pattern.test(id));
        },
      },
    },
  },
});
