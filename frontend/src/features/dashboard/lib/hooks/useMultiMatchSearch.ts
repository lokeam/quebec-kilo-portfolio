import { useMemo, useState } from 'react';

interface SearchableItem {
  searchString: string;
}

interface UseMultiMatchSearchResult<T> {
  searchTerm: string;
  setSearchTerm: React.Dispatch<React.SetStateAction<string>>;
  filteredItems: T[];
}

interface MultiMatchSearchConfig {
  minFuzzyLength?: number;
  caseSensitive?: boolean;
}

function multiMatchSearch(
  searchString: string,
  term: string,
  config: MultiMatchSearchConfig = {}
): boolean {
  if (!term || !searchString) return false;

  const normalizedSearch = searchString.toLowerCase();
  const normalizedTerm = term.toLowerCase().trim();

  // Simple substring match
  return normalizedSearch.includes(normalizedTerm);
}

export function useMultiMatchSearch<T extends { searchString: string }>(
  items: T[],
  config: MultiMatchSearchConfig = {}
): UseMultiMatchSearchResult<T> {
  const [searchTerm, setSearchTerm] = useState('');

  const filteredItems = useMemo(() => {
    if (!searchTerm) return items;

    // Split search term into words
    const words = searchTerm.toLowerCase().trim().split(/\s+/).filter(Boolean);
    console.log('Search words:', words);

    return items.filter(item => {
      // All words must match
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