# Icon System

This marketing site uses a simple icon system that supports both predefined icons and custom SVG content.

## Usage

### Using Predefined Icons

```astro
---
import Icon from './ui/Icon.astro';
---

<Icon name="gameController" size={24} class="text-blue-500" />
<Icon name="chartPie" size={32} class="text-green-500" />
<Icon name="users" size={16} class="text-purple-500" />
```

### Using Custom SVG Icons

```astro
---
import Icon from './ui/Icon.astro';
---

<Icon size={24} class="text-red-500">
  <svg xmlns="http://www.w3.org/2000/svg" class="w-full h-full" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
    <path d="M12 5v14"/>
    <path d="M5 12h14"/>
  </svg>
</Icon>
```

## Available Predefined Icons

- `gameController` - Video game controller icon
- `cloudData` - Cloud data connection icon
- `chartPie` - Pie chart icon
- `users` - Users/people icon
- `headset` - Headset/support icon
- `gift` - Gift/early access icon
- `rocket` - Rocket/get started icon
- `messageQuestionMark` - Question mark icon
- `blocks` - Building blocks icon
- `iconCoin` - Coin/money icon
- `chartBar` - Bar chart icon
- `iconHeadset` - Headset icon
- `chevronDown` - Chevron down arrow
- `plus` - Plus sign
- `check` - Checkmark

## Adding New Icons

To add new icons, edit `icons.ts` and add your SVG path to the `icons` object:

```typescript
export const icons = {
  // ... existing icons
  newIcon: `
    <svg xmlns="http://www.w3.org/2000/svg" class="w-full h-full" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
      <!-- Your SVG path here -->
    </svg>
  `
} as const;
```

## Icon Sources

Icons are manually extracted from:
- **Tabler Icons** - Professional icon set with consistent stroke width
- **Lucide Icons** - Clean, modern icon set

SVG paths are copied directly from the official websites to avoid unnecessary dependencies.

## Performance Benefits

- **No React dependencies** - No unnecessary bundle size
- **Tree-shakeable** - Only the icons you use are included
- **Lightweight** - Raw SVG paths are very small
- **Fast loading** - No external library downloads