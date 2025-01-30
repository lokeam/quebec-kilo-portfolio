
// Components
import { FormContainer } from '@/features/dashboard/components/templates/FormContainer';

// Shadcn UI Components
import { Input } from "@/shared/components/ui/input"
import { Switch } from "@/shared/components/ui/switch"

import { Button } from "@/shared/components/ui/button";

import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/components/ui/form"

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/components/ui/select"

// Hooks
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { toast } from "sonner"

// Zod
import { z } from "zod"

// Icons
import { House, Building, Building2, Warehouse } from 'lucide-react';
import { IconCar } from '@tabler/icons-react';

export const FormSchema = z.object({
  locationName: z
    .string({
      required_error: "Please enter a location name",
    })
    .min(3, {
      message: "Location name must be at least 3 characters long",
    }),
  locationType: z
    .string({
      required_error: "Please select a location type",
    }),
    coordinates: z.object({
      enabled: z.boolean().default(false),
      value: z.string().optional().superRefine((val, ctx) => {
        if (val === undefined && (ctx as z.RefinementCtx & { parent: { enabled: boolean } }).parent.enabled) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            message: "Coordinates are required when enabled",
          });
        }
      }),
    }).default({ enabled: false, value: undefined }),
});

interface MediaPageLocationFormProps {
  onSuccess?: (data: z.infer<typeof FormSchema>) => void;
  defaultValues?: z.infer<typeof FormSchema>;
  buttonText?: string;
}

export function MediaPageLocationForm({
  buttonText = "Submit",
  onSuccess,
  defaultValues = {
    locationName: '',
    locationType: '',
    coordinates: {
      enabled: false,
      value: '' // Ensure value is never undefined
    }
  }
}: MediaPageLocationFormProps) {
  /* Specific form components creates their own useForm hook instances */
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues
  });

  const handleSubmit = (data: z.infer<typeof FormSchema>) => {
    toast(`Form submitted with the following data: ${JSON.stringify(data)}`, {
      className: 'bg-green-500 text-white',
      duration: 2500,
    });
    onSuccess?.(data);
  };

  return (
    <FormContainer form={form} onSubmit={handleSubmit}>

      {/* Location Name */}
      <FormField
        control={form.control}
        name="locationName"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Location Name</FormLabel>
            <FormControl>
              <Input placeholder="Enter a location name" {...field} />
            </FormControl>
            <FormDescription>
              This is the name of the location where the media is stored.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      {/* Location Type */}
      <FormField
        control={form.control}
        name="locationType"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Location Type</FormLabel>
            <Select onValueChange={field.onChange} defaultValue={field.value}>
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select a location type" />
                </SelectTrigger>
              </FormControl>

              <SelectContent>
                <SelectItem value="house">
                <div className="flex items-center gap-2">
                  <House size={20} color='#fff' className='mr-2'/>
                  <span>House</span>
                </div>

                </SelectItem>
                <SelectItem value="apartment">
                  <div className="flex items-center gap-2">
                    <Building size={20} color='#fff' className='mr-2'/>
                    <span>Apartment</span>
                  </div>
                </SelectItem>
                <SelectItem value="office">
                  <div className="flex items-center gap-2">
                    <Building2 size={20} color='#fff' className='mr-2'/>
                    <span>Office</span>
                  </div>
                </SelectItem>
                <SelectItem value="commercialStorage">
                  <div className="flex items-center gap-2">
                    <Warehouse size={20} color='#fff' className='mr-2'/>
                    <span>Commercial Storage</span>
                  </div>
                </SelectItem>
                <SelectItem value="vehicle">
                  <div className="flex items-center gap-2">
                    <IconCar size={25} color='#fff' className='mr-2'/>
                    <span>Vehicle</span>
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
            <FormDescription>
              Think of this as the venue where the media is stored.
            </FormDescription>

            <FormMessage />
          </FormItem>
        )}
      />

      {/* Coordinates */}
      <FormField
        control={form.control}
        name="coordinates"
        render={({ field }) => (
          <FormItem className="space-y-4">
            <div className="flex flex-row items-center justify-between rounded-lg border p-4">
              <div className="space-y-0.5">
                <FormLabel className="text-base">
                  Coordinates
                </FormLabel>
                <FormDescription>
                  Optionally add map coordinates for the location.
                </FormDescription>
              </div>

              <FormControl>
                <Switch
                  checked={field.value?.enabled}
                  onCheckedChange={(checked) => {
                    field.onChange({
                      enabled: checked,
                      value: checked ? field.value?.value : undefined
                    });
                  }}
                />
              </FormControl>
            </div>
            {field.value?.enabled && (
              <FormItem>
                <FormControl>
                  <Input
                    placeholder="Enter coordiantes"
                    value={field.value.value ?? ''}
                    onChange={(event) => {
                      field.onChange({
                        enabled: true,
                        value: event.target.value || '' // Ensure value is never undefined
                      });
                    }}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          </FormItem>
        )}
      />

      <Button type="submit" className="w-full">{buttonText}</Button>
    </FormContainer>
  )
}
