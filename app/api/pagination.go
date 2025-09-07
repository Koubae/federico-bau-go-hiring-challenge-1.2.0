package api

import (
	"errors"
	"net/http"
	"strconv"
)

const (
	DefaultLimit = 10
	MaxLimit     = 100
	MinLimit     = 1

	DefaultOffset = 0
	MinOffset     = 0
)

type Pagination struct {
	Limit  int
	Offset int
}

func PaginationRequest(r *http.Request) (*Pagination, error) {
	limit := DefaultLimit
	offset := DefaultOffset

	if limitQuery := r.URL.Query().Get("limit"); limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err != nil {
			return nil, errors.New("invalid limit parameter: must be a number")
		}
		limit = parsedLimit
	}

	if offsetQuery := r.URL.Query().Get("offset"); offsetQuery != "" {
		parsedOffset, err := strconv.Atoi(offsetQuery)
		if err != nil {
			return nil, errors.New("invalid offset parameter: must be a number")
		}
		offset = parsedOffset
	}

	if limit < MinLimit {
		return nil, errors.New("invalid limit parameter: must be greater than 0")
	} else if limit > MaxLimit {
		return nil, errors.New("invalid limit parameter: must be less than " + strconv.Itoa(MaxLimit))
	}
	if offset < MinOffset {
		return nil, errors.New("invalid offset parameter: must be greater than or equal to 0")
	}

	pagination := &Pagination{
		Limit:  limit,
		Offset: offset,
	}
	return pagination, nil
}
