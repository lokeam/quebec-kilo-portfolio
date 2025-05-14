package types

type IGDBQuery struct {
	search  string
	fields  []string
	where   []string
	limit   int
}