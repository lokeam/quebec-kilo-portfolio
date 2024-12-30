import { useEffect, useState } from 'react';
import { cn } from '@/shared/components/ui/utils';
import { SidebarTrigger } from '@/shared/components/ui/sidebar'
import { Separator } from '@/shared/components/ui/separator';

interface TopNavigationProps extends React.HTMLAttributes<HTMLElement> {
  fixed?: boolean
  ref?: React.Ref<HTMLElement>
};

export function TopNavigation({
  className,
  fixed,
  children,
  ...props
}: TopNavigationProps) {
  const [offset, setOffset] = useState<number>(0);

  useEffect(() => {
    const onScroll = () => {
      setOffset(document.body.scrollTop || document.documentElement.scrollTop);
    };

    document.addEventListener('scroll', onScroll, { passive: true });

    return () => document.removeEventListener('scroll', onScroll);
  }, []);

  return (
    <header
      id="top-nav"
      className={cn(
        'flex items-center gap-3 sm:gap-4 bg-background p-4 h-16',
        fixed && 'header-fixed peer/header w-[inherit] fixed z-50 rounded-md',
        offset > 10 && fixed ? 'shadow' : 'shadow-none',
        className
      )}
      {...props}
    >
      <SidebarTrigger variant='outline' className='scale-125 sm:scale-100' />
      <Separator orientation='vertical' className='h-6' />
      {children}
    </header>
  );
};
