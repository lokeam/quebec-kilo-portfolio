/**
 * Utility function to combine class names (similar to cn utility from clsx/classnames)
 * Filters out falsy values and joins the remaining classes with spaces
 */
export function cn(...classes: (string | undefined | null | false)[]): string {
  return classes.filter(Boolean).join(" ");
}