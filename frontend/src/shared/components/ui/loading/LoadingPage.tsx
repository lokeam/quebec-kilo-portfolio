import QKOLoader from '@/shared/components/ui/LogoMap/brand/qko_loader';

export function LoadingPage() {
  return (
    <div
      aria-hidden="true"
      className="flex items-center justify-center min-h-screen bg-background text-foreground loading-container"
      role="status"
      style={{
        // Additional inline styles as backup
        backgroundColor: 'hsl(var(--background))',
        color: 'hsl(var(--foreground))',
      }}
    >
      <QKOLoader />
    </div>
  );
}
