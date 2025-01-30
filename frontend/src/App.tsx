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
import { TooltipProvider } from '@/shared/components/ui/tooltip';

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

const MediStoragePage = lazy(() => import(
  /* webpackChunkName: "MediaStoragepage" */
  '@/features/dashboard/pages/MediaStoragePage/MediaStoragePage'
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

const OnboardingIntro = lazy(() => import(
  /* webpackChunkName: "OnboardingIntro" */
  '@/features/dashboard/components/organisms/OnboardingPage/OnboardingIntro'
));

const OnboardingPageComplete = lazy(() => import(
  /* webpackChunkName: "OnboardingPageComplete" */
  '@/features/dashboard/components/organisms/OnboardingPage/OnboardingPageComplete'
));

const OnboardingLocationSelection = lazy(() => import(
  /* webpackChunkName: "OnboardingLocationSelection" */
  '@/features/dashboard/components/organisms/OnboardingPage/OnboardingLocationSelection'
));

const OnboardingPagePhysicalLocations = lazy(() => import(
  /* webpackChunkName: "OnboardingPagePhysicalLocations" */
  '@/features/dashboard/components/organisms/OnboardingPage/OnboardingPagePhysicalLocations'
));

const OnboardingPagePhysicalSublocations = lazy(() => import(
  /* webpackChunkName: "OnboardingPagePhysicalSublocations" */
  '@/features/dashboard/components/organisms/OnboardingPage/OnboardingPagePhysicalSublocations'
));

const OnboardingPageDigital = lazy(() => import(
  /* webpackChunkName: "OnboardingPageDigital" */
  '@/features/dashboard/components/organisms/OnboardingPage/OnboardingDigital'
));


function App() {
  return (
    <ErrorBoundaryProvider>
      <BrowserRouter>
        <ThemeProvider enableSystemPreference>
          <TooltipProvider delayDuration={300}>
            <NetworkStatusProvider>
              <SidebarProvider defaultOpen={true}>
                <Suspense fallback={<Loading />}>
                  <Routes>
                    <Route path="/onboarding/welcome" element={<OnboardingPage />} />
                    <Route path="/onboarding/intro" element={<OnboardingIntro />} />
                    <Route path="/onboarding/locations" element={<OnboardingLocationSelection />} />
                    <Route path="/onboarding/physical/location" element={<OnboardingPagePhysicalLocations />} />
                    <Route path="/onboarding/physical/sublocation" element={<OnboardingPagePhysicalSublocations />} />
                    <Route path="/onboarding/digital" element={<OnboardingPageDigital />} />
                    <Route path="/onboarding/complete" element={<OnboardingPageComplete />} />

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
                      <Route path="/media-storage" element={<MediStoragePage />} />
                      <Route path="/spend-tracking" element={<SpendTrackingPage />} />
                      <Route path="/notifications" element={<NotificationsPage /> } />
                      <Route path="/settings" element={<SettingsPage />} />
                    </Route>
                  </Routes>
                </Suspense>
              </SidebarProvider>
            </NetworkStatusProvider>
          </TooltipProvider>
        </ThemeProvider>
      </BrowserRouter>
    </ErrorBoundaryProvider>
  );
}

export default App;
