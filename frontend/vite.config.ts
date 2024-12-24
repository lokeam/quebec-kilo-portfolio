import { defineConfig, mergeConfig } from 'vitest/config';
import viteConfig from './vite.base.config';

export default mergeConfig(
  viteConfig,
  defineConfig({
    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: ['./test/setup/setupTests.ts'],
      include: ['src/**/*.{test,spec}.{ts,tsx}']
    }
  })
);