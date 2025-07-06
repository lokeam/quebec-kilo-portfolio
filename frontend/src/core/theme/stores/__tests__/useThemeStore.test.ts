import { describe, it, expect, beforeEach } from 'vitest'
import { renderHook, act } from '@testing-library/react';
import {
  useThemeStore,
  useThemeMode,
  useThemeActions,
  useIsSystemPreference
} from '../useThemeStore';

describe('useThemeStore', () => {
   beforeEach(() => {
    // Clear store before each test
    const { actions } = useThemeStore.getState();
    actions.changeTheme('light');
    actions.disableSystemPreference();

    // Clear localStorage
    localStorage.clear();
  });

  describe('initial state', () => {
    it('should initialize with default theme mode', () => {
      const { result } = renderHook(() => useThemeMode());
      expect(result.current).toBe('light');
    });

    it('should initialize with system preference disabled', () => {
      const { result } = renderHook(() => useIsSystemPreference());
      expect(result.current).toBe(false);
    });

    it('should enable system preference when explicitly set', () => {
      const { result } = renderHook(() => ({
        actions: useThemeActions(),
        isSystemPreference: useIsSystemPreference(),
      }));

      act(() => {
        result.current.actions.enableSystemPreference();
      });

      expect(result.current.isSystemPreference).toBe(true);
    });
  });

  describe('theme actions', () => {
    it('should set a specific theme', () => {
      const { result } = renderHook(() => ({
        mode: useThemeMode(),
        actions: useThemeActions(),
      }));

      act(() => {
        result.current.actions.changeTheme('dark');
      });

      expect(result.current.mode).toBe('dark');
    });
  });

  describe('system preference', () => {
    it('should enable system preference', () => {
      const { result } = renderHook(() => useThemeStore())

      act(() => {
        result.current.actions.enableSystemPreference();
      });

      expect(result.current.isSystemPreference).toBe(true);
    })

    it('should disable system preference', () => {
      const { result } = renderHook(() => useThemeStore())

      act(() => {
        result.current.actions.enableSystemPreference();
        result.current.actions.disableSystemPreference();
      });

      expect(result.current.isSystemPreference).toBe(false);
    })
  })
})