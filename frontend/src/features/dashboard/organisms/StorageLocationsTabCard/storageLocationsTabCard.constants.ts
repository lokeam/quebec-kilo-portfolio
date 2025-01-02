export const ITEMS_PER_PAGE = 5;

export enum TabValue {
  Physical = 'physical',
  Digital = 'digital'
};

export const TAB_LABELS = {
  [TabValue.Physical]: 'Physical Storage',
  [TabValue.Digital]: 'Online Storage'
} as const;
