package types

type IGDBQuery struct {
	Search  string
	Fields  []string
	Where   []string
	Limit   int
}