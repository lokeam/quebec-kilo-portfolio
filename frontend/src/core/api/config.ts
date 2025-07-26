export const QUERY_STALE_TIME = 1000 * 60 * 5; // 5 minutes
export const QUERY_GARBAGE_COLLECTION_TIME = 1000 * 60 * 5; // 5 minutes

// Force development to always use proxy settings
const isDevelopment = import.meta.env.DEV;
export const BASE_API_URL = isDevelopment ? '' : (import.meta.env.VITE_API_URL || 'https://api.q-ko.com');
export const API_BASE_PATH = '/api/v1';