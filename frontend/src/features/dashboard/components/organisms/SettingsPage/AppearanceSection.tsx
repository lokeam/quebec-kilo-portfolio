import { useFormContext } from 'react-hook-form';
import { Monitor, Moon, Sun } from '@/shared/components/ui/icons';
import { Button } from '@/shared/components/ui/button';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/shared/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
} from '@/shared/components/ui/dropdown-menu';

import type { FormValues } from '@/features/dashboard/pages/SettingsPage/SettingsForm';
import { useThemeStore } from '@/core/theme/stores/useThemeStore';


export function AppearanceSection() {
  const form = useFormContext<FormValues>();
  const { watch, setValue } = form;

  const { actions } = useThemeStore();

  return (
    <Card>
      <CardHeader>
        <CardTitle>Appearance</CardTitle>
        <CardDescription>Choose your preferred color scheme.</CardDescription>
      </CardHeader>
      <CardContent className="grid gap-6">
        {/* Theme */}
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">Color Scheme</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Available in light or dark mode.
            </p>
          </div>
          <DropdownMenu modal={false}>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm">
                {watch("theme") === "light" && <Sun className="h-4 w-4 mr-2" />}
                {watch("theme") === "dark" && <Moon className="h-4 w-4 mr-2" />}
                {watch("theme") === "system" && <Monitor className="h-4 w-4 mr-2" />}
                {watch("theme")[0].toUpperCase() + watch("theme").slice(1)}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              align="end"
              avoidCollisions={false}
            >
              <DropdownMenuRadioGroup
                value={watch("theme")}
                onValueChange={(value) => {
                  setValue('theme', value as FormValues['theme']);

                  if (value === 'system') {
                    actions.enableSystemPreference();
                  } else {
                    actions.disableSystemPreference();
                    actions.changeTheme(value as 'light' | 'dark');
                  }
                }}
              >
                <DropdownMenuRadioItem value="light">
                  <Sun className="h-4 w-4 mr-2" />
                  Light
                </DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="dark">
                  <Moon className="h-4 w-4 mr-2" />
                  Dark
                </DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="system">
                  <Monitor className="h-4 w-4 mr-2" />
                  System
                </DropdownMenuRadioItem>
              </DropdownMenuRadioGroup>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardContent>
    </Card>
  )
}
