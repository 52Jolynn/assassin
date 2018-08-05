package core

type Pagination struct {
	Limit     int         `json:"limit"`
	Offset    int         `json:"offset"`
	Total     int         `json:"total"`
	PageNum   int         `json:"page_num"`
	TotalPage int         `json:"total_page"`
	Data      interface{} `json:"data"`
}

func NewPagination(limit, offset, total int, data interface{}) *Pagination {
	pagination := Pagination{Limit: limit, Offset: offset, Total: total, Data: data}
	if total % limit == 0 {
		pagination.TotalPage = total / limit
	} else {
		pagination.TotalPage = total / limit + 1
	}
	pagination.PageNum = offset / limit + 1
	return &pagination
}
