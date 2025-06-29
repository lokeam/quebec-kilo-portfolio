import { useState, type ImgHTMLAttributes, type ReactNode } from 'react';
import { Image } from 'lucide-react';

interface ImageWithFallbackProps extends Omit<ImgHTMLAttributes<HTMLImageElement>, 'src'> {
  fallback?: ReactNode;
  className?: string;
  width?: number;
  height?: number;
  alt?: string;
  src?: string | undefined | null;
  onLoad?: () => void;
  /** IGDB size token for optimization (e.g., "cover_big", "cover_small", "thumb") */
  igdbSize?: string;
}

export function ImageWithFallback({
  src,
  alt,
  fallback = <Image className="h-6 w-6 text-muted-foreground" />,
  className,
  onLoad,
  igdbSize,
  ...props
}: ImageWithFallbackProps) {
  const [error, setError] = useState(false);

  // Optimize IGDB URLs for better image quality
  const optimizedSrc = src ? optimizeIGDBUrl(src, igdbSize) : undefined;

  if (!src || error) {
    return (
      <div className={`bg-muted flex items-center justify-center ${className}`}>
        {fallback}
      </div>
    );
  }

  return (
    <img
      src={optimizedSrc}
      alt={alt}
      className={className}
      onError={() => setError(true)}
      onLoad={onLoad}
      {...props}
    />
  );
}

/**
 * Optimizes IGDB URLs for better image quality
 * @param url - The original IGDB URL
 * @param size - The desired IGDB size token
 * @returns Optimized URL with better quality
 */
function optimizeIGDBUrl(url: string, size?: string): string {
  // Check if this is an IGDB URL
  if (url.includes('images.igdb.com')) {
    // Add https: prefix if missing (IGDB omits protocol)
    const urlWithProtocol = url.startsWith('//') ? `https:${url}` : url;

    // If size is specified, replace the size token
    if (size) {
      return urlWithProtocol.replace(/t_[^/]+/, `t_${size}`);
    }

    // Default to cover_big for better quality if no size specified
    return urlWithProtocol.replace(/t_[^/]+/, 't_cover_big');
  }

  // Return original URL if not an IGDB URL
  return url;
}
