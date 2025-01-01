import { Suspense } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import HomePage from '@/pages/HomePage/HomePage';
import { ErrorBoundaryProvider } from '@/core/error/providers/ErrorBoundaryProvider';
import { ThemeProvider } from '@/core/theme/providers/ThemeProvider';
import AuthenticatedLayout from '@/core/auth/components/AuthenicatedLayout/AuthenticatedLayout';
import ProtectedRoute from '@/core/auth/components/ProtectedRoute/ProtectedRoute';
import { SidebarProvider } from './shared/components/ui/sidebar';
import { Loading } from '@/shared/components/ui/loading/Loading';

// TODO: Lazy load secondary routes

function App() {
  return (
    <ErrorBoundaryProvider>
      <BrowserRouter>
        <ThemeProvider enableSystemPreference>
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
        </ThemeProvider>
      </BrowserRouter>
    </ErrorBoundaryProvider>
  );
}

export default App;
