package digital

import (
	"github.com/lokeam/qko-beta/internal/types"
)

// DigitalServicesCatalog contains the complete list of available digital services.
// This is a read-only list and should not be modified at runtime.
var DigitalServicesCatalog = []types.DigitalServiceItem{
	{
		ID:                    "amazonluna",
		Name:                  "Amazon Luna",
		Logo:                  "amazon",
		IsSubscriptionService: true,
		URL:                   "https://luna.amazon.com/",
	},
	{
		ID:                    "applearcade",
		Name:                  "Apple Arcade",
		Logo:                  "apple",
		IsSubscriptionService: true,
		URL:                   "https://www.apple.com/apple-arcade/",
	},
	{
		ID:                    "blizzard",
		Name:                  "Blizzard Battle.net",
		Logo:                  "blizzard",
		IsSubscriptionService: false,
		URL:                   "https://www.blizzard.com/en-us/",
	},
	{
		ID:                    "ea",
		Name:                  "EA Play",
		Logo:                  "ea",
		IsSubscriptionService: true,
		URL:                   "https://www.ea.com/ea-play",
	},
	{
		ID:                    "epicgames",
		Name:                  "Epic Games",
		Logo:                  "epicgames",
		IsSubscriptionService: false,
		URL:                   "https://store.epicgames.com/en-US/",
	},
	{
		ID:                    "fanatical",
		Name:                  "Fanatical",
		Logo:                  "fanatical",
		IsSubscriptionService: false,
		URL:                   "https://www.fanatical.com/en/",
	},
	{
		ID:                    "gog",
		Name:                  "GOG",
		Logo:                  "gog",
		IsSubscriptionService: false,
		URL:                   "https://www.gog.com/en/",
	},
	{
		ID:                    "googleplaypass",
		Name:                  "Google Play Pass",
		Logo:                  "google",
		IsSubscriptionService: true,
		URL:                   "https://play.google.com/store/pass/getstarted/",
	},
	{
		ID:                    "greenmangaming",
		Name:                  "Green Man Gaming",
		Logo:                  "greenmangaming",
		IsSubscriptionService: false,
		URL:                   "https://www.greenmangaming.com/",
	},
	{
		ID:                    "humblebundle",
		Name:                  "Humble Bundle",
		Logo:                  "humblebundle",
		IsSubscriptionService: false,
		URL:                   "https://www.humblebundle.com/",
	},
	{
		ID:                    "itchio",
		Name:                  "itch.io",
		Logo:                  "itchio",
		IsSubscriptionService: false,
		URL:                   "https://itch.io/",
	},
	{
		ID:                    "meta",
		Name:                  "Meta",
		Logo:                  "meta",
		IsSubscriptionService: false,
		URL:                   "https://www.meta.com/nz/meta-quest-plus/",
	},
	{
		ID:                    "nintendo",
		Name:                  "Nintendo Switch Online",
		Logo:                  "nintendo",
		IsSubscriptionService: true,
		URL:                   "https://www.nintendo.com/",
	},
	{
		ID:                    "nvidia",
		Name:                  "NVIDIA",
		Logo:                  "nvidia",
		IsSubscriptionService: true,
		URL:                   "https://www.nvidia.com/en-us/geforce-now/",
	},
	{
		ID:                    "primegaming",
		Name:                  "Prime Gaming",
		Logo:                  "prime",
		IsSubscriptionService: true,
		URL:                   "https://gaming.amazon.com/home",
	},
	{
		ID:                    "playstation",
		Name:                  "PlayStation Network",
		Logo:                  "ps",
		IsSubscriptionService: true,
		URL:                   "https://www.playstation.com/en-us/playstation-network/",
	},
	{
		ID:                    "shadow",
		Name:                  "Shadow",
		Logo:                  "shadow",
		IsSubscriptionService: true,
		URL:                   "https://shadow.tech/",
	},
	{
		ID:                    "steam",
		Name:                  "Steam",
		Logo:                  "steam",
		IsSubscriptionService: false,
		URL:                   "https://store.steampowered.com/",
	},
	{
		ID:                    "ubisoft",
		Name:                  "Ubisoft",
		Logo:                  "ubisoft",
		IsSubscriptionService: false,
		URL:                   "https://www.ubisoft.com/en-us/",
	},
	{
		ID:                    "xboxgamepass",
		Name:                  "Xbox Game Pass",
		Logo:                  "xbox",
		IsSubscriptionService: true,
		URL:                   "https://www.xbox.com/en-US/xbox-game-pass",
	},
}
