import { type ReactNode } from 'react';
import { useAuth } from '@/core/auth/hooks/useAuth';
import { AuthContext } from '@/core/auth/context-provider/AuthContext';

export function AuthProvider({ children }: { children: ReactNode }) {
  const auth = useAuth();
  return (
    <AuthContext.Provider value={auth}>
      {children}
    </AuthContext.Provider>
  );
}