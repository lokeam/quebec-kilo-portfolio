/** @type {import("tailwindcss").Config} */
export default {
  content: ["./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}"],
  theme: {
    extend: {
      colors: {
        // Brand Colors - Semantic naming
        brand: {
          primary: {
            50: '#f7fdf4',
            100: '#eefbe0',
            200: '#ddf5c7',
            300: '#c7ed9e',
            400: '#a8e06f',
            500: '#7CC254', // Your base green
            600: '#6bb344',
            700: '#5a9a3a',
            800: '#4a7d30',
            900: '#3d6628',
            950: '#1f3314',
          },
          secondary: {
            50: '#f0f4fc',
            100: '#e1e9f8',
            200: '#c3d3f1',
            300: '#9bb7e8',
            400: '#6b95d8',
            500: '#3D6FB7', // Your base blue
            600: '#365ea3',
            700: '#2f4f8a',
            800: '#284172',
            900: '#22355e',
            950: '#111a2f',
          },
          accent: {
            50: '#f7f4fc',
            100: '#eee7f8',
            200: '#ddd3f1',
            300: '#c4b3e8',
            400: '#a38dd8',
            500: '#674099', // Your base purple
            600: '#5d3888',
            700: '#4e2f72',
            800: '#40275c',
            900: '#352149',
            950: '#1a1024',
          },
        },
        // Semantic Colors - For UI states
        success: {
          50: '#f0fdf4',
          100: '#dcfce7',
          200: '#bbf7d0',
          300: '#86efac',
          400: '#4ade80',
          500: '#22c55e',
          600: '#16a34a',
          700: '#15803d',
          800: '#166534',
          900: '#14532d',
        },
        warning: {
          50: '#fffbeb',
          100: '#fef3c7',
          200: '#fde68a',
          300: '#fcd34d',
          400: '#fbbf24',
          500: '#f59e0b',
          600: '#d97706',
          700: '#b45309',
          800: '#92400e',
          900: '#78350f',
        },
        error: {
          50: '#fef2f2',
          100: '#fee2e2',
          200: '#fecaca',
          300: '#fca5a5',
          400: '#f87171',
          500: '#ef4444',
          600: '#dc2626',
          700: '#b91c1c',
          800: '#991b1b',
          900: '#7f1d1d',
        },
      },
      // Custom spacing for brand consistency
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
      },
      // Custom font sizes for brand typography
      fontSize: {
        'xs': ['0.75rem', { lineHeight: '1rem' }],
        'sm': ['0.875rem', { lineHeight: '1.25rem' }],
        'base': ['1rem', { lineHeight: '1.5rem' }],
        'lg': ['1.125rem', { lineHeight: '1.75rem' }],
        'xl': ['1.25rem', { lineHeight: '1.75rem' }],
        '2xl': ['1.5rem', { lineHeight: '2rem' }],
        '3xl': ['1.875rem', { lineHeight: '2.25rem' }],
        '4xl': ['2.25rem', { lineHeight: '2.5rem' }],
        '5xl': ['3rem', { lineHeight: '1' }],
        '6xl': ['3.75rem', { lineHeight: '1' }],
      },
    },
  },
  plugins: [],
}
