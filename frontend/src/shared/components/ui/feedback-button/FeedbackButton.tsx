import { useState } from 'react';

// Shadcn UI components
import { Button } from '@/shared/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/shared/components/ui/dialog';
import { Label } from '@/shared/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';

// Icons
import { MessageSquare, Send } from '@/shared/components/ui/icons';

// Custom hooks - Sentry feedback service
import { submitFeedback, type FeedbackData } from '@/core/monitoring/sentry-feedback-service';
import { useSentryTracking } from '@/shared/hooks/useSentryTracking';

// Utils
import { cn } from '@/shared/components/ui/utils';

interface FeedbackButtonProps {
  variant?: 'default' | 'outline' | 'ghost';
  size?: 'default' | 'sm' | 'lg' | 'icon';
  className?: string;
}

type FeedbackType = 'bug' | 'feature' | 'general' | 'other';
type UserExperience = 'excellent' | 'good' | 'fair' | 'poor';

export function FeedbackButton({
  variant = 'outline',
  size = 'default',
  className
}: FeedbackButtonProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [feedbackType, setFeedbackType] = useState<FeedbackType>('general');
  const [message, setMessage] = useState('');
  const [userExperience, setUserExperience] = useState<UserExperience | undefined>(undefined);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { trackUserInteraction } = useSentryTracking();

  const handleSubmit = async () => {
    if (!message.trim()) return;

    setIsSubmitting(true);
    const startTime = performance.now();

    try {
      // Submit feedback to Sentry
      const feedback: FeedbackData = {
        type: feedbackType,
        message: message.trim(),
        userExperience,
        currentPage: window.location.pathname,
      };

      submitFeedback(feedback);

      // Track user interactions
      const duration = performance.now() - startTime;
      trackUserInteraction('feedback_submitted', duration, {
        feedbackType,
        messageLength: message.length,
        hasUserExperience: !!userExperience,
      });

      // Reset form and close dialog
      setMessage('');
      setFeedbackType('general');
      setUserExperience(undefined);
      setIsOpen(false);
    } catch (error) {
      console.error('Failed to submit feedback:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleOpenChange = (open: boolean) => {
    setIsOpen(open);
    if (open) {
      trackUserInteraction('feedback_dialog_opened', undefined, {
        currentPage: window.location.pathname,
      });
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button variant={variant} size={size} className={className}>
          <MessageSquare className="h-4 w-4 mr-2" />
          Give feedback or Report an Issue
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Help us improve Q-Ko</DialogTitle>
          <DialogDescription>
            Tell us what went wrong or share your feedback. We'll review it and get back to you.
          </DialogDescription>
        </DialogHeader>

        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="feedback-type">What type of feedback is this?</Label>
            <Select value={feedbackType} onValueChange={(value: FeedbackType) => setFeedbackType(value)}>
              <SelectTrigger>
                <SelectValue placeholder="Select feedback type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="bug">Bug Report</SelectItem>
                <SelectItem value="feature">Feature Request</SelectItem>
                <SelectItem value="general">General Feedback</SelectItem>
                <SelectItem value="other">Other</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="message">Your feedback</Label>
            <textarea
              id="message"
              placeholder="Describe what happened or what you'd like to see improved..."
              value={message}
              onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setMessage(e.target.value)}
              rows={4}
              className={cn(
                "flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-none"
              )}
            />
          </div>

          <div className="grid gap-2">
            <Label htmlFor="experience">How would you rate your experience? (Optional)</Label>
            <Select value={userExperience || ''} onValueChange={(value: string) => setUserExperience(value as UserExperience || undefined)}>
              <SelectTrigger>
                <SelectValue placeholder="Select rating" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="excellent">Excellent</SelectItem>
                <SelectItem value="good">Good</SelectItem>
                <SelectItem value="fair">Fair</SelectItem>
                <SelectItem value="poor">Poor</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => setIsOpen(false)}
            disabled={isSubmitting}
          >
            Cancel
          </Button>
          <Button
            onClick={handleSubmit}
            disabled={!message.trim() || isSubmitting}
          >
            {isSubmitting ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                Sending...
              </>
            ) : (
              <>
                <Send className="h-4 w-4 mr-2" />
                Send Feedback
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}