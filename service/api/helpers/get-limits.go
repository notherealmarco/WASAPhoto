package helpers

import (
	"net/url"
	"strconv"
)

const (
	DEFAULT_LIMIT  = 15 // don't know if should be moved to config
	DEFAULT_OFFSET = 0
)

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
