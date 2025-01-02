import React from 'react';
import { cn } from '@/shared/components/ui/utils';

interface PageGridProps extends React.HTMLAttributes<HTMLElement> {
  className?: string;
  children?: React.ReactNode;
};

export function PageGrid({
  className,
  children,
  ...props
}: PageGridProps) {
  return (
    <div
      className={cn(
        'grid gap-4 sm:grid-cols-2 lg:grid-cols-4', className
      )}
      {...props}
    >
      {children}
    </div>
  );
};
