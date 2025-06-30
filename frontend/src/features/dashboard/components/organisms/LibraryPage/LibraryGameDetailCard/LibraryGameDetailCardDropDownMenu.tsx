import { memo, useState } from 'react';
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import { Button } from '@/shared/components/ui/button';
import { Settings } from 'lucide-react';
import { LibraryGameDeletionDialog } from '@/features/dashboard/components/organisms/LibraryPage/LibraryGameDeletionDialog/LibraryGameDeletionDialog';
import type { LibraryGameItemRefactoredResponse } from '@/types/domain/library-types';

interface LibraryGameDetailCardDropDownMenuProps {
  game: LibraryGameItemRefactoredResponse;
  onRemoveFromLibrary: () => void;
}

function LibraryGameDetailCardDropDownMenu({
  game,
  onRemoveFromLibrary
}: LibraryGameDetailCardDropDownMenuProps) {
  const [removeDialogOpen, setRemoveDialogOpen] = useState(false);

  const handleRemoveClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setRemoveDialogOpen(true);
  };

  const handleDeleteSuccess = () => {
    onRemoveFromLibrary();
    setRemoveDialogOpen(false);
  };

  return (
    <>
      <DropdownMenuPrimitive.Root modal={false}>
        <DropdownMenuPrimitive.Trigger asChild>
          <Button
            variant="ghost"
            size="icon"
          >
            <Settings className="h-4 w-4" />
          </Button>
        </DropdownMenuPrimitive.Trigger>
        <DropdownMenuPrimitive.Portal>
          <DropdownMenuPrimitive.Content
            className="z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"
            sideOffset={4}
            align="end"
            side="top"
            avoidCollisions={false}
          >
            <DropdownMenuPrimitive.Item
              className="relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
              onClick={handleRemoveClick}
            >
              Remove from Library
            </DropdownMenuPrimitive.Item>
          </DropdownMenuPrimitive.Content>
        </DropdownMenuPrimitive.Portal>
      </DropdownMenuPrimitive.Root>

      <LibraryGameDeletionDialog
        game={game}
        open={removeDialogOpen}
        onOpenChange={setRemoveDialogOpen}
        onDeleteSuccess={handleDeleteSuccess}
      />
    </>
  );
}

export const MemoizedLibraryGameDetailCardDropDownMenu = memo(LibraryGameDetailCardDropDownMenu);