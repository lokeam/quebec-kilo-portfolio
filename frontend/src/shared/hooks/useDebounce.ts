import { useEffect, useState } from "react";
import { debounce } from "lodash";

export function useDebounce<T>(
  value: T,
  delay: number
): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = debounce(() => {
      setDebouncedValue(value);
    }, delay);

    handler();

    return () => handler.cancel();
  }, [value, delay]);

  return debouncedValue;
}
