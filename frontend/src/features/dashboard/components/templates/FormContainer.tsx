import { type ReactNode } from 'react';

// Components
import { Form } from '@/shared/components/ui/form';

// Hooks
import { type UseFormReturn } from 'react-hook-form';

// Zod
import { z } from 'zod';

interface FormContainerProps<T extends z.ZodType> {
  form: UseFormReturn<z.infer<T>>;
  onSubmit: (data: z.infer<T>) => void;
  children: ReactNode;
}

export function FormContainer<T extends z.ZodType>({
  form,
  onSubmit,
  children,
}: FormContainerProps<T>) {
  return (
    /* Pass all useForm hook props to Shadcn Form component */
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8 p-4">
        {children}
      </form>
    </Form>
  );
}
