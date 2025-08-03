import { Suspense, lazy, useEffect } from 'react';
import { BrowserRouter, Route, Routes, useLocation } from 'react-router-dom';
import * as Sentry from '@sentry/react';
import ProtectedRoute from '@/core/auth/components/ProtectedRoute/ProtectedRoute';

// Preloading hooks
import { usePreloadNavigation } from '@/shared/hooks/usePreloadNavigation';

// Pages
import HomePage from '@/features/dashboard/pages/HomePage/HomePage';
import { LoadingPage } from '@/shared/components/ui/loading/LoadingPage';


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
// import { OnboardingDebugPanel } from '@/core/utils/debug/OnboardingDebugPanel';

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

const OnboardingToastSetup = lazy(() => import(
  /* webpackChunkName: "OnboardingIntro" */
  '@/features/dashboard/components/organisms/OnboardingPage/OnboardingToastSetup'
));

// Login
const LoginPage = lazy(() => import(
  /* webpackChunkName: "LoginPage" */
  '@/features/login/pages/LoginPage'
));

// Signup
const SignupPage = lazy(() => import(
  /* webpackChunkName: "SignupPage" */
  '@/features/login/pages/SignupPage'
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

// Deleted Account Protected Route
const DeletedAccountProtectedRoute = lazy(() => import(
  /* webpackChunkName: "DeletedAccountProtectedRoute" */
  '@/core/auth/components/DeletedAccountProtectedRoute/DeletedAccountProtectedRoute'
));

// Marked for Deletion Protected Route
const MarkedForDeletionProtectedRoute = lazy(() => import(
  /* webpackChunkName: "MarkedForDeletionProtectedRoute" */
  '@/core/auth/components/MarkedForDeletionProtectedRoute/MarkedForDeletionProtectedRoute'
));

// Error Page
const ErrorPage = lazy(() => import(
  /* webpackChunkName: "ErrorPage" */
  '@/core/error/pages/ErrorPage'
));

// Route monitoring component
function RouteMonitor() {
  const location = useLocation();

  // Initialize preloading hooks
  usePreloadNavigation();

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
                <Suspense fallback={<LoadingPage />}>
                  <Routes>
                    {/* Public routes */}
                    <Route path="/login" element={<LoginPage />} />
                    <Route path="/signup" element={<SignupPage />} />

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
                      <Route path="/spend-tracking" element={<SpendTrackingPage />} />
                      <Route path="/settings" element={<SettingsPage />} />
                      <Route path="/error" element={<ErrorPage />} />
                    </Route>

                    {/* Onboarding Routes - New users only */}
                    <Route path="/onboarding/*" element={
                      <OnboardingProtectedRoute>
                        <OnboardingLayout />
                      </OnboardingProtectedRoute>
                    } >
                      <Route path="welcome" element={<OnboardingPage />} />
                      <Route path="name" element={<OnboardingName />} />
                      <Route path="messages" element={<OnboardingToastSetup />} />
                    </Route>
                    {/* Marked For Deletion Protected Routes */}
                    <Route path="/account-recovery" element={
                      <MarkedForDeletionProtectedRoute>
                        <AccountRecoveryPage />
                      </MarkedForDeletionProtectedRoute>
                    }/>

                    {/* Deleted Account Protected Routes */}
                    <Route path="/deleted" element={
                      <DeletedAccountProtectedRoute>
                        <DeletedAccountPage />
                      </DeletedAccountProtectedRoute>
                    }/>

                  </Routes>
                </Suspense>
              </NetworkStatusProvider>
            </TooltipProvider>
          </ThemeProvider>
        </ErrorBoundaryProvider>
      </BrowserRouter>
  );
}

export default App;
