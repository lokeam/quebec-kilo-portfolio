package igdb

import (
	"fmt"
	"strings"

	"github.com/lokeam/qko-beta/internal/interfaces"
)

const (
	MinLimit = 1
	MaxLimit = 50
)

type IGDBQuery struct {
	search  string
	fields  []string
	where   []string
	limit   int
}

type QueryBuilder struct {
	query IGDBQuery
	logger interfaces.Logger
}

// NewIGDBQueryBuilder creates a new query builder
func NewIGDBQueryBuilder(logger interfaces.Logger) *QueryBuilder {
	igdbQueryBuilder := &QueryBuilder{
		// Init slices to prevent nil pointer dereferences
		query: IGDBQuery{
			fields: make([]string, 0),
			where:  make([]string, 0),
		},
		logger: logger,
	}

	if logger != nil {
    logger.Debug("Created new QueryBuilder", map[string]any{
        "builder": igdbQueryBuilder,
    })
	}

	return igdbQueryBuilder
}

// Search sets search term for the query.
// Returns QueryBuilder for method chaining.
func (qb *QueryBuilder) Search(searchTerm string) *QueryBuilder {
	if qb.logger != nil {
		qb.logger.Debug("Query builder setting search term", map[string]any{
			"searchTerm": searchTerm,
		})
	}
	qb.query.search = searchTerm
	return qb
}

// Fields sets fields that will be returned by the query.
// Returns QueryBuilder for method chaining.
func (qb *QueryBuilder) Fields(fields ...string) *QueryBuilder {
	if qb.logger != nil {
		qb.logger.Debug("Query builder setting fields", map[string]any{
			"fields": fields,
		})
	}
	qb.query.fields = fields
	return qb
}

// Where adds condition to the query.
// Current use case is for filtering by game_type
// Returns QueryBuilder for method chaining.
func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	if qb.logger != nil {
		qb.logger.Debug("Query building adding where condition", map[string]any{
			"conditions": condition,
		})
	}
	qb.query.where = append(qb.query.where, condition)
	return qb
}

// Limit sets max number of results returned by query
// Returns QueryBuilder for method chaining.
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	if qb.logger != nil {
		qb.logger.Debug("Query builder setting limit", map[string]any{
			"limit": limit,
		})
	}
	qb.query.limit = limit
	return qb
}

// Build constructs the final IGDB API query string.
// Returns an error if query is invalid
func (qb *QueryBuilder) Build() (string, error) {
	// Validate fields
	if len(qb.query.fields) == 0 {
		return "", ErrNoFields
	}

	// Validate search term
	if qb.query.search == "" {
		return "", ErrInvalidSearchTerm
	}

	// Validate limit
	if qb.query.limit < MinLimit || qb.query.limit > MaxLimit {
		return "", NewInvalidLimitError(qb.query.limit)
	}

	// Validate where conditions
	if len(qb.query.where) == 0 {
		return "", ErrInvalidWhereCondition
	}

	// Validate logger
	if qb.logger != nil {
		qb.logger.Debug("Building query", map[string]any{
			"query": qb.query,
		})
	}

	queryParts := qb.buildQueryParts()
	query := strings.Join(queryParts, " ")

	// Validate
	if qb.logger != nil {
		qb.logger.Debug("IGDB Query built successfully", map[string]any{
			"query": query,
		})
	}

	return query, nil
}

// Helper fn - buildQueryParts constructs the query parts and returns them as a slice of strings
func (qb *QueryBuilder) buildQueryParts() []string {
	var queryParts []string

	// Add fields
	if len(qb.query.fields) > 0 {
		queryParts = append(queryParts, fmt.Sprintf("fields %s;", strings.Join(qb.query.fields, ",")))
	}

	// Add search
	if qb.query.search != "" {
		queryParts = append(queryParts, fmt.Sprintf("search \"%s\";", qb.query.search))
	}

	// Add where conditions
	if len(qb.query.where) > 0 {
		queryParts = append(queryParts, fmt.Sprintf("where %s;", strings.Join(qb.query.where, " & ")))
	}

	// Add limit
	if qb.query.limit > 0 {
		queryParts = append(queryParts, fmt.Sprintf("limit %d;", qb.query.limit))
	}

	return queryParts
}