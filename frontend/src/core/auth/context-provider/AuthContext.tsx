import { createContext, useContext } from 'react';
import type { UseAuthReturn } from '@/core/auth/hooks/useAuth'; // Adjust import if needed

export const AuthContext = createContext<UseAuthReturn | undefined>(undefined);

export function useAuthContext() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuthContext must be used within AuthProvider');
  return ctx;
}