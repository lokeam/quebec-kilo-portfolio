export const mockThemeConfig = {
  light: {
    background: '0 0% 100%',
    foreground: '222.2 47.4% 11.2%',
    primary: '222.2 47.4% 11.2%',
    'primary-foreground': '210 40% 98%',
  },
  dark: {
    background: '222.2 84% 4.9%',
    foreground: '210 40% 98%',
    primary: '210 40% 98%',
    'primary-foreground': '222.2 47.4% 11.2%',
  },
} as const;

// Helper to check if theme variables are applied correctly
export const getThemeVariables = () => {
  const root = document.documentElement;
  return {
    background: root.style.getPropertyValue('--background'),
    foreground: root.style.getPropertyValue('--foreground'),
    primary: root.style.getPropertyValue('--primary'),
    'primary-foreground': root.style.getPropertyValue('--primary-foreground'),
  };
};