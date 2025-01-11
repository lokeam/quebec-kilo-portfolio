import { Suspense, lazy } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import ProtectedRoute from '@/core/auth/components/ProtectedRoute/ProtectedRoute';

// Pages
import HomePage from '@/features/dashboard/pages/HomePage/HomePage';
import { Loading } from '@/shared/components/ui/loading/Loading';

// Layouts
import AuthenticatedLayout from '@/core/auth/components/AuthenicatedLayout/AuthenticatedLayout';

// Providers
import { ErrorBoundaryProvider } from '@/core/error/providers/ErrorBoundaryProvider';
import { ThemeProvider } from '@/core/theme/providers/ThemeProvider';
import { SidebarProvider } from './shared/components/ui/sidebar';
import { NetworkStatusProvider } from '@/core/network-status/providers/NetworkStatusProvider';

// TODO: Lazy load secondary routes
const OnlineServicesPage = lazy(() => import(
  /* webpackChunkName: "OnlineServicesPage" */
  '@/features/dashboard/pages/OnlineServices/OnlineServicesPage'
));

const LibraryPage = lazy(() => import(
  /* webpackChunkName: "LibraryPage" */
  '@/features/dashboard/pages/LibraryPage/LibraryPage'
));

const WishListPage = lazy(() => import(
  /* webpackChunkName: "WishListPage" */
  '@/features/dashboard/pages/WishListPage/WishListPage'
));


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
                    <Route path="/library" element={<LibraryPage />}/>
                    <Route path="/online-services" element={<OnlineServicesPage />} />
                    <Route path="/wishlist" element={<WishListPage />} />
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
