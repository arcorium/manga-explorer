package mapper

import (
	"manga-explorer/internal/common/dto"
	"math"
)

func NewResponsePage[T any](result []T, total uint64, query *dto.PagedQueryInput) dto.ResponsePage {
	if query.Element == 0 {
		return dto.ResponsePage{
			Elements:      uint64(len(result)),
			TotalElements: total,
		}
	}
	return dto.ResponsePage{
		Elements:      uint64(len(result)),
		TotalElements: total,
		TotalPage:     uint64(math.Ceil(float64(total) / float64(query.Element))),
		CurrentPage:   (query.Offset() / query.Element) + 1,
	}
}
