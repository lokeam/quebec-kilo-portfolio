import { useState } from 'react';

// RHF
import { useForm } from 'react-hook-form';

// Zod Validation
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';

// Shadcn UI Components
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/components/ui/form';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';
import { Button } from '@/shared/components/ui/button';
import { Calendar } from '@/shared/components/ui/calendar';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/shared/components/ui/popover';
import { Textarea } from '@/shared/components/ui/textarea';
import { CalendarIcon, ChevronDown } from 'lucide-react';
import { format } from 'date-fns';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/shared/components/ui/collapsible';
import { cn } from '@/shared/components/ui/utils';

import { CARD_PARENT_CLASS } from '@/features/dashboard/lib/constantsstyle.constants';

// Form Schema
const formSchema = z.object({

});


export function OnlineServiceForm() {
  const [date, setDate] = useState<Date>()
  const [selectedColor, setSelectedColor] = useState<string>("")
  const [isOpen, setIsOpen] = useState(false)

  const colors = [
    "Red",
    "Blue",
    "Green",
    "Yellow",
    "Purple",
    "Orange",
    "Pink",
    "Brown",
    "Gray",
    "Black",
  ];

  // Define form
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      serviceName: '',
      color: '',
      serviceURL: '',
      cost: 0.00,
      subscriptionType: 'recurring',
      period: 'month',
      nextPayment: new Date(),
    },
  });

  // Submit Handler
  function onSubmit(values: z.infer<typeof formSchema>) {
    // Do something w/ form values
    console.log(`form values: ${JSON.stringify(values)}`);
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
      >

        <div className="w-full max-w-2xl mx-auto p-6 space-y-8">
          {/* Title */}
          <div className="flex items-center gap-2 mb-6">
            <Button variant="ghost" className="w-8 h-8 p-0">
              <ChevronDown className="h-4 w-4" />
            </Button>
            <h1 className="text-xl font-semibold">New subscription</h1>
          </div>

          <div className="space-y-6">
            {/* General Section */}
            <div className={CARD_PARENT_CLASS}>
              <h2 className="text-base font-medium mb-4">General</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Online Service Name <span className="text-red-500">*</span></Label>
                  <Input id="serviceName" placeholder="Enter Online Service Name" />
                </div>
                <div className="space-y-2">
                  <Label>Color</Label>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="outline" className="w-[42px] p-2">
                        <div
                          className={cn(
                            "h-full w-full rounded",
                            selectedColor ? "bg-" + selectedColor.toLowerCase() + "-500" : "bg-gray-200"
                          )}
                        />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent>
                      {colors.map((color) => (
                        <DropdownMenuItem
                          key={color}
                          onClick={() => setSelectedColor(color)}
                          className="flex items-center gap-2"
                        >
                          <div className={`w-4 h-4 rounded bg-${color.toLowerCase()}-500`} />
                          {color}
                        </DropdownMenuItem>
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="website">Website</Label>
                  <Input id="serviceURL" defaultValue="https://adobe.com" />
                </div>
              </div>
            </div>

            {/* Expense Section */}
            <div className={CARD_PARENT_CLASS}>
              <h2 className="text-base font-medium mb-4">Expense</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="cost">Cost *</Label>
                  <div className="flex">
                    <div className="flex-none flex items-center px-3 border border-r-0 rounded-l-md bg-gray-50">
                      <span className="text-sm text-gray-500">$</span>
                    </div>
                    <Input id="cost" defaultValue={0.00} className="rounded-l-none" />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="expense-type">Subscription Type *</Label>
                  <Select>
                    <SelectTrigger>
                      <SelectValue placeholder="Select type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="recurring">Recurring</SelectItem>
                      <SelectItem value="one-time">Free</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </div>

            {/* Billing Section */}
            <div className={CARD_PARENT_CLASS}>
              <h2 className="text-base font-medium mb-4">Billing</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Period</Label>
                  <div className="flex gap-2">
                    <Input type="number" defaultValue="1" className="w-20" />
                    <Select defaultValue="month">
                      <SelectTrigger>
                        <SelectValue placeholder="Select period" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="month">Month</SelectItem>
                        <SelectItem value="3months">3 Months</SelectItem>
                        <SelectItem value="year">Year</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <div className="space-y-2">
                  <Label>Next Payment</Label>
                  <Button
                    variant="outline"
                    className="w-full justify-start text-left font-normal"
                    onClick={() => setDate(new Date())}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4" />
                    {date ? format(date, "PPP") : <span>Pick a date</span>}
                  </Button>
                </div>

                <div className="space-y-4">
                <div className="space-y-2">
                  <Label>Payment Method</Label>
                  <Select>
                    <SelectTrigger>
                      <SelectValue placeholder="Select a Payment Method" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="credit-card">Credit Card</SelectItem>
                      <SelectItem value="debit-card">Debit Card</SelectItem>
                      <SelectItem value="paypal">PayPal</SelectItem>
                    </SelectContent>
                  </Select>
                  <Button variant="link" className="h-auto p-0 text-xs">
                    Edit Payment Methods in settings
                  </Button>
                </div>
                <div className="space-y-2">
                  <Label>Notes</Label>
                  <Textarea placeholder="Add any additional notes here..." className="min-h-[100px]" />
                </div>
              </div>
              </div>
            </div>

            {/* Reminders Section */}
            <div className={CARD_PARENT_CLASS}>
              <h2 className="text-base font-medium mb-4">Reminders</h2>
              <p className="text-sm text-muted-foreground mb-4">
                Set up to 3 reminders - you must enter a next payment date
              </p>
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label>Reminder 1</Label>
                  <Select>
                    <SelectTrigger>
                      <SelectValue placeholder="None" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="none">None</SelectItem>
                      <SelectItem value="1-day">1 day before</SelectItem>
                      <SelectItem value="3-days">3 days before</SelectItem>
                      <SelectItem value="1-week">1 week before</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>Reminder 2</Label>
                  <Select>
                    <SelectTrigger>
                      <SelectValue placeholder="None" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="none">None</SelectItem>
                      <SelectItem value="1-day">1 day before</SelectItem>
                      <SelectItem value="3-days">3 days before</SelectItem>
                      <SelectItem value="1-week">1 week before</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>Reminder 3</Label>
                  <Select>
                    <SelectTrigger>
                      <SelectValue placeholder="None" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="none">None</SelectItem>
                      <SelectItem value="1-day">1 day before</SelectItem>
                      <SelectItem value="3-days">3 days before</SelectItem>
                      <SelectItem value="1-week">1 week before</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </div>

            <Button className="w-full text-white hover:bg-black/90">
              Add Subscription
            </Button>
          </div>
        </div>

      </form>
    </Form>
  )
}

