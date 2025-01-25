import { memo } from 'react';
import { Badge } from '@/shared/components/ui/badge';

interface DashboardBadgeProps {
  variant?: 'default' | 'secondary' | 'outline';
  className?: string;
  children: React.ReactNode;
}

export const MemoizedDashboardBadge = memo(function DashboardBadge({
  variant = 'secondary',
  className,
  children
}: DashboardBadgeProps) {
  return (
    <Badge variant={variant} className={className}>
      {children}
    </Badge>
  );
});
