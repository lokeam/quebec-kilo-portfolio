import { useState } from 'react';

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
import { Button } from '@/shared/components/ui/button';

import * as RadioGroup from "@radix-ui/react-radio-group";
import * as CheckboxPrimitive from "@radix-ui/react-checkbox";
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

interface AddToLibraryFormLocationEntry {
  platformName: string;
  type: 'physical' | 'digital';
  location: {
    sublocationId?: string;
    digitalLocationId?: string;
  };
}

interface AddToLibraryFormPayload {
  gameId: number;
  gamesByPlatformAndLocation: AddToLibraryFormLocationEntry[];
}


// Zod
import { z } from 'zod';

// Types
import type { Game } from '@/types/domain/game';
import type { PhysicalLocation } from '@/types/domain/physical-location';
import type { DigitalLocation } from '@/types/domain/digital-location';
import { BoxIcon, CheckCircle2 } from 'lucide-react';
import { IconCloudDataConnection } from '@tabler/icons-react';

// AddGameToLibraryFormSchema
export const AddGameToLibraryFormSchema = z.object({
  storageType: z.enum(['physical', 'digital', 'both'], {
    required_error: 'Please select where you will store this game'
  }),
  selectedLocations: z.record(z.boolean()).refine(
    (locations) => Object.values(locations).some(loc => loc === true),
    { message: 'Please select at least one location' }
  ),
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
  const [isFormSubmitting, setIsFormSubmitting] = useState(false);

  const form = useForm<z.infer<typeof AddGameToLibraryFormSchema>>({
    resolver: zodResolver(AddGameToLibraryFormSchema),
    defaultValues: {
      storageType: undefined,
      gameLocations: {},
      selectedLocations: {}
    },
  });

  const selectedStorageType = form.watch('storageType');
  const selectedLocations = form.watch('selectedLocations');
  const selectedLocationsCount = Object.values(selectedLocations).filter(Boolean).length;

  const onSubmit = (data: z.infer<typeof AddGameToLibraryFormSchema>) => {
    try {
      setIsFormSubmitting(true);
      const { gameLocations } = data;

      // Transform data to an efficient payload
      const gamesByPlatformAndLocation = Object.entries(gameLocations).flatMap(([locationId, platforms]) => {
        // Determine if this is a digital location
        const isDigital = digitalLocations.some(loc => loc.id === locationId);

        // Map each platform to an entry
        return platforms.map((platformName): AddToLibraryFormLocationEntry => ({
          platformName,
          type: isDigital ? 'digital' : 'physical',
          location: isDigital
            ? { digitalLocationId: locationId }
            : { sublocationId: locationId }
        }));
      });

      const payload: AddToLibraryFormPayload = {
        gameId: selectedGame.id,
        gamesByPlatformAndLocation
      };

      // TODO: Replace with proper library post query
      console.log('Form Payload:', payload);
    } catch (error) {
      console.error('Error submitting form:', error);
    } finally {
      setIsFormSubmitting(false);
    }
  };

  return (
    <ErrorBoundary FallbackComponent={FormErrorFallback}>
      <FormContainer form={form} onSubmit={onSubmit}>
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
          <>
            <FormField
              control={form.control}
              name="selectedLocations"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Please tell us which version(s) of the game you have and where they are kept</FormLabel>
                  <FormControl>
                    <div className="space-y-4">
                      {selectedStorageType !== 'physical' && digitalLocations.length > 0 && (
                        <div className="space-y-4">
                          <div className="font-bold">Digital Locations</div>
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                            {digitalLocations.map((digiLocation) => (
                              <CheckboxPrimitive.Root
                                key={digiLocation.id}
                                checked={field.value[digiLocation.id] || false}
                                onCheckedChange={(checked) => {
                                  field.onChange({
                                    ...field.value,
                                    [digiLocation.id]: checked
                                  });
                                }}
                                className="relative ring-[1px] ring-border rounded-lg px-4 py-3 text-start text-muted-foreground data-[state=checked]:ring-2 data-[state=checked]:ring-primary data-[state=checked]:text-primary"
                              >
                                <IconCloudDataConnection className="mb-3" />
                                <span className="font-medium tracking-tight">Add game to {digiLocation.name}</span>
                                <CheckboxPrimitive.Indicator className="absolute top-2 right-2">
                                  <CheckCircle2 className="fill-primary text-primary-foreground" />
                                </CheckboxPrimitive.Indicator>
                              </CheckboxPrimitive.Root>
                            ))}
                          </div>
                        </div>
                      )}

                      {selectedStorageType !== 'digital' && (
                        <div className="space-y-4">
                          <div className="font-bold">Physical Locations</div>
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                            {physicalLocations.map((location) => (
                              location.sublocations?.map((sublocation) => (
                                <CheckboxPrimitive.Root
                                  key={sublocation.id}
                                  checked={field.value[sublocation.id] || false}
                                  onCheckedChange={(checked) => {
                                    field.onChange({
                                      ...field.value,
                                      [sublocation.id]: checked
                                    });
                                  }}
                                  className="relative ring-[1px] ring-border rounded-lg px-4 py-3 text-start text-muted-foreground data-[state=checked]:ring-2 data-[state=checked]:ring-primary data-[state=checked]:text-primary"
                                >
                                  <BoxIcon className="mb-3" />
                                  <span className="font-medium tracking-tight">Add game to {sublocation.name}</span>
                                  <CheckboxPrimitive.Indicator className="absolute top-2 right-2">
                                    <CheckCircle2 className="fill-primary text-primary-foreground" />
                                  </CheckboxPrimitive.Indicator>
                                </CheckboxPrimitive.Root>
                              ))
                            ))}
                          </div>
                        </div>
                      )}
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="gameLocations"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <div className="space-y-4">
                      {selectedStorageType !== 'physical' && digitalLocations.length > 0 && (
                        <div className="space-y-4">
                          {digitalLocations.map((digiLocation) => (
                            selectedLocations[digiLocation.id] && (
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
                            )
                          ))}
                        </div>
                      )}

                      {selectedStorageType !== 'digital' && (
                        <div className="space-y-4">
                          {physicalLocations.map((location) => (
                            location.sublocations?.map((sublocation) => (
                              selectedLocations[sublocation.id] && (
                                <div key={sublocation.id} className="space-y-2">
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
                              )
                            ))
                          ))}
                        </div>
                      )}
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </>
        )}

        {/* Submit Button */}
        <div className="mt-6">
          <Button
            type="submit"
            className="w-full md:w-auto"
            disabled={!form.formState.isValid || !selectedStorageType || isFormSubmitting}
          >
            {isFormSubmitting
              ? 'Adding to Library...'
              : selectedLocationsCount > 0
                ? `Add to Library (${selectedLocationsCount} location${selectedLocationsCount === 1 ? '' : 's'})`
                : 'Add to Library'}
          </Button>
        </div>
      </FormContainer>
    </ErrorBoundary>
  );
}
