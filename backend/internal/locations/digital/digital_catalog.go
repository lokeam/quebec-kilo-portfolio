package digital

// DigitalServicesCatalog contains the complete list of available digital services.
// This is a read-only list and should not be modified at runtime.
var DigitalServicesCatalog = []DigitalServiceCatalogItem{
	{
		ID:                    "amazonluna",
		Name:                  "Amazon Luna",
		Logo:                  "amazon",
		IsSubscriptionService: true,
	},
	{
		ID:                    "applearcade",
		Name:                  "Apple Arcade",
		Logo:                  "apple",
		IsSubscriptionService: true,
	},
	{
		ID:                    "blizzard",
		Name:                  "Blizzard Battle.net",
		Logo:                  "blizzard",
		IsSubscriptionService: false,
	},
	{
		ID:                    "ea",
		Name:                  "EA Play",
		Logo:                  "ea",
		IsSubscriptionService: true,
	},
	{
		ID:                    "epicgames",
		Name:                  "Epic Games",
		Logo:                  "epicgames",
		IsSubscriptionService: false,
	},
	{
		ID:                    "fanatical",
		Name:                  "Fanatical",
		Logo:                  "fanatical",
		IsSubscriptionService: false,
	},
	{
		ID:                    "gog",
		Name:                  "GOG",
		Logo:                  "gog",
		IsSubscriptionService: false,
	},
	{
		ID:                    "googleplaypass",
		Name:                  "Google Play Pass",
		Logo:                  "google",
		IsSubscriptionService: true,
	},
	{
		ID:                    "greenmangaming",
		Name:                  "Green Man Gaming",
		Logo:                  "greenmangaming",
		IsSubscriptionService: false,
	},
	{
		ID:                    "humblebundle",
		Name:                  "Humble Bundle",
		Logo:                  "humblebundle",
		IsSubscriptionService: false,
	},
	{
		ID:                    "itchio",
		Name:                  "itch.io",
		Logo:                  "itchio",
		IsSubscriptionService: false,
	},
	{
		ID:                    "meta",
		Name:                  "Meta",
		Logo:                  "meta",
		IsSubscriptionService: false,
	},
	{
		ID:                    "netflix",
		Name:                  "Netflix",
		Logo:                  "netflix",
		IsSubscriptionService: true,
	},
	{
		ID:                    "nintendo",
		Name:                  "Nintendo",
		Logo:                  "nintendo",
		IsSubscriptionService: true,
	},
	{
		ID:                    "nvidia",
		Name:                  "NVIDIA",
		Logo:                  "nvidia",
		IsSubscriptionService: true,
	},
	{
		ID:                    "primegaming",
		Name:                  "Prime Gaming",
		Logo:                  "prime",
		IsSubscriptionService: true,
	},
	{
		ID:                    "playstation",
		Name:                  "PlayStation Network",
		Logo:                  "ps",
		IsSubscriptionService: true,
	},
	{
		ID:                    "shadow",
		Name:                  "Shadow",
		Logo:                  "shadow",
		IsSubscriptionService: true,
	},
	{
		ID:                    "steam",
		Name:                  "Steam",
		Logo:                  "steam",
		IsSubscriptionService: false,
	},
	{
		ID:                    "ubisoft",
		Name:                  "Ubisoft",
		Logo:                  "ubisoft",
		IsSubscriptionService: false,
	},
	{
		ID:                    "xboxlive",
		Name:                  "Xbox Live",
		Logo:                  "xbox",
		IsSubscriptionService: true,
	},
	{
		ID:                    "xboxgamepass",
		Name:                  "Xbox Game Pass",
		Logo:                  "xbox",
		IsSubscriptionService: true,
	},
}
