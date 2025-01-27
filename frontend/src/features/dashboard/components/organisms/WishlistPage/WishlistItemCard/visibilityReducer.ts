import type { CardVisibility } from '@/features/dashboard/lib/types/wishlist/cards';

/**
 * Action type for visibility state changes
 * SET_VISIBILITY: Updates all visibility flags at once based on breakpoint changes
 */
type VisibilityAction = {
  type: 'SET_VISIBILITY';
  payload: CardVisibility;
};

/**
 * Reducer to manage WishlistItemCard's visibility state
 *
 * Performance benefits:
 * 1. Batch updates - All visibility changes happen in a single state update
 * 2. Predictable updates - State changes only happen through a reducer dispatch
 * 3. Stable dispatch reference - useReducer's dispatch function never changes
 * 4. Memoization friendly - Ideally this makes it easier to memoize callbacks that use dispatch
 *
 * @param state - Current visibility state
 * @param action - Action describing the state change
 * @returns New visibility state
 *
 * @example
 * const [visibility, dispatch] = useReducer(visibilityReducer, initialState);
 * // Update all visibility flags at once
 * dispatch({
 *   type: 'SET_VISIBILITY',
 *   payload: { showTags: false, showRating: false, ... }
 * });
 */
export const visibilityReducer = (
  state: CardVisibility,
  action: VisibilityAction
): CardVisibility => {
  switch (action.type) {
    case 'SET_VISIBILITY':
      return action.payload; // Single update for all visibility changes
    default:
      return state;
  }
};
