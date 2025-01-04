import React from 'react';
import { cn } from '@/shared/components/ui/utils';

interface PageGridProps extends React.HTMLAttributes<HTMLElement> {
  className?: string;
  children?: React.ReactNode;
  // Add new prop with default value
  columns?: {
    sm?: number;
    md?: number;
    lg?: number;
  };
};

export function PageGrid({
  className,
  children,
  columns = { sm: 2, lg: 4 }, // Default to original values
  ...props
}: PageGridProps) {
  // Construct grid classes dynamically
  const gridColumns = cn(
    'grid gap-4',
    columns.sm && `sm:grid-cols-${columns.sm}`,
    columns.md && `md:grid-cols-${columns.md}`,
    columns.lg && `lg:grid-cols-${columns.lg}`,
    className
  );

  return (
    <div className={gridColumns} {...props}>
      {children}
    </div>
  );
}
