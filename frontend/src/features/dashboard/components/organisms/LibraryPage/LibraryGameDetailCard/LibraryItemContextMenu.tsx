import type { ReactNode } from 'react';

// ShadCN Components
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from '@/shared/components/ui/context-menu';

// Icons
import { Trash2 } from '@/shared/components/ui/icons';

interface LibraryItemContextMenuProps {
  children: ReactNode;
  onRemoveFromLibrary: () => void;
}

export function LibraryItemContextMenu({ children, onRemoveFromLibrary }: LibraryItemContextMenuProps) {
  return (
    <ContextMenu>
      <ContextMenuTrigger asChild>
        {children}
      </ContextMenuTrigger>
      <ContextMenuContent className="w-64">
        <ContextMenuItem
          onClick={onRemoveFromLibrary}
          className="flex items-center gap-2 text-destructive focus:text-destructive"
        >
          <Trash2 className="h-4 w-4" />
          <span>Remove game from library</span>
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  );
}