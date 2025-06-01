package domain

type PageInfo struct {
	HasNextPage     bool
	HasPreviousPage bool
	StartCursor     *PageCursor
	EndCursor       *PageCursor
}

func NewPageInfo(hasNextPage, hasPreviousPage bool, startCursor, endCursor *PageCursor) *PageInfo {
	return &PageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		StartCursor:     startCursor,
		EndCursor:       endCursor,
	}
}
