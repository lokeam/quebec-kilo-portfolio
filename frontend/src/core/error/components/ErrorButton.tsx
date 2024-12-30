import { Button } from '@/shared/components/ui/button/button';
import { cn } from '@/shared/components/ui/utils';

interface ErrorButtonProps {
  onClick: () => void;
  label?: string;
  disabled?: boolean;
  'aria-label'?: string;
  className?: string;
  onKeyDown?: (event: React.KeyboardEvent) => void;
};

export const ErrorButton = ({
  onClick,
  label = 'Try Again',
  disabled = false,
  'aria-label': ariaLabel,
  className = '',
  onKeyDown,
} : ErrorButtonProps) => {
  console.log('ErrorButton');
  return (
    <Button
      variant="default"
      onClick={onClick}
      onKeyDown={onKeyDown}
      disabled={disabled}
      aria-label={ariaLabel}
      className={cn(
        "min-w-[100px] m-4",
        className
      )}
    >
      {label}
    </Button>
  );
};


