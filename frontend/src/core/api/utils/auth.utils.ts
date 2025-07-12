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
 * 1. App.tsx calls saveAuth0Token(getAccessTokenSilently)
 * 2. Other code calls getAuth0Token() to get tokens
 *
 * @example
 * // 1. Initialize in App.tsx
 * function App() {
 *   const { getAccessTokenSilently } = useAuth0();
 *   useEffect(() => {
 *     saveAuth0Token(() => getAccessTokenSilently());
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
 * Stores the function that grabs tokens from Auth0.
 * Empty until we set it up by calling saveAuth0Token
 */
let getTokenFn: (() => Promise<string>) | null = null;

/**
 * Saves the function that gets tokens from Auth0
 *
 * Why we need this:
 * - Auth0's function only works in React components
 * - We need to use it in axios interceptors (not React components)
 * - We need to be able to pass options to the function
 *
 * Call this once when the app starts
 *
 * @param fn - The function that gets tokens from Auth0
 */
export function saveAuth0Token(fn: () => Promise<string>) {
  getTokenFn = fn;
}

/**
 * Gets an Auth0 access token from anywhere in the app.
 *
 * @throws Error if saveAuth0Token hasn't been called
 * @returns Promise<string> A valid Auth0 access token
 */
export async function getAuth0Token() {
  if (!getTokenFn) throw new Error('Auth not initialized');
  return await getTokenFn();
}