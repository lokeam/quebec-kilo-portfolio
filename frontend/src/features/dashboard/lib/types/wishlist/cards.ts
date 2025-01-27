/**
 * Controls the visibility state of different sections within a WishlistCard component.
 * This interface manages responsive behavior and conditional rendering based on viewport
 * size and user preferences.
 */
export interface CardVisibility {
  readonly showTags: boolean;         // Controls tag section visibility
  readonly showRating: boolean;       // Controls rating section visibility
  readonly showReleaseDate: boolean;  // Controls release date visibility
  readonly showMoreDeals: boolean;    // Controls additional deals visibility
  readonly stackPriceContent: boolean;// Controls price content layout

  // Optional display configurations
  readonly stackInfoContent?: boolean;      // Controls info content layout
  readonly showLocationInfo?: boolean;      // Controls location info visibility
  readonly showSublocationInfo?: boolean;   // Controls sublocation info visibility
  readonly isMobile?: boolean;              // Controls mobile-specific rendering
}

/**
 * Default visibility configuration for wishlist cards.
 * Used as initial state and for reset operations.
 */
export const DEFAULT_CARD_VISIBILITY: Readonly<CardVisibility> = {
  showTags: true,
  showRating: true,
  showReleaseDate: true,
  showMoreDeals: true,
  stackPriceContent: false,
  stackInfoContent: false,
  showLocationInfo: false,
  showSublocationInfo: false,
  isMobile: false,
};
