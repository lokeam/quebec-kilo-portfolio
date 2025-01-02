import React from 'react';
import { cn } from '@/shared/components/ui/utils';

interface PagelineProps extends React.HTMLAttributes<HTMLElement> {
  className?: string;
  children?: React.ReactNode;
};

export function PageHeadline({
  className,
  children,
  ...props
}: PagelineProps) {
  return (
    <div
      className={cn(
        'mb-2 flex items-center justify-between space-y-2', className
      )}
      {...props}
    >
      {children}
    </div>
  );
};
