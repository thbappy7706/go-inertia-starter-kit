package models

type PaginatedResponse[T any] struct {
	Data        []T   `json:"data"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	From        int   `json:"from"`
	To          int   `json:"to"`
}

func NewPaginatedResponse[T any](data []T, total int64, page, perPage int) PaginatedResponse[T] {
	if perPage <= 0 {
		perPage = 10
	}
	if page <= 0 {
		page = 1
	}

	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}

	from := 0
	to := 0
	if total > 0 {
		from = (page-1)*perPage + 1
		to = from + len(data) - 1
		if to > int(total) {
			to = int(total)
		}
	}

	return PaginatedResponse[T]{
		Data:        data,
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     perPage,
		Total:       total,
		From:        from,
		To:          to,
	}
}
