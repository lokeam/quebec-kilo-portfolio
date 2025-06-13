package types

// GameType represents a game type with its display and normalized text
type GameType struct {
	ID             int64    `json:"id"`
  Type           string `json:"type"`
	DisplayText    string `json:"display_text"`
	NormalizedText string `json:"normalized_text"`
}

type GameTypeResponse struct {
	DisplayText    string `json:"display_text"`
	NormalizedText string `json:"normalized_text"`
}

// GameTypes maps IGDB game type IDs to their display and normalized text
var GameTypes = map[int64]GameType{
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

