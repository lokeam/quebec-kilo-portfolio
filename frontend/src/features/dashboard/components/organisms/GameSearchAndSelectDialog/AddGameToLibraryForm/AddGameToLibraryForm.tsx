// Components
import { FormContainer } from '@/features/dashboard/components/templates/FormContainer';

// Shadcn UI Components
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/components/ui/form';

import * as RadioGroup from "@radix-ui/react-radio-group";
import { ErrorBoundary } from 'react-error-boundary';
import { FormErrorFallback } from '@/shared/components/ui/FormErrorFallback/FormErrorFallback';

// Components
import { MultiSelect } from '@/shared/components/ui/MultiSelect/MultiSelect';

// Hooks
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';

import { cn } from '@/shared/components/ui/utils';

// Component specific types and interfaces
type AddToLibraryFormStorageType = 'physical' | 'digital' | 'both';

interface AddToLibraryFormStorageOption {
  value: AddToLibraryFormStorageType;
  label: string;
  icon: React.ReactNode;
}

// Zod
import { z } from 'zod';

// Types
import type { Game } from '@/types/domain/game';
import type { PhysicalLocation } from '@/types/domain/physical-location';
import type { DigitalLocation } from '@/types/domain/digital-location';
import { BoxIcon } from 'lucide-react';
import { IconCloudDataConnection } from '@tabler/icons-react';

// AddGameToLibraryFormSchema
export const AddGameToLibraryFormSchema = z.object({
  platformName: z.string({
    required_error: 'Please select a game platform'
  }),
  storageType: z.enum(['physical', 'digital', 'both'], {
    required_error: 'Please select where you will store this game'
  }),
  gameLocations: z.record(z.array(z.string())).refine(
    (locations) => Object.values(locations).some(loc => loc.length > 0),
    { message: 'Please select at least one platform for a location' }
  )
});

const storageOptions: AddToLibraryFormStorageOption[] = [
  {
    value: 'physical',
    label: 'In a physical location',
    icon: <BoxIcon />
  },
  {
    value: 'digital',
    label: 'In an online service',
    icon: <IconCloudDataConnection />
  },
  {
    value: 'both',
    label: 'I have both a digital and physical copy of this game',
    icon: (
      <div className="flex items-center gap-1">
        <BoxIcon />
        <span className="text-sm">+</span>
        <IconCloudDataConnection />
      </div>
    )
  }
];

interface AddGameToLibraryFormProps {
  selectedGame: Game;
  physicalLocations: PhysicalLocation[];
  digitalLocations: DigitalLocation[];
}

export function AddGameToLibraryForm({
  selectedGame,
  physicalLocations,
  digitalLocations,
}: AddGameToLibraryFormProps) {
  const form = useForm<z.infer<typeof AddGameToLibraryFormSchema>>({
    resolver: zodResolver(AddGameToLibraryFormSchema),
    defaultValues: {
      platformName: selectedGame.platformNames?.length === 1 ? selectedGame.platformNames[0] : "",
      storageType: undefined,
      gameLocations: {}
    },
  });

  const selectedStorageType = form.watch('storageType');

  return (
    <ErrorBoundary FallbackComponent={FormErrorFallback}>
      <FormContainer form={form} onSubmit={() => { console.log("Form submitted") }}>
        {/* Storage Type Selection */}
        <FormField
          control={form.control}
          name="storageType"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Where will you store this game?</FormLabel>
              <RadioGroup.Root
                onValueChange={field.onChange}
                value={field.value}
                className="flex flex-col gap-4"
              >
                {storageOptions.map((option) => (
                  <FormControl key={option.value}>
                    <RadioGroup.Item
                      value={option.value}
                      className={cn(
                        "relative flex items-center gap-4 p-4 rounded-lg border",
                        "data-[state=checked]:border-primary data-[state=checked]:bg-primary/5"
                      )}
                    >
                      <div className="flex items-center gap-4 w-full">
                        {option.icon && (
                          <div className="flex-shrink-0">{option.icon}</div>
                        )}
                        <span className="font-medium">{option.label}</span>
                        <RadioGroup.Indicator className="absolute right-4 top-4 flex h-4 w-4 items-center justify-center">
                          <div className="h-2.5 w-2.5 rounded-full bg-primary" />
                        </RadioGroup.Indicator>
                      </div>
                    </RadioGroup.Item>
                  </FormControl>
                ))}
              </RadioGroup.Root>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Location Selection */}
        {selectedStorageType && (
          <FormField
            control={form.control}
            name="gameLocations"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Please tell us which version(s) of the game you have and where they are kept</FormLabel>
                <FormControl>
                  <div className="space-y-4">
                    {selectedStorageType !== 'physical' && digitalLocations.length > 0 && (
                      <div className="space-y-4">
                        <div className="font-bold">Digital Locations</div>
                        {digitalLocations.map((digiLocation) => (
                          <div key={digiLocation.id} className="space-y-2">
                            <div className="font-bold">{digiLocation.name}</div>
                            <MultiSelect
                              value={field.value[digiLocation.id] || []}
                              onChange={(value) => {
                                field.onChange({
                                  ...field.value,
                                  [digiLocation.id]: value
                                });
                              }}
                              onBlur={field.onBlur}
                              name={`gameLocations.${digiLocation.id}`}
                              mainPlaceholder="Which game version are you adding?"
                              secondaryPlaceholder="Available platforms"
                              platformNames={selectedGame.platformNames || []}
                            />
                          </div>
                        ))}
                      </div>
                    )}

                    {selectedStorageType !== 'digital' && (
                      <div className="space-y-4">
                        <div className="font-bold">Physical Locations</div>
                        {physicalLocations.map((location) => (
                          <div key={location.id} className="space-y-2">
                            {location.sublocations?.map((sublocation) => (
                              <div key={sublocation.id}>
                                <div className="font-bold">{sublocation.name} | {location.name}</div>
                                <MultiSelect
                                  value={field.value[sublocation.id] || []}
                                  onChange={(value) => {
                                    field.onChange({
                                      ...field.value,
                                      [sublocation.id]: value
                                    });
                                  }}
                                  onBlur={field.onBlur}
                                  name={`gameLocations.${sublocation.id}`}
                                  mainPlaceholder="Which game version are you adding?"
                                  secondaryPlaceholder="Available platforms"
                                  platformNames={selectedGame.platformNames || []}
                                />
                              </div>
                            ))}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        )}
      </FormContainer>
    </ErrorBoundary>
  );
}
