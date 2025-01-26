import { memo, type ReactNode } from 'react';

interface InfoSectionProps {
  icon: ReactNode;
  label: string;
  value: string;
  hasStackedContent?: boolean;
  isVisible?: boolean;
  isMobile?: boolean;
  isCardView?: boolean;
};

export const InfoSection = memo(({
  icon,
  label,
  value,
  hasStackedContent = false,
  isVisible = true,
  isMobile = false,
  isCardView = false,
}: InfoSectionProps) => {

  if (!isVisible || isMobile) return null;

  return (
    <div className={`flex flex-row items-center gap-2 ${
      hasStackedContent ? 'flex-col max-w-[70px] overflow-x-hidden' : ''
    }`}>
      {icon}
      <div className={`flex flex-col ${isCardView ? 'ml-[5px]' : ''}`}>
        <span className={`mr-2 text-xs uppercase ${hasStackedContent ? 'hidden' : ''}`}>{label}</span>
        <span className={`text-sm text-white ${hasStackedContent ? 'max-w-[70px]' : 'max-w-[105px]'} overflow-x-hidden truncate`}>{value.charAt(0).toUpperCase() + value.slice(1)}</span>
      </div>
    </div>
  );
}, (prevProps, nextProps) => {
  return prevProps.hasStackedContent === nextProps.hasStackedContent &&
    prevProps.isMobile === nextProps.isMobile &&
    prevProps.isVisible === nextProps.isVisible &&
    prevProps.label === nextProps.label &&
    prevProps.value === nextProps.value &&
    prevProps.icon === nextProps.icon;
});
