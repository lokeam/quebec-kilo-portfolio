import { type FallbackProps } from 'react-error-boundary';
import { Button } from '@/shared/components/ui/button/Button';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';

export function HomePageErrorFallback({ error, resetErrorBoundary}: FallbackProps) {
  console.log('HomePageErrorFallback - error:', error);

  return (
    <PageMain>
      <PageHeadline>
        <h1 className="text-2xl font-bold tracking-tight">Dashboard Error</h1>
      </PageHeadline>
      <div className="flex flex-col items-center justify-center gap-4 py-8">
        <p className="text-muted-foreground">
          There was an error loading the dashboard. Please try again.
        </p>
        <Button onClick={resetErrorBoundary}>Retry Dashboard</Button>
      </div>
    </PageMain>
  )
}
