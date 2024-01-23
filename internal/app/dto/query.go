package dto

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/infrastructure/repository"
	"manga-explorer/internal/util"
	"strconv"
)

type PagedQueryInput struct {
	Element uint64 `form:"element"`
	Page    uint64 `form:"page"`
}

func (p PagedQueryInput) ConstructQuery(ctx *gin.Context) {
	p.Element = util.DropError(strconv.ParseUint(ctx.Query("element"), 10, 64))
	p.Page = util.DropError(strconv.ParseUint(ctx.Query("page"), 10, 64))
}

func (p PagedQueryInput) Offset() uint64 {
	return (p.Page - 1) * p.Element
}

func (p PagedQueryInput) IsPaged() bool {
	return p.Element == 0
}

func (p PagedQueryInput) ToQueryParam() repository.QueryParameter {
	if p.IsPaged() {
		return repository.NoQueryParameter
	}
	return repository.QueryParameter{
		Offset: p.Offset(),
		Limit:  p.Element,
	}
}
