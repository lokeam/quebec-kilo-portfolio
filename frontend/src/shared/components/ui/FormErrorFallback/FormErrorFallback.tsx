import { type FallbackProps } from 'react-error-boundary';
import { Button } from '@/shared/components/ui/button';

export function FormErrorFallback({ error, resetErrorBoundary }: FallbackProps) {
  return (
    <div role="alert" className="p-4 border border-destructive rounded-lg">
      <h2 className="text-lg font-semibold text-destructive">Something went wrong</h2>
      <p className="text-sm text-muted-foreground mt-2">{error.message}</p>
      <Button
        onClick={resetErrorBoundary}
        className="mt-4"
        variant="default"
      >
        Try Again
      </Button>
    </div>
  );
}