import { forwardRef } from 'react';
import { Button } from '@/shared/components/ui/button';
import type { ButtonProps } from '@/shared/components/ui/button';
import { cn } from '@/shared/components/ui/utils';

interface StateControlledButtonProps extends ButtonProps {
  isDisabled?: boolean;
};

export const StateControlledButton = forwardRef<HTMLButtonElement, ValidatedButtonProps>(
  ({
    children,
    isDisabled = false,
    className= '',
    ...props
  }, ref) => {
    return (
      <Button
        ref={ref}
        disabled={isDisabled}
        className={cn(className, isDisabled && 'cursor-not-allowed opacity-50')}
        {...props}
      >
        {children}
      </Button>
    )
  }
);

StateControlledButton.displayName = 'StateControlledButton';
