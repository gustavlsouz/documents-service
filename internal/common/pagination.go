package common

const MaxPageSize = 50

type Pagination interface {
	Page() int
	Size() int
	SetPage(page int)
	SetSize(limit int)
}

type pagination struct {
	page  int
	limit int
}

func noZero(v int) int {
	if v == 0 {
		return 1
	}
	return v
}

func max(v int, m int) int {
	if v > m {
		return m
	}
	return v
}

func NewPagination() Pagination {
	return &pagination{page: 1, limit: 1}
}

func (p *pagination) Page() int {
	return p.page
}

func (p *pagination) SetPage(page int) {
	p.page = noZero(page)
}

func (p *pagination) Size() int {
	return p.limit
}

func (p *pagination) SetSize(limit int) {
	p.limit = max(noZero(limit), MaxPageSize)
}
