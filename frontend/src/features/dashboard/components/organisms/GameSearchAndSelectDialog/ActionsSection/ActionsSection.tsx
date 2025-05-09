import { useCallback } from 'react';

// Shadcn Components
import { Button } from '@/shared/components/ui/button';

interface ActionsSectionProps {
  selectedGames: Set<string>;
  onConfirm: () => void;
}

export function ActionsSection({ selectedGames, onConfirm }: ActionsSectionProps) {
  const handleConfirm = useCallback(() => {
    onConfirm();
  }, [onConfirm]);

  if (selectedGames.size === 0) return null;

  return (
    <div className="flex justify-end gap-2 pt-4 border-t">
      <Button variant="outline" onClick={handleConfirm}>
        Cancel
      </Button>
      <Button onClick={handleConfirm}>
        Add {selectedGames.size} {selectedGames.size === 1 ? 'Game' : 'Games'}
      </Button>
    </div>
  );
}