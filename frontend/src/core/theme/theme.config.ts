export const themeConfig = {
  light: {
    // Default light theme values
    background: '0 0% 100%',
    foreground: '222.2 47.4% 11.2%',
    card: '0 0% 100%',
    'card-foreground': '222.2 47.4% 11.2%',
    popover: '0 0% 100%',
    'popover-foreground': '222.2 47.4% 11.2%',
    primary: '222.2 47.4% 11.2%',
    'primary-foreground': '210 40% 98%',
    secondary: '210 40% 96.1%',
    'secondary-foreground': '222.2 47.4% 11.2%',
    muted: '210 40% 96.1%',
    'muted-foreground': '215.4 16.3% 46.9%',
    accent: '210 40% 96.1%',
    'accent-foreground': '222.2 47.4% 11.2%',
    destructive: '0 84.2% 60.2%',
    'destructive-foreground': '210 40% 98%',
    border: '214.3 31.8% 91.4%',
    input: '214.3 31.8% 91.4%',
    ring: '222.2 84% 4.9%',
  },
  dark: {
    // Main backgrounds
    background: '220 20% 12%',    // The deep charcoal background (#161C24)
    foreground: '220 10% 98%',    // Text color on background

    // Cards and elevated surfaces
    card: '220 22% 15%',          // Slightly lighter card backgrounds (#1F262E)
    'card-foreground': '0 0% 100%',

    // Popover/dropdowns
    popover: '220 22% 15%',       // Matching card background
    'popover-foreground': '0 0% 100%',

    // Primary actions
    primary: '210 100% 50%',      // Bright blue (#0072F5) used in charts/buttons
    'primary-foreground': '0 0% 100%',

    // Secondary elements
    secondary: '220 25% 20%',     // Slightly lighter than background
    'secondary-foreground': '220 10% 80%',

    // Muted elements
    muted: '220 20% 18%',         // Used for sidebar and less prominent UI
    'muted-foreground': '220 10% 60%',

    // Accent colors
    accent: '210 100% 50%',       // Same as primary for consistency
    'accent-foreground': '0 0% 100%',

    // Destructive actions
    destructive: '0 85% 60%',     // Red for delete/warning (#FF3B3B)
    'destructive-foreground': '0 0% 100%',

    // Borders and inputs
    border: '220 20% 22%',        // Subtle borders between elements
    input: '220 20% 22%',
    ring: '210 100% 50%',         // Focus rings in primary blue
  },
  system: {
    // System mode uses the same values as light (will be overridden by actual system preference)
    background: '0 0% 100%',
    foreground: '222.2 47.4% 11.2%',
    card: '0 0% 100%',
    'card-foreground': '222.2 47.4% 11.2%',
    popover: '0 0% 100%',
    'popover-foreground': '222.2 47.4% 11.2%',
    primary: '222.2 47.4% 11.2%',
    'primary-foreground': '210 40% 98%',
    secondary: '210 40% 96.1%',
    'secondary-foreground': '222.2 47.4% 11.2%',
    muted: '210 40% 96.1%',
    'muted-foreground': '215.4 16.3% 46.9%',
    accent: '210 40% 96.1%',
    'accent-foreground': '222.2 47.4% 11.2%',
    destructive: '0 84.2% 60.2%',
    'destructive-foreground': '210 40% 98%',
    border: '214.3 31.8% 91.4%',
    input: '214.3 31.8% 91.4%',
    ring: '222.2 84% 4.9%',
  },
} as const;

export type ThemeConfig = typeof themeConfig;
