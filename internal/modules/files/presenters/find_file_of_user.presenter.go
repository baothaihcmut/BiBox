package presenters

type FindFileOfUserInput struct {
	IsInFolder *bool  `form:"is_in_folder"`
	IsFolder   *bool  `form:"is_folder"`
	SortBy     string `form:"sort_by" bind:"required"`
	IsAsc      bool   `form:"is_asc" bind:"required"`
	Offset     int    `form:"offset" bind:"required"`
	Limit      int    `form:"limit" bind:"required"`
}

type FindFileOfUserOuput struct {
	Files []*FileOutput `json:"files"`
}
