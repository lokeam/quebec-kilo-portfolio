
/**
 * Colors and names used for for graphs on the:
 *
 * TotalMonthlySpendingCard (Spend Tracking)
 * MonthlySpendingCard (Dashboard)
 *
 * Centralized here to avoid duplication and ensure consistency
 *
 */
export const GRAPH_CATEGORY_COLORS: Record<string, string> = {
  hardware: "var(--color-hardware)",
  dlc: "var(--color-dlc)",
  in_game_purchase: "var(--color-in-game-purchase)",
  one_time_purchase: "var(--color-one-time-purchase)",
  subscription: "var(--color-subscription)",
  physical_game: "var(--color-physical)",
  digital_game: "var(--color-digital)",
  misc: "var(--color-misc)",
};

export const GRAPH_CATEGORY_DISPLAY_NAMES: Record<string, string> = {
  hardware: "Hardware",
  dlc: "DLC",
  in_game_purchase: "In-game Purchase",
  one_time_purchase: "One-time Purchase",
  subscription: "Subscription",
  physical_game: "Physical Game",
  digital_game: "Digital Game",
  misc: "Misc",
};