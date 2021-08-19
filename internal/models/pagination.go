package models

type PaginationInfo struct {
	Page             uint64    // Current page number
	TotalPages       uint64    // Total pages
	TotalItems       uint64    // Total items
	PerPage          uint32    // Items per page - max 10k
	HasNextPage      bool      // Has next page
	HasPreviousPage  bool      // Has previous page
}

type PaginationInput struct {
	Cursor    uint64
	Take      uint32
	Skip      uint64
}
