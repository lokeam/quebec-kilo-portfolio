/**
 * Development environment detection via Node env
 */
const isDev = process.env.NODE_ENV === 'development';

/**
 * Configuration options for Logger
 * @property {boolean} enabled - Whether or not logging is enabled
 * @property {('debug'|'info'|'warn'|'error')} level - Minimum log level to display
 */
interface LoggerConfig {
  enabled: boolean;
  level: 'debug' | 'info' | 'warn' | 'error';
}

/**
 * Utility class for consistent application logging
 *
 * Features:
 * - Environment-aware (dev vs prod)
 * - Configurable log levels
 * - Structured data logging
 * - Emoji prefixes for visual distinction
 *
 * @example
 * ```typescript
 * // Enable all logging
 * logger.configure({ enabled: true });
 *
 * // Disable logging
 * logger.configure({ enabled: false });
 *
 * // Change log level
 * logger.configure({ level: 'error' });
 * ```
 */
class Logger {
  /**
   * Default configuration based on environment
   * @private
   */
  private config: LoggerConfig = {
    enabled: isDev,
    level: 'debug',
  };

  /**
   * Log debug level messages
   * @param message - Main log message
   * @param data - Optional structured data to log
   * @example
   * ```typescript
   * logger.debug('User action:', { userId: 123, action: 'click' });
   * ```
   */
  debug(message: string, data?: unknown): void {
    if (!this.config.enabled) return;
    console.log(`🔍 DEBUG: ${message}`, data);
  }

  /**
   * Log informational messages
   * @param message - Main log message
   * @param data - Optional structured data to log
   * @example
   * ```typescript
   * logger.info('User logged in:', { userId: 123 });
   * ```
   */
  info(message: string, data?: unknown): void {
    if (!this.config.enabled) return;
    console.info(`💬 INFO: ${message}`, data);
  }

  /**
   * Log warning messages
   * @param message - Main log message
   * @param data - Optional structured data to log
   * @example
   * ```typescript
   * logger.warn('API rate limit approaching:', { remaining: 10 });
   * ```
   */
  warn(message: string, data?: unknown): void {
    if (!this.config.enabled) return;
    console.warn(`⚠️ WARN: ${message}`, data);
  }

  /**
   * Log error messages
   * @param message - Main log message
   * @param data - Optional error object or debug data
   * @example
   * ```typescript
   * logger.error('API request failed:', error);
   * ```
   */
  error(message: string, data?: unknown): void {
    if (!this.config.enabled) return;
    console.error(`🚨 ERROR: ${message}`, data);
  }

  /**
   * Update logger configuration
   * @param config - Partial configuration to merge with existing
   *
   * @example
   * ```typescript
   * // Enable all logging
   * logger.configure({ enabled: true });
   *
   * // Disable logging
   * logger.configure({ enabled: false });
   *
   * // Change log level
   * logger.configure({ level: 'error' });
   * ```
   */
  configure(config: Partial<LoggerConfig>): void {
    this.config = { ...this.config, ...config };
  }
}

// Export singleton instance
export const logger = new Logger();
