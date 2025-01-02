import React from 'react';
import { cn } from '@/shared/components/ui/utils';

interface PageMainProps extends React.HTMLAttributes<HTMLElement> {
  fixed?: boolean
  ref?: React.Ref<HTMLElement>
}

export function PageMain({
  fixed,
  ...props
}: PageMainProps) {
  return (
    <main
      className={cn(
        'peer-[.header-fixed]/header:mt-16',
        'px-4 py-6',
        fixed && 'fixed-main flex flex-col flex-grow overflow-hidden'
      )}
      {...props}
    />
  );
}

