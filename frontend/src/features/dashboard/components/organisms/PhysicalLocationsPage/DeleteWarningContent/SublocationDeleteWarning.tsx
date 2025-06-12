import { memo } from 'react';
import { truncateWarningText } from '@/features/dashboard/lib/utils/truncateWarningText';
import type { LocationsBFFSublocationResponse } from '@/types/domain/physical-location';

interface SublocationDeleteWarningProps {
  location: LocationsBFFSublocationResponse;
}

export const SublocationDeleteWarning = memo(({ location }: SublocationDeleteWarningProps) => {
  const hasGames = location.storedGames?.length > 0;

  if (!hasGames) {
    return (
      <p>
        Are you sure you want to delete this location? This action cannot be undone.
      </p>
    );
  }

  return (
    <div className="space-y-2">
      <p>
        Are you sure to want to delete this location? You will be deleting the following games stored here:
      </p>
      <ul className="list-disc pl-4 space-y-1">
        {location.storedGames.slice(0, 5).map((game) => (
          <li key={game.id} className="text-sm">
            "{truncateWarningText(game.name)}": ({truncateWarningText(game.platform)})
            {game.isUniqueCopy && " only copy"}
          </li>
        ))}
        {location.storedGames.length > 5 && (
          <li className="text-sm italic">
            Complete game listing for this location available on the media storage page
          </li>
        )}
      </ul>
      <p className="mt-2">
        This action cannot be undone.
      </p>
    </div>
  );
});

SublocationDeleteWarning.displayName = 'SublocationDeleteWarning';
