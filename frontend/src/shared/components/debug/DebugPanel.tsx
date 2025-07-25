/**
 * Debug Utilities for Development
 *
 * This file provides debugging tools for testing authentication, toasts, and optimistic state.
 * Available in development mode only.
 *
 * USAGE:
 * 1. Console Access: debugUtils.clearToastState()
 * 2. URL Parameters: http://localhost:3000/?debug=toasts
 * 3. Component Access: const debugUtils = useDebugUtils()
 *
 * @example
 * // Console access (pure functions only)
 * debugUtils.clearToastState()
 * debugUtils.testOptimisticIntroToasts()
 *
 * // Component access (all functions including hooks)
 * const debugUtils = useDebugUtils()
 * debugUtils.debugAuth0Claims()
 *
 * // URL parameter triggers
 * http://localhost:3000/?debug=toasts
 * http://localhost:3000/?debug=optimistic
 */

import { useDebugUtils } from "@/shared/utils/debug-utils";

export const DebugPanel = () => {
  const debugUtils = useDebugUtils();

  return (
    <div className="fixed bottom-4 right-4 z-50 bg-white border rounded-lg shadow-lg p-4">
      <h3 className="text-sm font-bold mb-2">Debug Tools</h3>
      <div className="space-y-1">

        {/* Debug Auth0 Claims - Check user claims and metadata */}
        <button
          onClick={debugUtils.debugAuth0Claims}
          className="px-3 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Debug Auth0 Claims
        </button>

        {/* Clear Toast State - Reset localStorage and reload */}
        <button
          onClick={debugUtils.clearToastState}
          className="px-3 py-1 text-sm bg-red-500 text-white rounded hover:bg-red-600"
        >
          Clear Toast State
        </button>

        {/* Test Optimistic Toasts - Set optimistic state to true */}
        <button
          onClick={debugUtils.testOptimisticIntroToasts}
          className="px-3 py-1 text-sm bg-green-500 text-white rounded hover:bg-green-600"
        >
          Test Optimistic Toasts
        </button>

        {/* Reset Optimistic State - Clear all optimistic state */}
        <button
          onClick={debugUtils.resetOptimisticState}
          className="px-3 py-1 text-sm bg-yellow-500 text-white rounded hover:bg-yellow-600"
        >
          Reset Optimistic State
        </button>

        {/* Debug Current State - Log current application state */}
        <button
          onClick={debugUtils.debugCurrentState}
          className="px-3 py-1 text-sm bg-purple-500 text-white rounded hover:bg-purple-600"
        >
          Debug Current State
        </button>

        {/* Test No Toasts - Set optimistic state to false */}
        <button
          onClick={debugUtils.testNoToasts}
          className="px-3 py-1 text-sm bg-gray-500 text-white rounded hover:bg-gray-600"
        >
          Test No Toasts
        </button>
      </div>
    </div>
  );
};