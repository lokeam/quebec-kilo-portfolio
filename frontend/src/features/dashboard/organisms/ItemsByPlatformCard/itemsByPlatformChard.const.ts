import { type ChartConfig } from "@/shared/components/ui/chart";

export const chartConfig = {
  itemCount: {
    label: "Visitors",
  },
  ps4: {
    label: "PS4",
    color: "hsl(var(--chart-1))",
  },
  pc: {
    label: "PC",
    color: "hsl(var(--chart-2))",
  },
  rom: {
    label: "ROM",
    color: "hsl(var(--chart-3))",
  },
  xbox: {
    label: "XBOX",
    color: "hsl(var(--chart-4))",
  },
  switch: {
    label: "Switch",
    color: "hsl(var(--chart-5))",
  },
} satisfies ChartConfig;
