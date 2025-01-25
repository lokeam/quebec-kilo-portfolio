import { format } from 'date-fns';

export function formatDate(dateString: string | undefined): { dayStr: string; monthStr: string } {  // Updated return type
  if (!dateString) {
    return { dayStr: '--', monthStr: '---' };
  }

  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      return { dayStr: '--', monthStr: '---' };
    }
    return {
      dayStr: format(date, 'd'),    // Changed from day to dayStr
      monthStr: format(date, 'MMM') // monthStr stays the same
    };
  } catch {
    return { dayStr: '--', monthStr: '---' };
  }
}