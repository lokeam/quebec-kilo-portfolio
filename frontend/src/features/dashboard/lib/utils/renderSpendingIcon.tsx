import { MediaIcon } from './getMediaIcon';
import { DigitalLocationIcon } from './getDigitalLocationIcon';
import { MediaCategory } from '@/types/domain/spend-tracking';
import type { SpendTrackingItemResponse } from '@/types/domain/spend-tracking';

interface RenderSpendingIconProps {
  item: SpendTrackingItemResponse;
  size?: 'sm' | 'md' | 'lg';  // For different size variants
  showLogs?: boolean;         // Optional logging
}

/**
 * Renders the appropriate icon for a spending item based on its type.
 * For subscriptions, uses DigitalLocationIcon with the provider name.
 * For other media types, uses MediaIcon with the media type.
 *
 * @param props.item - The spending item to render an icon for
 * @param props.size - The size variant of the icon (sm: 24px, md: 40px, lg: 56px)
 * @param props.showLogs - Whether to show debug logs
 * @returns The rendered icon component
 */
export const renderSpendingIcon = ({
  item,
  size = 'sm',  // Default to small size
  showLogs = false
}: RenderSpendingIconProps) => {
  // Size mapping
  const sizeClasses = {
    sm: 'h-6 w-6',
    md: 'h-10 w-10',
    lg: 'h-14 w-14'
  };

  // For subscriptions, use the digital location icon
  if (item.mediaType === MediaCategory.SUBSCRIPTION) {
    if (showLogs) {
      console.log(`renderIcon, subscription,: `, item);
    }

    return (
      <DigitalLocationIcon
        name={item.provider}
        className={sizeClasses[size]}
      />
    );
  }

  // For other media types, use the media icon
  return (
    <MediaIcon
      mediaType={item.mediaType}
      className={sizeClasses[size]}
    />
  );
};