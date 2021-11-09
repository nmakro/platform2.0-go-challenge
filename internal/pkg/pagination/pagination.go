package pagination

// PageInfo struct.
type PageInfo struct {
	Page       uint32 `json:"page"`
	PageSize   uint32 `json:"per_page"`
	TotalPages uint32 `json:"total_pages,omitempty"`
	TotalItems uint32 `json:"total_items,omitempty"`
}

// PageInfoRequest : Pangeinfo from a request.
type PageInfoRequest struct {
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"per_page"`
}

// PageInfoResponse : Pageinfo for a response.
type PageInfoResponse struct {
	TotalPages uint32 `json:"total_pages,omitempty"`
	TotalItems uint32 `json:"total_items,omitempty"`
}

// Offset returns the offset based on page and pagesize.
func (p *PageInfo) Offset() uint32 {
	return (p.Page - 1) * p.PageSize
}

// GetOrDefault checks for zero page or pagesize and replaces with default.
func (p *PageInfo) GetOrDefault(defaultPage, defaultSize uint32) PageInfo {
	page := defaultPage
	pageSize := defaultSize

	if p.Page > 0 {
		page = p.Page
	}

	if p.PageSize > 0 {
		pageSize = p.PageSize
	}

	return PageInfo{
		Page:     page,
		PageSize: pageSize,
	}
}

// NewPageInfo constructor also calculates the amount or total pages.
func NewPageInfo(page, pageSize, totalCount uint32) PageInfo {
	if pageSize == 0 {
		pageSize = 1
	}
	tp := totalCount / pageSize
	if totalCount%pageSize != 0 {
		tp++
	}

	return PageInfo{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: tp,
		TotalItems: totalCount,
	}
}

// CalculatePagination calculates the amount or total pages.
func CalculatePagination(request PageInfo, totalCount uint32) PageInfo {
	return NewPageInfo(request.Page, request.PageSize, totalCount)
}

// SliceSplitIndexes returnes the start and end index of  a selected page of a slice.
func SliceSplitIndexes(pi PageInfo, sliceLength int) (int, int) {
	start := int((pi.Page - 1) * pi.PageSize)
	if start >= sliceLength {
		return 0, 0
	}

	end := start + int(pi.PageSize)

	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}
