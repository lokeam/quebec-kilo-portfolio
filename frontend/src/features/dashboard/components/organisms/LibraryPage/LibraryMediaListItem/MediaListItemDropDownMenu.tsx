import { memo } from 'react';

// ShadCN Components
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu';
import { Button } from '@/shared/components/ui/button';

// Icons
import { Settings } from 'lucide-react';

interface MediaListItemDropDownMenuProps {
  onRemoveFromLibrary: () => void;
}

function MediaListItemDropDownMenu({ onRemoveFromLibrary }: MediaListItemDropDownMenuProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon">
          <Settings className="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={onRemoveFromLibrary}>
          Remove game from library
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

export const MemoizedMediaListItemDropDownMenu = memo(MediaListItemDropDownMenu);
