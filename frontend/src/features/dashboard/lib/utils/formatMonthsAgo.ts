/**
 * Formats a timestamp into a human-readable "months ago" string
 * @param timestamp - The timestamp in milliseconds
 * @returns A string in the format "from X months ago"
 */

export function formatMonthsAgo(timestamp: number): string {
  const now = new Date();
  const date = new Date(timestamp);

  // Check if millisecond timestamp is within the current month
  if (now.getFullYear() === date.getFullYear() && now.getMonth() === date.getMonth()) {
    return 'from this month';
  }

  const monthsDiff = (now.getFullYear() - date.getFullYear()) * 12 +
    (now.getMonth() - date.getMonth());

  // Check if timestamp is over a year old
  if (monthsDiff > 11) {
    return 'from over a year ago';
  }

  return `from ${monthsDiff} month${monthsDiff === 1 ? '' : 's'} ago`;
}