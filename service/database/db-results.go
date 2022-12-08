package database

type QueryResult int

// Constants used to represent the result of queries
const (
	SUCCESS       = 0
	ERR_NOT_FOUND = 1
	ERR_EXISTS    = 2
	ERR_INTERNAL  = 3
)
