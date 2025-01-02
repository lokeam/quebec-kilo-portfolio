
import type { OnlineService, OnlineServicesData } from './onlineServicesCard.types';

export const onlineServices: OnlineService[] = [
  {
    name: "Steam",
    url: "https://store.steampowered.com/",
    monthlyFee: "FREE",
    annualFee: "FREE",
    plan: "",
    renewalMonth: "",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "GOG",
    url: "https://www.gog.com/en",
    monthlyFee: "FREE",
    annualFee: "FREE",
    plan: "",
    renewalMonth: "",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Nintendo Switch Online",
    url: "https://accounts.nintendo.com/",
    monthlyFee: "$6.67/month",
    annualFee: "$80.04/year",
    plan: "12 Month",
    renewalMonth: "Mar 2025",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Playstation Network",
    url: "https://account.sonyentertainmentnetwork.com/",
    monthlyFee: "$15/month",
    annualFee: "$180/year",
    plan: "3 Month",
    renewalMonth: "Dec 2024",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "XBOX Live",
    url: "https://www.xbox.com/en-US/live",
    monthlyFee: "$10/month",
    annualFee: "$120/year",
    plan: "1 Month",
    renewalMonth: "Aug 2024",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Epic Games Store",
    url: "https://www.epicgames.com/store/en-US/",
    monthlyFee: "FREE",
    annualFee: "FREE",
    plan: "",
    renewalMonth: "",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Uplay",
    url: "https://uplay.ubisoft.com/",
    monthlyFee: "FREE",
    annualFee: "FREE",
    plan: "",
    renewalMonth: "",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Origin",
    url: "https://www.origin.com/",
    monthlyFee: "FREE",
    annualFee: "FREE",
    plan: "",
    renewalMonth: "",
    avatar: "/placeholder.svg?height=36&width=36"
  },
];

export const onlineServicesData: OnlineServicesData = {
  // Calculate total annual fees from non-free services
  totalAnnual: "$380.04", // Sum of Nintendo ($80.04), PSN ($180), and Xbox ($120)

  // Assuming current month is March 2024, these services renew this month
  renewsThisMonth: ["Nintendo Switch Online", "XBOX Live"],

  // Total count of services
  totalServices: onlineServices.length,

  // All services
  services: onlineServices
};
