import { useState, useEffect } from 'react';

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
import { useCreateLibraryGame } from '@/core/api/queries/gameLibrary.queries';

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
  platformId: number;
  type: 'physical' | 'digital';
  location: {
    sublocationId?: string;
    digitalLocationId?: string;
  };
}

export interface AddToLibraryFormPayload {
  gameId: number; // IGDB game ID
  gameName: string;
  gameCoverUrl: string;
  gameFirstReleaseDate: number;
  gameType: {
    displayText: string;
    normalizedText: string;
  };
  gameThemeNames: string[];
  gamesByPlatformAndLocation: AddToLibraryFormLocationEntry[];
  gameRating: number;
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
  gameLocations: z.record(z.array(z.object({
    id: z.number(),
    name: z.string(),
  }))).refine(
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
  onClose?: () => void;
}

export function AddGameToLibraryForm({
  selectedGame,
  physicalLocations,
  digitalLocations,
  onClose,
}: AddGameToLibraryFormProps) {
  const [isFormSubmitting, setIsFormSubmitting] = useState(false);

  // Add detailed debug logging without stringifying
  // console.log('AddGameToLibraryForm - Raw Props:', {
  //   selectedGame,
  //   physicalLocations,
  //   digitalLocations,
  // });

  // Type check the physical locations
  // console.log('Physical Locations Type Check:', {
  //   isArray: Array.isArray(physicalLocations),
  //   length: physicalLocations?.length,
  //   firstLocation: physicalLocations?.[0],
  //   firstLocationSublocations: physicalLocations?.[0]?.sublocations,
  // });

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

  // Add debug logging for form state
  console.log('Form state:', {
    selectedStorageType,
    selectedLocations,
    selectedLocationsCount,
    formState: form.formState,
    isValid: form.formState.isValid,
    errors: form.formState.errors,
  });

  // Add effect to log when storage type changes
  useEffect(() => {
    console.log('Storage type changed:', {
      selectedStorageType,
      physicalLocations,
      digitalLocations,
    });
  }, [selectedStorageType, physicalLocations, digitalLocations]);

  const createMutation = useCreateLibraryGame();

  const onSubmit = async (data: z.infer<typeof AddGameToLibraryFormSchema>) => {
    try {
      setIsFormSubmitting(true);
      const { gameLocations } = data;

      // Log the form data before transformation
      console.log('🔍 DEBUG: Form submission - Raw form data:', {
        data,
        gameLocations,
        physicalLocations,
        digitalLocations,
      });

      // Transform data to an efficient payload
      const gamesByPlatformAndLocation = Object.entries(gameLocations).flatMap(([locationId, platforms]) => {
        const isDigital = digitalLocations.some(loc => loc.id === locationId);

        // Log each location transformation
        console.log('🔍 DEBUG: Transforming location:', {
          locationId,
          platforms,
          isDigital,
          physicalLocation: physicalLocations.find(loc => loc.sublocations?.some(sub => sub.id === locationId)),
          digitalLocation: digitalLocations.find(loc => loc.id === locationId),
        });

        return platforms.map((platform): AddToLibraryFormLocationEntry => ({
          platformName: platform.name,
          platformId: platform.id,
          type: isDigital ? 'digital' : 'physical',
          location: isDigital
            ? { digitalLocationId: locationId }
            : { sublocationId: locationId },
        }));
      });

      // Log the transformed data
      console.log('🔍 DEBUG: Form submission - Transformed data:', {
        gamesByPlatformAndLocation,
      });

      const payload: AddToLibraryFormPayload = {
        gameId: selectedGame.id,
        gameName: selectedGame.name,
        gameCoverUrl: selectedGame.coverUrl || '',
        gameFirstReleaseDate: selectedGame.firstReleaseDate || 0,
        gameType: selectedGame.gameType || {
          displayText: '',
          normalizedText: ''
        },
        gameThemeNames: selectedGame.themeNames || [],
        gamesByPlatformAndLocation,
        gameRating: selectedGame.rating || 0,
      };

      // Log the final payload
      console.log('🔍 DEBUG: Form submission - Final payload:', payload);

      await createMutation.mutateAsync(payload);
      if (onClose) onClose();

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
          render={({ field }) => {
            // Log the field state
            // console.log('Storage type field:', {
            //   value: field.value,
            //   onChange: field.onChange,
            // });

            return (
              <FormItem>
                <FormLabel>Where will you store this game?</FormLabel>
                <RadioGroup.Root
                  onValueChange={(value) => {
                    console.log('Storage type changed to:', value);
                    field.onChange(value);
                  }}
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
            );
          }}
        />

        {/* Location Selection */}
        {selectedStorageType && (
          <>
            <FormField
              control={form.control}
              name="selectedLocations"
              render={({ field }) => {
                // Add debug logging for field render
                // console.log('selectedLocations field:', {
                //   value: field.value,
                //   onChange: field.onChange,
                // });

                return (
                  <FormItem>
                    <FormLabel>Please tell us which version(s) of the game you have and where they are kept</FormLabel>
                    <FormControl>
                      <div className="space-y-4">
                        {selectedStorageType !== 'physical' && digitalLocations.length > 0 && (
                          <div className="space-y-4">
                            <div className="font-bold">Digital Locations</div>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                              {digitalLocations.map((digiLocation) => {
                                // Add debug logging for digital locations
                                console.log('Rendering digital location:', digiLocation);
                                return (
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
                                );
                              })}
                            </div>
                          </div>
                        )}

                        {selectedStorageType !== 'digital' && (
                          <div className="space-y-4">
                            <div className="font-bold">Physical Locations</div>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                              {Array.isArray(physicalLocations) && physicalLocations.map((location) => {
                                // Add detailed debug logging for physical locations
                                // console.log('Rendering physical location:', {
                                //   locationId: location.id,
                                //   locationName: location.name,
                                //   sublocations: location.sublocations,
                                //   sublocationsLength: location.sublocations?.length,
                                //   sublocationsType: typeof location.sublocations,
                                //   isArray: Array.isArray(location.sublocations),
                                // });

                                return Array.isArray(location.sublocations) && location.sublocations.map((sublocation) => {
                                  // Add detailed debug logging for sublocations
                                  // console.log('Rendering sublocation:', {
                                  //   sublocationId: sublocation.id,
                                  //   sublocationName: sublocation.name,
                                  //   sublocationType: sublocation.type,
                                  //   parentLocation: location.name,
                                  // });

                                  return (
                                    <CheckboxPrimitive.Root
                                      key={sublocation.id}
                                      checked={field.value[sublocation.id] || false}
                                      onCheckedChange={(checked) => {
                                        // console.log('Checkbox changed:', {
                                        //   sublocationId: sublocation.id,
                                        //   checked,
                                        //   currentValue: field.value,
                                        // });
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
                                  );
                                });
                              })}
                            </div>
                          </div>
                        )}
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                );
              }}
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
                                  platforms={selectedGame.platforms || []}
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
                                    platforms={selectedGame.platforms || []}
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
            className="w-full"
            disabled={!form.formState.isValid || !selectedStorageType || isFormSubmitting || createMutation.isPending}
          >
            {isFormSubmitting || createMutation.isPending
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
