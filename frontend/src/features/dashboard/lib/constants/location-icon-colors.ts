import type { LocationIconBgColor } from '@/types/domain/location-types';

/**
 * Type for the color mapping object that contains both light and dark mode colors
 */
type ColorMapping = {
  light: {
    text: string;
    background: string;
  };
  dark: {
    text: string;
    background: string;
  };
};

/**
 * Mapping of location icon background colors to their respective hex values
 * for both light and dark themes
 */
export const LOCATION_ICON_COLORS: Record<LocationIconBgColor, ColorMapping> = {
  red: {
    light: {
      text: '#C4554D',
      background: '#FAECEC',
    },
    dark: {
      text: '#BE524B',
      background: '#BE524B',
    },
  },
  green: {
    light: {
      text: '#548164',
      background: '#EEF3ED',
    },
    dark: {
      text: '#4F9768',
      background: '#242B26',
    },
  },
  blue: {
    light: {
      text: '#487CA5',
      background: '#E9F3F7',
    },
    dark: {
      text: '#447ACB',
      background: '#447ACB',
    },
  },
  orange: {
    light: {
      text: '#CC782F',
      background: '#F8ECDF',
    },
    dark: {
      text: '#CB7B37',
      background: '#36291F',
    },
  },
  gold: {
    light: {
      text: '#C29343',
      background: '#FAF3DD',
    },
    dark: {
      text: '#C19138',
      background: '#372E20',
    },
  },
  purple: {
    light: {
      text: '#8A67AB',
      background: '#F6F3F8',
    },
    dark: {
      text: '#865DBB',
      background: '#2A2430',
    },
  },
  brown: {
    light: {
      text: '#976D57',
      background: '#F3EEEE',
    },
    dark: {
      text: '#A27763',
      background: '#2E2724',
    },
  },
  gray: {
    light: {
      text: '#787774',
      background: '#F1F1EF',
    },
    dark: {
      text: '#9B9B9B',
      background: '#252525',
    },
  },
  pink: {
    light: {
      text: '#B35488',
      background: '#F9F2F5',
    },
    dark: {
      text: '#BA4A78',
      background: '#2E2328',
    },
  },
};

/**
 * Default colors for when no specific color is provided
 */
export const DEFAULT_COLORS: ColorMapping = {
  light: {
    text: '#373530',
    background: '#FFFFFF',
  },
  dark: {
    text: '#D4D4D4',
    background: '#191919',
  },
};