/**
 * Global token management for Auth0 authentication.
 *
 * This module solves a key architectural challenge:
 * - Auth0's getAccessTokenSilently is a React hook
 * - Hooks can only be used in React components
 * - But we need tokens in non-React code (such as axios interceptors)
 *
 * Solution:
 * 1. Store Auth0's token getter globally
 * 2. Initialize it once at app startup
 * 3. Make it available everywhere
 *
 * Flow:
 * 1. App.tsx calls setAuth0TokenFn(getAccessTokenSilently)
 * 2. Other code calls getAuth0Token() to get tokens
 *
 * @example
 * // 1. Initialize in App.tsx
 * function App() {
 *   const { getAccessTokenSilently } = useAuth0();
 *   useEffect(() => {
 *     setAuth0TokenFn(() => getAccessTokenSilently());
 *   }, [getAccessTokenSilently]);
 * }
 *
 * // 2. Use in any non-React code
 * async function someApiCall() {
 *   const token = await getAuth0Token();
 *   // Use token...
 * }
 */

/**
 * Global storage for Auth0's token getter function.
 * Null until initialized by setAuth0TokenFn.
 */
let getTokenFn: (() => Promise<string>) | null = null;

/**
 * Stores Auth0's getAccessTokenSilently function globally.
 * Must be called once at app initialization.
 *
 * @param fn - Auth0's getAccessTokenSilently function
 */
export function setAuth0TokenFn(fn: () => Promise<string>) {
  getTokenFn = fn;
}

/**
 * Gets an Auth0 access token from anywhere in the app.
 *
 * @throws Error if setAuth0TokenFn hasn't been called
 * @returns Promise<string> A valid Auth0 access token
 */
export async function getAuth0Token() {
  if (!getTokenFn) throw new Error('Auth not initialized');
  return await getTokenFn();
}