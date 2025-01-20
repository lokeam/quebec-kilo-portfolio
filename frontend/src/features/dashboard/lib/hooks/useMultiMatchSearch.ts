import { useMemo, useState } from 'react';

/**
 * Type definition for the hook's return value
 * @template T - The type of items being filtered
 */
interface UseMultiMatchSearchResult<T> {
  /** Current search term */
  searchTerm: string;
  /** Function to update the search term */
  setSearchTerm: React.Dispatch<React.SetStateAction<string>>;
  /** Array of items that match the search criteria */
  filteredItems: T[];
}

/**
 * Configuration options for the search functionality
 */
interface MultiMatchSearchConfig {
  /** Minimum length required for fuzzy search (not currently implemented) */
  minFuzzyLength?: number;
  /** Whether to perform case-sensitive search (defaults to false) */
  caseSensitive?: boolean;
}


/**
 * Helper function to perform the actual search matching
 *
 * @param searchString - The string to search within
 * @param term - The term to search for
 * @param config - Search configuration options
 * @returns boolean indicating if the term was found
 *
 * @internal
 */
function multiMatchSearch(
  searchString: string,
  term: string,
  config: MultiMatchSearchConfig = {}
): boolean {
  if (!term || !searchString) return false;

  const normalizedSearch = searchString.toLowerCase();
  const normalizedTerm = term.toLowerCase().trim();

  /* Simple substring match */
  return normalizedSearch.includes(normalizedTerm);
}

/**
 * Hook for performing multi-word search across an array of items
 *
 * This hook provides search functionality where:
 * - Multiple words in the search term are treated as AND conditions
 * - Each item must contain all words to be considered a match
 * - Search is case-insensitive by default
 *
 * @template T - Type of items to search through. Must include a searchString property
 *
 * @param items - Array of items to search through
 * @param config - Configuration options for the search
 *
 * @returns {UseMultiMatchSearchResult<T>} Object containing:
 *   - searchTerm: Current search term
 *   - setSearchTerm: Function to update the search term
 *   - filteredItems: Array of items that match the search criteria
 *
 * @example
 * ```tsx
 * interface SearchableItem {
 *   id: string;
 *   searchString: string;
 * }
 *
 * const items: SearchableItem[] = [
 *   { id: '1', searchString: 'playstation gaming console' },
 *   { id: '2', searchString: 'xbox gaming system' }
 * ];
 *
 * const { searchTerm, setSearchTerm, filteredItems } = useMultiMatchSearch(items, {
 *   caseSensitive: false
 * });
 *
 * // Searching for "gaming console" will return the playstation item
 * ```
 */
export function useMultiMatchSearch<T extends { searchString: string }>(
  items: T[],
  config: MultiMatchSearchConfig = {}
): UseMultiMatchSearchResult<T> {
  /* State for search term */
  const [searchTerm, setSearchTerm] = useState('');

  /* Filter items based on search term */
  const filteredItems = useMemo(() => {
    if (!searchTerm) return items;

    /* Split search term into words and scrub empty strings */
    const words = searchTerm.toLowerCase().trim().split(/\s+/).filter(Boolean);
    console.log('Search words:', words);

    return items.filter(item => {
      /* All words must match for search item to be included */
      const matches = words.every(word => multiMatchSearch(item.searchString, word, config));
      console.log(`Item "${item.searchString}" ${matches ? 'matches' : 'does not match'} search terms`);
      return matches;
    });
  }, [items, searchTerm, config]);

  return {
    searchTerm,
    setSearchTerm,
    filteredItems,
  };
}
