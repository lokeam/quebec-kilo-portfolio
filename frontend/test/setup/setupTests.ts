import '@testing-library/jest-dom';
import { cleanup } from '@testing-library/react';
import { afterEach, beforeEach, vi } from 'vitest';

// Define type for mocked matchMedia fn
type MockMatchMedia = {
  (query: string): MediaQueryList;
  mockClear: () => void;
};

// Automatically clean up after each test
afterEach(() => {
  cleanup();
  vi.clearAllMocks();

  // Specifically clean up matchMedia mock
  const matchMediaMock = window.matchMedia as MockMatchMedia;
  if (matchMediaMock?.mockClear) {
    matchMediaMock.mockClear();
  }

  // Remove the matchMedia property
  Object.defineProperty(window, 'matchMedia', {
    value: undefined,
    writable: true,
  })
});

beforeEach(() => {
  // Create global mock obj for window.matchMedia
  const matchMediaMock = vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })) as MockMatchMedia;

  // Mock clear method
  matchMediaMock.mockClear = vi.fn();

  // Define the mock as a property of window
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: matchMediaMock,
  })
})

