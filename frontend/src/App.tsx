import { Suspense } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import HomePage from '@/features/dashboard/pages/HomePage/HomePage';
import { ErrorBoundaryProvider } from '@/core/error/providers/ErrorBoundaryProvider';
import { ThemeProvider } from '@/core/theme/providers/ThemeProvider';
import AuthenticatedLayout from '@/core/auth/components/AuthenicatedLayout/AuthenticatedLayout';
import ProtectedRoute from '@/core/auth/components/ProtectedRoute/ProtectedRoute';
import { SidebarProvider } from './shared/components/ui/sidebar';
import { Loading } from '@/shared/components/ui/loading/Loading';
import { NetworkStatusProvider } from '@/core/network-status/providers/NetworkStatusProvider';

// TODO: Lazy load secondary routes

function App() {
  return (
    <ErrorBoundaryProvider>
      <BrowserRouter>
        <ThemeProvider enableSystemPreference>
          <NetworkStatusProvider>
            <SidebarProvider defaultOpen={true}>
              <Suspense fallback={<Loading />}>
                <Routes>
                  <Route
                    element={
                      <ProtectedRoute>
                        <AuthenticatedLayout />
                      </ProtectedRoute>
                    }
                  >
                    <Route path="/" element={<HomePage />} />
                  </Route>
                </Routes>
              </Suspense>
            </SidebarProvider>
          </NetworkStatusProvider>
        </ThemeProvider>
      </BrowserRouter>
    </ErrorBoundaryProvider>
  );
}

export default App;
