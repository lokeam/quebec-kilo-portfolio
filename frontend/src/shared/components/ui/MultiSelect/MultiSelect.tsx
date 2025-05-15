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
  platformNames: string[];
  value?: string[];
  onChange?: (value: string[]) => void;
  onBlur?: () => void;
  name?: string;
}

export const MultiSelect = forwardRef<HTMLButtonElement, MultiSelectProps>(({
  mainPlaceholder,
  secondaryPlaceholder,
  platformNames,
  value: controlledValue,
  onChange,
  onBlur,
  name,
}, ref) => {
  const [open, setOpen] = useState(false);
  const [internalValue, setInternalValue] = useState<string[]>([]);

  // Use controlled value if provided, otherwise use internal state
  const value = controlledValue ?? internalValue;

  const handleSetValue = (mSelectValue: string) => {
    const newValue = value.includes(mSelectValue)
      ? value.filter((item) => item !== mSelectValue)
      : [...value, mSelectValue];

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
                  value.map((val, i) => (
                      <div
                        key={i}
                        className="px-2 py-1 rounded-md border bg-slate-600 text-xs font-medium">
                          {/* {platformNames?.find((platform: string) => platform === val)?.label} */}
                          {val}
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
            <CommandEmpty>No framework found.</CommandEmpty>
            <CommandGroup>
              <CommandList>
                {platformNames?.map((platform: string) => (
                  <CommandItem
                    key={platform}
                    value={platform}
                    onSelect={() => {
                        handleSetValue(platform)
                    }}
                  >
                    <Check
                      className={cn(
                        "mr-2 h-4 w-4",
                        value.includes(platform) ? "opacity-100" : "opacity-0"
                      )}
                    />
                    {platform}
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
