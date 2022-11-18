package db_errors

import "strings"

// Returns true if the query result has no rows
func EmptySet(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no rows in result set")
}
