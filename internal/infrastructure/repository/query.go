package repository

import (
  "github.com/uptrace/bun"
  "manga-explorer/internal/util"
  "strconv"
)

func NewResult[T any, U util.Integral](data T, total U) PagedQueryResult[T] {
  return PagedQueryResult[T]{Data: data, Total: uint64(total)}
}

type PagedQueryResult[T any] struct {
  Data  T
  Total uint64
}

var NoQueryParameter = QueryParameter{0, 0} // Used for to get all of them without offset and limit

type QueryParameter struct {
  Offset uint64
  Limit  uint64 // Mostly ignored when the item is less than the Limit, it should not be used to check error
}

func (p QueryParameter) Insert(query *bun.SelectQuery) *bun.SelectQuery {
  return query.Offset(int(p.Offset)).Limit(int(p.Limit))
}

func (p QueryParameter) RawQuery() string {
  str := " OFFSET " + strconv.Itoa(int(p.Offset))

  if p.Limit != 0 {
    str += " LIMIT " + strconv.Itoa(int(p.Limit))
  }
  return str
}
