// Tailwind Form classes
export const CARD_PARENT_CLASS = "rounded-lg border bg-card text-card-foreground shadow-sm relative cursor-pointer w-full p-4 bg-gradient-to-b from-slate-900 to-slate-950 border-slate-800" as const;

export const BADGE_STYLES = {
  mediaType: {
    hardware: "bg-green-700/50 text-slate-200",
    dlc: "bg-orange-700/50 text-slate-200",
    inGamePurchase: "bg-blue-600/50 text-slate-200",
    disc: "bg-blue-400/50 text-slate-200",
    physical: "bg-yellow-400/50 text-slate-200",
    subscription: "bg-red-800/50 text-slate-200"
  },
  spendTransactionType: {
    subscription: "bg-purple-900/50 text-purple-200",
    "one-time": "bg-slate-700/50 text-slate-200"
  }
} as const;
