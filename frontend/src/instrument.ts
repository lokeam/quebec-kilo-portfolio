import * as Sentry from "@sentry/react";

Sentry.init({
  dsn: import.meta.env.VITE_SENTRY_DSN,
  environment: import.meta.env.VITE_SENTRY_ENVIRONMENT,

  // Adds request headers and IP for users
  sendDefaultPii: true,

  // Enable logs (optional)
  _experiments: {
    enableLogs: true,
  },

  // Use tunneling to avoid ad blockers
  tunnel: import.meta.env.VITE_SENTRY_TUNNEL,

  // Console logging integration - sends console.log, console.error, and console.warn calls as logs to Sentry
  integrations: [
    Sentry.consoleLoggingIntegration({ levels: ["log", "error", "warn"] }),
  ],
});