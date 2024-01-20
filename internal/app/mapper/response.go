package mapper

import (
	"manga-explorer/internal/app/dto"
	"math"
)

func NewResponsePage[T any](result []T, total uint64, query *dto.PagedQueryInput) dto.ResponsePage {
	return dto.ResponsePage{
		Elements:      uint64(len(result)),
		TotalElements: total,
		TotalPage:     uint64(math.Ceil(float64(total) / float64(query.Element))),
		CurrentPage:   uint64(math.Ceil(float64(query.Offset()) / float64(query.Element))),
	}
}
