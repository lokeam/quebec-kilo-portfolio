import { BrowserRouter, Route, Routes } from 'react-router-dom';
import HomePage from '@/pages/HomePage/HomePage';
import { ErrorBoundaryProvider } from '@/core/error/providers/ErrorBoundaryProvider';
import { ThemeProvider } from '@/core/theme/providers/ThemeProvider';
import AuthenticatedLayout from '@/core/auth/components/AuthenicatedLayout/AuthenticatedLayout';
import ProtectedRoute from '@/core/auth/components/ProtectedRoute/ProtectedRoute';
import { SidebarProvider } from './shared/components/ui/sidebar';

function App() {
  return (
    <ErrorBoundaryProvider
      config={{
        severity: 'error',
        message: 'Application Error',
        actionLabel: 'Reload',
      }}
    >
      <ThemeProvider enableSystemPreference>
        <SidebarProvider defaultOpen={true}>
          <BrowserRouter>
            <Routes>
              <Route element={<ProtectedRoute>
                <AuthenticatedLayout />
                </ProtectedRoute>}>
                <Route path="/" element={<HomePage />} />
              </Route>
            </Routes>
          </BrowserRouter>
        </SidebarProvider>
      </ThemeProvider>
    </ErrorBoundaryProvider>
   );
}

export default App;
