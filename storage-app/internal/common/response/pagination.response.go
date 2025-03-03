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
