package dto

type PagedQueryInput struct {
	Element uint64 `form:"element"`
	Page    uint64 `form:"page"`
}

func (p PagedQueryInput) Offset() uint64 {
	return (p.Page - 1) * p.Element
}

func (p PagedQueryInput) IsPaged() bool {
	return p.Element == 0
}
