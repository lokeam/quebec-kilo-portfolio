import type { ReactElement } from 'react';
import { LOGO_MAP, useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';
import { memo, useMemo } from 'react';
import { IconCloudDataConnection } from '@tabler/icons-react';

/**
 * @file LogoMap.tsx
 * @description A memoized component for rendering SVG logos based on domain and name.
 * Optimized for performance with React.memo and useMemo to prevent unnecessary re-renders.
 */


type DomainName = keyof typeof LOGO_MAP;
export type LogoName<D extends DomainName> = keyof typeof LOGO_MAP[D];

interface SVGLogoProps<D extends DomainName> {
   /** Domain category (e.g., 'games', 'movies', 'music') */
  domain: D;

  /** The specific logo name within the domain */
  name: LogoName<D>;

  /** Optional className for styling */
  className?: string;
}

/**
 * Core SVG Logo component before memoization
 * @template D - Domain type constraint
 * @param {SVGLogoProps<D>} props - Component props
 * @returns {ReactElement | null} - Rendered logo or fallback icon
 */
function SVGLogoComponent<D extends DomainName>({
  domain,
  name,
  className,
}: SVGLogoProps<D>): ReactElement | null {
  const MEMOIZED_DOMAIN_MAPS = useDomainMaps();

  // Memoize entire logo selection + rendering process to prevent unnecessary re-renders
  return useMemo(() => {
    const domainMap = MEMOIZED_DOMAIN_MAPS[domain as keyof typeof MEMOIZED_DOMAIN_MAPS];
    const SelectedLogo = !domainMap || !domainMap[name as keyof typeof domainMap]
      ? IconCloudDataConnection // Fallback to a generic Network services icon if logo not found
      : domainMap[name as keyof typeof domainMap];

    return <SelectedLogo className={className} />;
  }, [domain, name, className, MEMOIZED_DOMAIN_MAPS]);
}

/**
 * Memoized SVG Logo component with strict prop comparison
 * Only re-renders when domain, name, or className changes
 */
const SVGLogo = memo(SVGLogoComponent, (prevProps, nextProps) => {
  return (
    prevProps.domain === nextProps.domain &&
    prevProps.name === nextProps.name &&
    prevProps.className === nextProps.className
  );
}) as typeof SVGLogoComponent;

export default SVGLogo;
