/* Online Services Page */
export const ITEMS_PER_PAGE = 5;

/* Library Page */
export const LIBRARY_MEDIA_ITEM_VISIBILITY = {
  showTags: true,
  showRating: true,
  showReleaseDate: true,
  showMoreDeals: true,
  stackPriceContent: false,
  stackInfoContent: false,
  isMobile: false,
} as const;

export const LIBRARY_MEDIA_ITEM_BREAKPOINT_RULES = [
  {
    breakpoint: 750,
    value: {
      ...LIBRARY_MEDIA_ITEM_VISIBILITY,
      stackInfoContent: true,
    },
  },
  {
    breakpoint: 590,
    value: {
      ...LIBRARY_MEDIA_ITEM_VISIBILITY,
      isMobile: true,
    },
  },

];
