/**
 * apiRequest wraps an arbitrary API call with:
 *  - async/await syntax
 *  - pre‑call debug log
 *  - post‑call success log
 *  - catch block with error log + optional Sentry/metrics
 *  - retry logic (if configured)
 *
 * Usage:
 *   return apiRequest('getUserById', () => axios.get(...));
 */

import { logger } from "@/core/utils/logger/logger";

/**
 * @param operation – human‑readable label for logs & metrics
 * @param fn – function that returns your API promise
 */
export async function apiRequest<T>(operation: string, fn: () => Promise<T>): Promise<T> {
  logger.debug(`→ ${operation}`);
  try {
    const result = await fn();
    logger.debug(`← ${operation} succeeded`);
    return result;
  } catch (err) {
    logger.error(`✖ ${operation} failed`, { error: err });
    throw err;
  }
}
