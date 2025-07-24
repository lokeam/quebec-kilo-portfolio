import { useState, useEffect, useMemo } from 'react';

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
import type { AddGameFormPhysicalLocationsResponse, AddGameFormDigitalLocationsResponse } from '@/types/domain/search';
import { BoxIcon, CheckCircle2, SquareDashed } from '@/shared/components/ui/icons';
import { IconCloudDataConnection } from '@/shared/components/ui/icons';

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
  physicalLocations: AddGameFormPhysicalLocationsResponse[];
  digitalLocations: AddGameFormDigitalLocationsResponse[];
  onClose?: () => void;
}

export function AddGameToLibraryForm({
  selectedGame,
  physicalLocations,
  digitalLocations,
  onClose,
}: AddGameToLibraryFormProps) {
  const [isFormSubmitting, setIsFormSubmitting] = useState(false);

  // Determine form type at component start
  const formType = useMemo(() => {
    if (physicalLocations.length > 0 && digitalLocations.length === 0) return 'physical';
    if (physicalLocations.length === 0 && digitalLocations.length > 0) return 'digital';
    return 'both';
  }, [physicalLocations.length, digitalLocations.length]);

  const form = useForm<z.infer<typeof AddGameToLibraryFormSchema>>({
    resolver: zodResolver(AddGameToLibraryFormSchema),
    defaultValues: {
      storageType: formType === 'both' ? undefined : formType,
      selectedLocations: {},
      gameLocations: {}
    }
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
        const isDigital = digitalLocations.some(loc => loc.digitalLocationId === locationId);

        // Log each location transformation
        console.log('🔍 DEBUG: Transforming location:', {
          locationId,
          platforms,
          isDigital,
          physicalLocation: physicalLocations.find(loc => loc.sublocationId === locationId),
          digitalLocation: digitalLocations.find(loc => loc.digitalLocationId === locationId),
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
        {/* Show message when no locations are available */}
        {physicalLocations.length === 0 && digitalLocations.length === 0 && (
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <SquareDashed className="h-12 w-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">Where do you keep your games IRL?</h3>
            <p className="text-muted-foreground">
              Create at least one <span className="font-bold text-primary">online service</span> or <span className="font-bold text-primary">physical location</span> before adding games to your library
            </p>
          </div>
        )}

        {/* Storage Type Selection - Only show if both types are available and locations exist */}
        {formType === 'both' && !(physicalLocations.length === 0 && digitalLocations.length === 0) && (
          <FormField
            control={form.control}
            name="storageType"
            render={({ field }) => (
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
                  {storageOptions
                    .filter((option: AddToLibraryFormStorageOption) => {
                      if (option.value === 'physical') return physicalLocations.length > 0;
                      if (option.value === 'digital') return digitalLocations.length > 0;
                      if (option.value === 'both') return physicalLocations.length > 0 && digitalLocations.length > 0;
                      return false;
                    })
                    .map((option) => (
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
        )}

        {/* Location Selection - Only show if locations exist */}
        {!(physicalLocations.length === 0 && digitalLocations.length === 0) && (formType === 'both' ? selectedStorageType : formType) && (
          <>
            <FormField
              control={form.control}
              name="selectedLocations"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className='mb-2'>Storage location and game platform:</FormLabel>
                  <FormControl>
                    <div className="space-y-4">
                      {(formType === 'both' ? selectedStorageType !== 'physical' : formType === 'digital') && digitalLocations.length > 0 && (
                        <div className="space-y-4">
                          <div className="font-bold">Digital Locations</div>
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                            {digitalLocations.map((digiLocation) => (
                              <CheckboxPrimitive.Root
                                key={digiLocation.digitalLocationId}
                                checked={field.value[digiLocation.digitalLocationId] || false}
                                onCheckedChange={(checked) => {
                                  field.onChange({
                                    ...field.value,
                                    [digiLocation.digitalLocationId]: checked
                                  });
                                }}
                                className="relative ring-[1px] ring-border rounded-lg px-4 py-3 text-start text-muted-foreground data-[state=checked]:ring-2 data-[state=checked]:ring-primary data-[state=checked]:text-primary"
                              >
                                <IconCloudDataConnection className="mb-3" />
                                <span className="font-medium tracking-tight">Add game to {digiLocation.digitalLocationName}</span>
                                <CheckboxPrimitive.Indicator className="absolute top-2 right-2">
                                  <CheckCircle2 className="fill-primary text-primary-foreground" />
                                </CheckboxPrimitive.Indicator>
                              </CheckboxPrimitive.Root>
                            ))}
                          </div>
                        </div>
                      )}

                      {(formType === 'both' ? selectedStorageType !== 'digital' : formType === 'physical') && (
                        <div className="space-y-4">
                          <div className="font-bold">Physical Locations</div>
                          <div className="grid grid-cols-1 gap-3">
                            {physicalLocations.map((location) => (
                              <CheckboxPrimitive.Root
                                key={location.sublocationId}
                                checked={field.value[location.sublocationId] || false}
                                onCheckedChange={(checked) => {
                                  field.onChange({
                                    ...field.value,
                                    [location.sublocationId]: checked
                                  });
                                }}
                                className="relative ring-[1px] ring-border rounded-lg px-4 py-3 text-start text-muted-foreground data-[state=checked]:ring-2 data-[state=checked]:ring-primary data-[state=checked]:text-primary"
                              >
                                <BoxIcon className="mb-3" />
                                <span className="font-medium tracking-tight">Add game to {location.sublocationName}</span>
                                <CheckboxPrimitive.Indicator className="absolute top-2 right-2">
                                  <CheckCircle2 className="fill-primary text-primary-foreground" />
                                </CheckboxPrimitive.Indicator>
                              </CheckboxPrimitive.Root>
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
                      {(formType === 'both' ? selectedStorageType !== 'physical' : formType === 'digital') && digitalLocations.length > 0 && (
                        <div className="space-y-4">
                          {digitalLocations.map((digiLocation) => (
                            selectedLocations[digiLocation.digitalLocationId] && (
                              <div key={digiLocation.digitalLocationId} className="space-y-2">
                                <div className="font-bold">{digiLocation.digitalLocationName}</div>
                                <MultiSelect
                                  value={field.value[digiLocation.digitalLocationId] || []}
                                  onChange={(value) => {
                                    field.onChange({
                                      ...field.value,
                                      [digiLocation.digitalLocationId]: value
                                    });
                                  }}
                                  onBlur={field.onBlur}
                                  name={`gameLocations.${digiLocation.digitalLocationId}`}
                                  mainPlaceholder="Which game version are you adding?"
                                  secondaryPlaceholder="Available platforms"
                                  platforms={selectedGame.platforms || []}
                                />
                              </div>
                            )
                          ))}
                        </div>
                      )}

                      {(formType === 'both' ? selectedStorageType !== 'digital' : formType === 'physical') && (
                        <div className="space-y-4">
                          {physicalLocations.map((location) => (
                            selectedLocations[location.sublocationId] && (
                              <div key={location.sublocationId} className="space-y-2">
                                <div className="font-bold">{location.sublocationName} | {location.parentLocationName}</div>
                                <MultiSelect
                                  value={field.value[location.sublocationId] || []}
                                  onChange={(value) => {
                                    field.onChange({
                                      ...field.value,
                                      [location.sublocationId]: value
                                    });
                                  }}
                                  onBlur={field.onBlur}
                                  name={`gameLocations.${location.sublocationId}`}
                                  mainPlaceholder="Which game version are you adding?"
                                  secondaryPlaceholder="Available platforms"
                                  platforms={selectedGame.platforms || []}
                                />
                              </div>
                            )
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

        {/* Submit Button - Only show if locations are available */}
        {!(physicalLocations.length === 0 && digitalLocations.length === 0) && (
          <div className="mt-6">
            <Button
              type="submit"
              className="w-full"
              disabled={!form.formState.isValid || (formType === 'both' ? !selectedStorageType : false) || isFormSubmitting || createMutation.isPending}
            >
              {isFormSubmitting || createMutation.isPending
                ? 'Adding to Library...'
                : selectedLocationsCount > 0
                  ? `Add to Library (${selectedLocationsCount} location${selectedLocationsCount === 1 ? '' : 's'})`
                  : 'Add to Library'}
            </Button>
          </div>
        )}
      </FormContainer>
    </ErrorBoundary>
  );
}
