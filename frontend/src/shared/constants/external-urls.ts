/**
 * External URLs Configuration
 *
 * Uses environment variables for all URLs. No hardcoded fallbacks for security.
 * Environment variables must be set in all environments (dev, staging, prod).
 */

/**
 * Marketing site URLs
 */
export const EXTERNAL_URLS = {
  // Marketing site base URL
  MARKETING_SITE: import.meta.env.VITE_MARKETING_SITE_URL,

  // Legal pages
  TERMS_AND_CONDITIONS: import.meta.env.VITE_TERMS_URL,
  PRIVACY_POLICY: import.meta.env.VITE_PRIVACY_URL,

  // Support pages
  SUPPORT: import.meta.env.VITE_SUPPORT_URL,
  CONTACT: import.meta.env.VITE_CONTACT_URL,
} as const;

/**
 * Type for external URL keys
 */
export type ExternalUrlKey = keyof typeof EXTERNAL_URLS;

/**
 * Helper function to get external URL by key with validation
 */
export const getExternalUrlByKey = (key: ExternalUrlKey): string => {
  const url = EXTERNAL_URLS[key];

  if (!url) {
    console.error(`External URL not configured: ${key}. Please set the corresponding environment variable.`);
    // Return a safe fallback that won't expose sensitive information
    return '#';
  }

  return url;
};

/**
 * Validation function to check if all required URLs are configured
 */
export const validateExternalUrls = (): void => {
  const missingUrls = Object.entries(EXTERNAL_URLS)
    .filter(([, url]) => !url)
    .map(([key]) => key);

  if (missingUrls.length > 0) {
    console.warn('Missing external URL environment variables:', missingUrls);
  }
};