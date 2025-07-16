/**
 * Utility function to convert props object to HTML attributes string
 * Maps object entries to key="value" format and joins with spaces
 */
export function buildAttributes(props: Record<string, any>): string {
  return Object.entries(props)
    .map(([key, value]) => `${key}="${value}"`)
    .join(" ");
}