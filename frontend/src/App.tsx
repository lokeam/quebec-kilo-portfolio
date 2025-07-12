import { Suspense, lazy, useEffect } from 'react';
import { BrowserRouter, Route, Routes, useLocation } from 'react-router-dom';
import * as Sentry from '@sentry/react';
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
import { TooltipProvider } from '@/shared/components/ui/tooltip';
import { TanstackMutationToast } from '@/shared/components/ui/TanstackMutationToast/TanstackMutationToast';
import OnboardingProtectedRoute from './core/auth/components/OnboardingProtectedRoute/OnboardingProtectedRoute';
import OnboardingLayout from './core/auth/components/OnboardingLayout/OnboardingLayout';

// Debug Panel (development only)
import { OnboardingDebugPanel } from '@/core/utils/debug/OnboardingDebugPanel';

// Lazy load secondary routes
const OnlineServicesPage = lazy(() => import(
  /* webpackChunkName: "OnlineServicesPage" */
  '@/features/dashboard/pages/OnlineServices/OnlineServicesPage'
));

const PhysicalLocationsPage = lazy(() => import(
  /* webpackChunkName: "PhysicalLocationsPage" */
  '@/features/dashboard/pages/PhysicalLocations/PhysicalLocationsPage'
));

const LibraryPage = lazy(() => import(
  /* webpackChunkName: "LibraryPage" */
  '@/features/dashboard/pages/LibraryPage/LibraryPage'
));

const WishListPage = lazy(() => import(
  /* webpackChunkName: "WishListPage" */
  '@/features/dashboard/pages/WishListPage/WishListPage'
));

const NotificationsPage = lazy(() => import(
  /* webpackChunkName: "NotificationsPage" */
  '@/features/dashboard/pages/NotificationsPage/NotificationsPage'
));

const SettingsPage = lazy(() => import(
  /* webpackChunkName: "SettingsPage" */
  '@/features/dashboard/pages/SettingsPage/SettingsPage'
));

const SpendTrackingPage = lazy(() => import(
  /* webpackChunkName: "SpendTrackingPage" */
  '@/features/dashboard/pages/SpendTrackingPage/SpendTrackingPage'
));

const OnboardingPage = lazy(() => import(
  /* webpackChunkName: "OnboardingPage" */
  '@/features/dashboard/pages/Onboarding/OnboardingPage'
));

const OnboardingName = lazy(() => import(
  /* webpackChunkName: "OnboardingName" */
  '@/features/dashboard/components/organisms/OnboardingPage/OnboardingName'
));

const OnboardingIntro = lazy(() => import(
  /* webpackChunkName: "OnboardingIntro" */
  '@/features/dashboard/components/organisms/OnboardingPage/OnboardingIntro'
));

// Login
const LoginPage = lazy(() => import(
  /* webpackChunkName: "LoginPage" */
  '@/features/login/pages/LoginPage'
));

// Account Recovery
const AccountRecoveryPage = lazy(() => import(
  /* webpackChunkName: "AccountRecoveryPage" */
  '@/features/login/pages/AccountRecoveryPage'
));

// Deleted Account
const DeletedAccountPage = lazy(() => import(
  /* webpackChunkName: "DeletedAccountPage" */
  '@/features/login/pages/DeletedAccountPage'
));

// Deletion Protected Route
const MarkedForDeletionProtectedRoute = lazy(() => import(
  /* webpackChunkName: "MarkedForDeletionProtectedRoute" */
  '@/core/auth/components/MarkedForDeletionProtectedRoute/MarkedForDeletionProtectedRoute'
));

// Route monitoring component
function RouteMonitor() {
  const location = useLocation();

  useEffect(() => {
    // Track route changes in Sentry
    Sentry.addBreadcrumb({
      category: 'navigation',
      message: `Route changed to: ${location.pathname}`,
      level: 'info',
      data: {
        from: location.pathname,
        to: location.pathname,
        search: location.search,
        hash: location.hash,
        timestamp: new Date().toISOString(),
      },
    });

    // Set route context for error reporting
    Sentry.setContext('navigation', {
      currentRoute: location.pathname,
      search: location.search,
      hash: location.hash,
    });
  }, [location.pathname, location.search, location.hash]);

  return null;
}

function App() {
  return (
      <BrowserRouter>
        <ErrorBoundaryProvider>
          <ThemeProvider enableSystemPreference={true}>
            <TooltipProvider delayDuration={300}>
              <NetworkStatusProvider>
                <TanstackMutationToast />
                <RouteMonitor />
                <Suspense fallback={<Loading />}>
                  <Routes>
                    {/* Public routes */}
                    <Route path="/login" element={<LoginPage />} />
                    <Route path="/deleted" element={<DeletedAccountPage />} />

                    {/* Protected routes */}
                    <Route
                      element={
                        <ProtectedRoute>
                          <SidebarProvider defaultOpen={true}>
                            <AuthenticatedLayout />
                          </SidebarProvider>
                        </ProtectedRoute>
                      }
                    >
                      <Route path="/" element={<HomePage />} />
                      <Route path="/library" element={<LibraryPage />}/>
                      <Route path="/online-services" element={<OnlineServicesPage />} />
                      <Route path="/physical-locations" element={<PhysicalLocationsPage />} />
                      <Route path="/wishlist" element={<WishListPage />} />
                      <Route path="/spend-tracking" element={<SpendTrackingPage />} />
                      <Route path="/notifications" element={<NotificationsPage /> } />
                      <Route path="/settings" element={<SettingsPage />} />
                    </Route>

                    {/* Onboarding Routes - New users only */}
                    <Route path="/onboarding/*" element={
                      <OnboardingProtectedRoute>
                        <OnboardingLayout />
                      </OnboardingProtectedRoute>
                    } >
                      <Route path="welcome" element={<OnboardingPage />} />
                      <Route path="name" element={<OnboardingName />} />
                      <Route path="intro" element={<OnboardingIntro />} />
                    </Route>
                    {/* Marked For Deletion Protected Routes */}
                    <Route path="/account-recovery" element={
                      <MarkedForDeletionProtectedRoute>
                        <AccountRecoveryPage />
                      </MarkedForDeletionProtectedRoute>
                    }/>
                  </Routes>
                </Suspense>

                {/* Debug Panel - Development Only */}
                <OnboardingDebugPanel />
              </NetworkStatusProvider>
            </TooltipProvider>
          </ThemeProvider>
        </ErrorBoundaryProvider>
      </BrowserRouter>
  );
}

export default App;
