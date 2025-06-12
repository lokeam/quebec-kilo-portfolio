/**
 * Truncates a deletion warning text to a maximum length and adds an ellipsis if the text is longer.
 * @param text - The text to truncate.
 * @param maxLength - The maximum length of the text.
 * @returns The truncated text.
 */
export const truncateWarningText = (text: string, maxLength: number = 20): string => {
  if (text.length <= maxLength) return text;
  return `${text.slice(0, maxLength)}...`
};
