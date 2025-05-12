package types

// GameType represents a game type with its display and normalized text
type GameType struct {
	ID             int    `json:"id"`
  Type           string `json:"type"`
	DisplayText    string `json:"displayText"`
	NormalizedText string `json:"normalizedText"`
}

type GameTypeResponse struct {
	DisplayText    string `json:"displayText"`
	NormalizedText string `json:"normalizedText"`
}

// GameTypes maps IGDB game type IDs to their display and normalized text
var GameTypes = map[int]GameType{
	0:  {DisplayText: "Main Game", NormalizedText: "main"},
	1:  {DisplayText: "DLC Addon", NormalizedText: "dlc"},
	2:  {DisplayText: "Expansion", NormalizedText: "expansion"},
	3:  {DisplayText: "Bundle", NormalizedText: "bundle"},
	4:  {DisplayText: "Standalone Expansion", NormalizedText: "standalone"},
	5:  {DisplayText: "Mod", NormalizedText: "mod"},
	6:  {DisplayText: "Episode", NormalizedText: "episode"},
	7:  {DisplayText: "Season", NormalizedText: "season"},
	8:  {DisplayText: "Remake", NormalizedText: "remake"},
	9:  {DisplayText: "Remaster", NormalizedText: "remaster"},
	10: {DisplayText: "Expanded Game", NormalizedText: "expanded"},
	11: {DisplayText: "Port", NormalizedText: "port"},
	12: {DisplayText: "Fork", NormalizedText: "fork"},
	13: {DisplayText: "Pack", NormalizedText: "pack"},
	14: {DisplayText: "Update", NormalizedText: "update"},
}

// GetGameType returns the GameType for a given ID, or a zero value if not found
func GetGameType(id int) GameType {
	if gameType, exists := GameTypes[id]; exists {
		return gameType
	}
	return GameType{} // Return zero value if not found
}
