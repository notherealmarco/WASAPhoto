package helpers

import (
	"net/url"
	"strconv"
)

const (
	DEFAULT_LIMIT  = 30
	DEFAULT_OFFSET = 0
)

// Get the start index and limit from the query.
// If they are not present, use the default values.
func GetLimits(query url.Values) (int, int, error) {

	limit := DEFAULT_LIMIT
	start_index := DEFAULT_OFFSET

	var err error

	if query.Get("limit") != "" {
		limit, err = strconv.Atoi(query.Get("limit"))
	}

	if query.Get("start_index") != "" {
		start_index, err = strconv.Atoi(query.Get("start_index"))
	}

	if err != nil {
		return 0, 0, err
	}

	return start_index, limit, nil
}
