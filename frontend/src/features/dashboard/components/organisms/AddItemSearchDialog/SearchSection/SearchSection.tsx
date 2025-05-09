import { useCallback } from 'react';

// Shadcn Components
import { Input } from '@/shared/components/ui/input';
import { Button } from '@/shared/components/ui/button';

// Icons
import { SearchIcon } from 'lucide-react';

interface SearchSectionProps {
  onSearch: (term: string) => void;
}

export function SearchSection({ onSearch }: SearchSectionProps) {
  const handleSearchChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    onSearch(event.target.value);
  }, [onSearch]);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex gap-2">
        <div className="relative flex-1">
          <SearchIcon className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search games..."
            className="pl-8"
            onChange={handleSearchChange}
          />
        </div>
        <Button variant="outline" size="icon">
          <SearchIcon className="h-4 w-4" />
        </Button>
      </div>

      {/* TODO: Add platform and genre filters */}
      <div className="flex gap-2">
        {/* Platform filter */}
        <div className="flex-1">
          {/* TODO: Implement platform filter */}
        </div>

        {/* Genre filter */}
        <div className="flex-1">
          {/* TODO: Implement genre filter */}
        </div>
      </div>
    </div>
  );
}