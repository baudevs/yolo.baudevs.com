package license

// RateLimit represents rate limiting configuration
type RateLimit struct {
	MaxQueries int   `json:"max_queries"` // Maximum number of queries
	Period     int64 `json:"period"`      // Period in seconds
}
