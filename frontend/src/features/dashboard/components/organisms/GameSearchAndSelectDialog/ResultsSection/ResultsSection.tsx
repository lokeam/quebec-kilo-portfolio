// Components
import { SearchResult } from '../SearchResult';

// Types
import type { SearchResult as SearchResultType } from '@/types/domain/search';

interface ResultsSectionProps {
  results?: SearchResultType[];
  isLoading: boolean;
  error: Error | null;
  onSelect: (gameId: string) => void;
}

export function ResultsSection({ results, isLoading, error, onSelect }: ResultsSectionProps) {
  if (isLoading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 3 }).map((_, index) => (
          <div key={index} className="h-24 animate-pulse bg-muted rounded-lg" />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-red-500 p-4">
        Error loading results: {error.message}
      </div>
    );
  }

  if (!results) return null;

  return (
    <div className="space-y-4">
      {results.map((result) => (
        <SearchResult
          key={result.game.id}
          game={result.game}
          onAction={() => onSelect(String(result.game.id))}
        />
      ))}
    </div>
  );
}