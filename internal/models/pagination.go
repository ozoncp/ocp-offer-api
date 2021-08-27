package models

import "math"

type PaginationInfo struct {
	HasPreviousPage bool   // Has previous page
	HasNextPage     bool   // Has next page
	PerPage         uint32 // Items per page - max 10k
	Page            uint64 // Current page number
	TotalPages      uint64 // Total pages
	TotalItems      uint64 // Total items
}

type PaginationInput struct {
	// Number of items per page
	Take uint32
	// The number of skipped elements
	Skip uint64
}

// Get pagination information.
func (p *PaginationInput) GetPaginationInfo(perPage uint32, totalItems uint64) *PaginationInfo {
	totalPages := uint64(math.Ceil(float64(totalItems) / float64(p.Take)))
	page := uint64(math.Ceil(float64(p.Skip)/float64(p.Take) + 1))

	hasPreviousPage := false
	if page > 1 {
		hasPreviousPage = true
	}

	hasNextPage := false
	if page < totalPages {
		hasNextPage = true
	}

	return &PaginationInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		PerPage:         perPage,
		Page:            page,
		TotalPages:      totalPages,
		TotalItems:      totalItems,
	}
}
