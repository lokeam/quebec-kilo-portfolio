// Custom middleware for proxying API requests
import axios from 'axios';

/**
 * Creates a middleware for the Vite dev server that proxies API requests
 * to the actual API server. This avoids cross-origin issues in development.
 *
 * @returns {import('connect').NextHandleFunction} Connect middleware function
 */

const API_VERSION = '/api/v1/';

export function createApiMiddleware() {
  return async function apiMiddleware(req, res, next) {
    // Only handle API requests
    // UPDATE: we now need to check for /api/v1/ instead of /v1/
    if (req.url && req.url.startsWith(API_VERSION)) {
      try {
        // console.log(`[API Proxy] Forwarding request: ${req.method} ${req.url}`);

        // Forward the request to the API server
        // NOTE: API_BASE_PATH in the req.url is set in frontend/src/core/api/config.ts
        const apiResponse = await axios({
          method: req.method,
          url: `http://api.localhost${req.url}`,
          headers: {
            ...(req.headers && req.headers.authorization && { 'Authorization': req.headers.authorization }),
            'Accept': 'application/json',
            'Content-Type': (req.headers && req.headers['content-type']) || 'application/json'
          },
          data: req.body,
          responseType: 'json',
          timeout: 30000 // 30 seconds
        });

        // console.log(`[API Proxy] Response received: ${apiResponse.status}`);

        // Return the API response
        res.setHeader('Content-Type', 'application/json');
        res.statusCode = apiResponse.status;
        res.end(JSON.stringify(apiResponse.data));
      } catch (error) {
        console.error('[API Proxy] Error:', error.message);

        // Return the error
        res.statusCode = error.response?.status || 500;
        res.setHeader('Content-Type', 'application/json');
        res.end(JSON.stringify({
          error: error.message,
          details: error.response?.data || 'Internal server error'
        }));
      }
    } else {
      // Not an API request, pass to next middleware
      next();
    }
  };
}