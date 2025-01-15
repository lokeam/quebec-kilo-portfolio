
// Components
import { FormContainer } from '@/features/dashboard/components/templates/FormContainer';

// Shadcn UI Components
import { Input } from "@/shared/components/ui/input"

import { Button } from "@/shared/components/ui/button";

// Icons
import { BookshelfIcon } from '@/shared/components/ui/CustomIcons/BookShelfIcon';
import { MediaConsoleIcon } from '@/shared/components/ui/CustomIcons/MediaConsoleIcon';

import { DrawerIcon } from '@/shared/components/ui/CustomIcons/DrawerIcon';
import { CabinetIcon } from '@/shared/components/ui/CustomIcons/CabinetIcon';
import { ClosetIcon } from '@/shared/components/ui/CustomIcons/ClosetIcon';
import { Package } from 'lucide-react';

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

const FormSchema = z.object({
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
  bgColor: z
    .string({
      required_error: "Please select a background color",
    }),
});

interface MediaPageSublocationFormProps {
  onSuccess?: () => void;
}

export function MediaPageSublocationForm({ onSuccess}: MediaPageSublocationFormProps) {
  /* Specific form components creates their own useForm hook instances */
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      locationName: '',
      locationType: '',
      bgColor: '',
    }
  });

  const handleSubmit = (data: z.infer<typeof FormSchema>) => {
    toast(`Form submitted with the following data: ${JSON.stringify(data)}`, {
      className: 'bg-green-500 text-white',
      duration: 2500,
    });
    onSuccess?.();
  };


  return (
    <FormContainer form={form} onSubmit={handleSubmit}>

      {/* Location Name */}
      <FormField
        control={form.control}
        name="locationName"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Sublocation Name</FormLabel>
            <FormControl>
              <Input placeholder="Example: Study bookshelf A" {...field} />
            </FormControl>
            <FormDescription>
              This is the area in your main location where the media is stored.
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
            <FormLabel>Storage Unit Type</FormLabel>
            <Select onValueChange={field.onChange} defaultValue={field.value}>
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Where are you keeping your media?" />
                </SelectTrigger>
              </FormControl>

              <SelectContent>
                <SelectItem value="shelf">
                  <div className="flex items-center gap-2">
                    <BookshelfIcon size={20} color='#fff' className='mr-2'/>
                    <span>Shelf / Shelving unit</span>
                  </div>
                </SelectItem>
                <SelectItem value="console">
                  <div className="flex items-center gap-2">
                    <MediaConsoleIcon size={20} color='#fff' className='mr-2'/>
                    <span>Media console</span>
                  </div>
                </SelectItem>
                <SelectItem value="cabinet">
                  <div className="flex items-center gap-2">
                    <CabinetIcon size={20} color='#fff' className='mr-2'/>
                    <span>Cabinet</span>
                  </div>
                </SelectItem>
                <SelectItem value="closet">
                  <div className="flex items-center gap-2">
                    <ClosetIcon size={20} color='#fff' className='mr-2'/>
                    <span>Closet</span>
                  </div>
                  </SelectItem>
                <SelectItem value="drawer">
                  <div className="flex items-center gap-2">
                    <DrawerIcon size={20} color='#fff' className='mr-2'/>
                    <span>Drawer</span>
                  </div>
                </SelectItem>
                <SelectItem value="box">
                  <div className="flex items-center gap-2">
                    <Package size={20} color='#fff' className='mr-2'/>
                    <span>Storage container</span>
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

      {/* Icon BG Color */}
      <FormField
        control={form.control}
        name="bgColor"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Icon Background Color</FormLabel>
            <Select onValueChange={field.onChange} defaultValue={field.value}>
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select a background color for your icon" />
                </SelectTrigger>
              </FormControl>

              <SelectContent>
                <SelectItem value="red">Red</SelectItem>
                <SelectItem value="blue">Blue</SelectItem>
                <SelectItem value="green">Green</SelectItem>
                <SelectItem value="purple">Purple</SelectItem>
                <SelectItem value="orange">Orange</SelectItem>
                <SelectItem value="yellow">Yellow</SelectItem>
                <SelectItem value="gray">Gray</SelectItem>
                <SelectItem value="white">White</SelectItem>
                <SelectItem value="black">Black</SelectItem>
              </SelectContent>
            </Select>
            <FormDescription>
              Customize the background color of your storage unit icon.
            </FormDescription>

            <FormMessage />
          </FormItem>
        )}
      />



      <Button type="submit" className="w-full">Submit</Button>
    </FormContainer>
  )
}
