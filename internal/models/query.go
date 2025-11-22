package models

// Range filter: date, timestamp, numeric ranges etc
type RangeFilter struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Fuzzy filter: LIKE %search%
type FuzzyFilter struct {
	Query     string  `json:"query"`
	Threshold float64 `json:"threshold,omitempty"`
}

// Advanced Query Payload
type QueryRequest struct {
	Filters map[string]interface{} `json:"filters,omitempty"`
	Range   map[string]RangeFilter `json:"range,omitempty"`
	Fuzzy   map[string]FuzzyFilter `json:"fuzzy,omitempty"`
	Limit   int                    `json:"limit,omitempty"`
	Offset  int                    `json:"offset,omitempty"`
	Sort    string                 `json:"sort,omitempty"`
	Format  string                 `json:"format,omitempty"`
}
