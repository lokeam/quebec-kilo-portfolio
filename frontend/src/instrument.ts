import * as Sentry from "@sentry/react";

// Only initialize Sentry if enabled
if (import.meta.env.VITE_SENTRY_ENABLED !== 'false') {
  Sentry.init({
    dsn: import.meta.env.VITE_SENTRY_DSN,
    environment: import.meta.env.VITE_SENTRY_ENVIRONMENT,

    // Adds request headers and IP for users
    sendDefaultPii: true,

    // Enable performance monitoring
    tracesSampleRate: 1.0, // Track 100% of transactions in development
    profilesSampleRate: 1.0, // Track 100% of profiles in development

      // Performance thresholds
    maxBreadcrumbs: 50, // Limit breadcrumbs to avoid spam

    // Core Web Vitals thresholds (standard metrics)
    // LCP: Largest Contentful Paint should be < 2.5s
    // FID: First Input Delay should be < 100ms
    // CLS: Cumulative Layout Shift should be < 0.1

    // Enable logs and Core Web Vitals monitoring
    _experiments: {
      enableLogs: true,
      // Enable Core Web Vitals tracking
      enableWebVitals: true,
    },

    // Use tunneling to avoid ad blockers
    tunnel: import.meta.env.VITE_SENTRY_TUNNEL_URL,

      // Integrations for comprehensive monitoring
    integrations: [
      // Console logging integration - sends console.log, console.error, and console.warn calls as logs to Sentry
      Sentry.consoleLoggingIntegration({ levels: ["log", "error", "warn"] }),

      // Browser performance monitoring
      Sentry.browserTracingIntegration(),
    ],

    // Filter out noise in development
    beforeSend(event, hint) {
      // Don't send errors in dev unless explicitly testing
      if (import.meta.env.DEV && hint.originalException && typeof hint.originalException === 'object' && 'message' in hint.originalException) {
        const message = String(hint.originalException.message);
        if (!message.includes('Sentry Test')) {
          return null;
        }
      }
      return event;
    },
  });
}