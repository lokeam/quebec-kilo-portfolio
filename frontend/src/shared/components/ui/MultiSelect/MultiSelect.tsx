import { useState, forwardRef } from "react";

// shad cn ui components
import { Button } from "@/shared/components/ui/button";

import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/shared/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/shared/components/ui/popover";

// icons
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/shared/components/ui/utils";

// Types
//import type { Game } from "@/types/domain/game";

// Test data:
// const frameworks = [
//   {
//       value: "next.js",
//       label: "Next.js",
//   },
//   {
//       value: "sveltekit",
//       label: "SvelteKit",
//   },
//   {
//       value: "nuxt.js",
//       label: "Nuxt.js",
//   },
//   {
//       value: "remix",
//       label: "Remix",
//   },
//   {
//       value: "astro",
//       label: "Astro",
//   },
// ]

interface MultiSelectProps {
  mainPlaceholder?: string;
  secondaryPlaceholder?: string;

  platformNames?: string[];
  platforms?: Array<{ id: number; name: string }>;
  value?: Array<{id: number, name: string}>;
  //value?: string[];
  //onChange?: (value: string[]) => void;
  onChange?: (value: Array<{id: number, name: string}>) => void;
  onBlur?: () => void;
  name?: string;
}

export const MultiSelect = forwardRef<HTMLButtonElement, MultiSelectProps>(({
  mainPlaceholder,
  secondaryPlaceholder,
  platforms,
  value: controlledValue,
  onChange,
  onBlur,
  name,
}, ref) => {
  const [open, setOpen] = useState(false);
  //const [internalValue, setInternalValue] = useState<string[]>([]);
  const [internalValue, setInternalValue] = useState<Array<{id: number, name: string}>>([]);

  // Use controlled value if provided, otherwise use internal state
  const value = controlledValue ?? internalValue;

  const handleSetValue = (platform: {id: number, name: string}) => {
    const newValue = value.some(v => v.id === platform.id)
      ? value.filter((item) => item.id !== platform.id)
      : [...value, platform];

    if (onChange) {
      onChange(newValue);
    } else {
      setInternalValue(newValue);
    }
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
          <Button
            ref={ref}
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="w-full justify-between"
            onBlur={onBlur}
            name={name}
          >
            <div className="flex gap-2 justify-start">
              {value?.length ?
                  value.map((platform, i) => (
                    <div
                      key={i}
                      className="px-2 py-1 rounded-md border bg-slate-600 text-xs font-medium">
                      {platform.name}
                    </div>
                  ))
                  : mainPlaceholder || `Which game version are you adding?`}
            </div>
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full p-0">
        <Command>
            <CommandInput placeholder={secondaryPlaceholder || "Available platforms"} />
            <CommandEmpty>No platforms found.</CommandEmpty>
            <CommandGroup>
              <CommandList>
                {platforms?.map((platform) => (
                  <CommandItem
                    key={platform.id}
                    value={platform.name}
                    onSelect={() => {
                        handleSetValue(platform)
                    }}
                  >
                    <Check
                      className={cn(
                        "mr-2 h-4 w-4",
                        value.some(v => v.id === platform.id) ? "opacity-100" : "opacity-0"  // Instead of value.includes(platform)
                      )}
                    />
                    {platform.name}
                  </CommandItem>
                ))}
              </CommandList>
            </CommandGroup>
        </Command>
      </PopoverContent>
    </Popover>
  )
});

MultiSelect.displayName = "MultiSelect";
