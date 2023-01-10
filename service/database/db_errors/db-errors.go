package db_errors

import "strings"

// Returns true if the error is a "no rows in result set" error
func EmptySet(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no rows in result set")
}

// Returns true if the error is a Unique constraint violation error
func UniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

// Returns true if the error is a Foreign Key constraint violation error
func ForeignKeyViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "FOREIGN KEY constraint failed")
}
