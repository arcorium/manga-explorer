package mapper

import (
	dto2 "manga-explorer/internal/common/dto"
	"math"
)

func NewResponsePage[T any](result []T, total uint64, query *dto2.PagedQueryInput) dto2.ResponsePage {
	return dto2.ResponsePage{
		Elements:      uint64(len(result)),
		TotalElements: total,
		TotalPage:     uint64(math.Ceil(float64(total) / float64(query.Element))),
		CurrentPage:   uint64(math.Ceil(float64(query.Offset()) / float64(query.Element))),
	}
}
