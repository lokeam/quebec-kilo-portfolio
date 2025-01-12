import { useState, type ImgHTMLAttributes, type ReactNode } from 'react';
import { Image } from 'lucide-react';

interface ImageWithFallbackProps extends Omit<ImgHTMLAttributes<HTMLImageElement>, 'src'> {
  fallback?: ReactNode;
  className?: string;
  width?: number;
  height?: number;
  alt?: string;
  src?: string | undefined | null;
}


export function ImageWithFallback({
  src,
  alt,
  fallback = <Image className="h-6 w-6 text-muted-foreground" />,
  className,
  ...props
}: ImageWithFallbackProps) {
  const [error, setError] = useState(false);

  if (!src || error) {
    return (
      <div className={`bg-muted flex items-center justify-center ${className}`}>
        {fallback}
      </div>
    );
  }

  return (
    <img
      src={src}
      alt={alt}
      className={className}
      onError={() => setError(true)}
      {...props}
    />
  );
}
