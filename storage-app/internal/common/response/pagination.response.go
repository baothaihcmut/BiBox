package response

type PaginationResponse struct {
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	Total      int64 `json:"total"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
	NextOffset *int  `json:"next_offset"`
	PrevOffset *int  `json:"prev_offset"`
}

func InitPaginationResponse(count, limit, offset int) PaginationResponse {
	pagination := PaginationResponse{
		Offset:  offset,
		Limit:   limit,
		Total:   int64(count),
		HasNext: false,
		HasPrev: false,
	}

	if offset+limit < int(count) {
		nextOffset := offset + limit
		pagination.HasNext = true
		pagination.NextOffset = &nextOffset
	}
	if offset > 0 {
		prevOffset := offset - limit
		pagination.HasPrev = true
		pagination.PrevOffset = &prevOffset
	}
	return pagination
}
