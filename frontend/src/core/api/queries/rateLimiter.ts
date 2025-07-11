/**
 * BACKEND PROTECTION SYSTEM - Rate Limiter
 *
 * This file prevents TanstackQuery from hammering the backend with too many requests.
 *
 * HOW IT WORKS:
 * - Keeps track of which API calls are failing
 * - If the same API call fails 3 times in a row, it blocks that call for 30 seconds
 * - When an API call succeeds, it resets the failure count
 * - Works for all users (global tracking)
 *
 * EXAMPLE:
 * - User tries to load profile → fails
 * - User tries again → fails
 * - User tries again → fails
 * - System blocks profile loading for 30 seconds
 * - After 30 seconds, user can try again
 *
 * WHY THIS MATTERS:
 * - Prevents infinite retry loops that could crash the backend
 * - Hopefully protects hosting bill costs from exploding
 * - Gives backend time to recover from issues
 */

type QueryKey = string;

interface RateLimitState {
  consecutiveFailures: number;  // How many times this query has failed in a row
  blockedUntil: number;        // When to stop blocking (timestamp in milliseconds)
}

// CONFIGURATION
const RATE_LIMIT_BLOCK_MS = 30_000; // How long to block failing queries (30 seconds)
const RATE_LIMIT_FAILURES = 3;       // How many failures before blocking

// Global storage - remembers which queries are blocked
const rateLimitMap = new Map<QueryKey, RateLimitState>();

/**
 * Converts a query key array (such as ['user-profile']) into a string for storage
 */
function getKey(queryKey: unknown): QueryKey {
  return JSON.stringify(queryKey);
}

/**
 * Checks if a query is currently blocked due to too many failures
 *
 * @param queryKey - The query to check (e.g. ['user-profile'])
 * @returns true if the query is blocked, false if it's allowed
 */
export function isRateLimited(queryKey: unknown): boolean {
  const key = getKey(queryKey);
  const state = rateLimitMap.get(key);

  // If we've never seen this query, it's not blocked
  if (!state) return false;

  // If the query is blocked and the block time hasn't expired, keep it blocked
  if (state.blockedUntil && Date.now() < state.blockedUntil) return true;

  // Otherwise allow the query
  return false;
}

/**
 * Records that a query failed - Will eventually block it if it fails too often
 *
 * @param queryKey - The query that failed (e.g. ['user-profile'])
 */
export function recordQueryFailure(queryKey: unknown) {
  const key = getKey(queryKey);
  const now = Date.now();

  // Get existing state or create new state
  const state = rateLimitMap.get(key) || { consecutiveFailures: 0, blockedUntil: 0 };

  // Increment failure count
  state.consecutiveFailures += 1;

  // If we've failed too many times, block this query
  if (state.consecutiveFailures >= RATE_LIMIT_FAILURES) {
    state.blockedUntil = now + RATE_LIMIT_BLOCK_MS;
  }

  // Save the updated state
  rateLimitMap.set(key, state);
}

/**
 * Records that a query succeeded - resets the failure count
 *
 * @param queryKey - The query that succeeded (e.g. ['user-profile'])
 */
export function recordQuerySuccess(queryKey: unknown) {
  const key = getKey(queryKey);

  // Reset the failure count and remove any block
  rateLimitMap.set(key, { consecutiveFailures: 0, blockedUntil: 0 });
}