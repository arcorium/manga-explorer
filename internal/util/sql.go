package util

import (
  "database/sql"
)

func CheckSqlResult(result sql.Result, err error) error {
  if err != nil {
    return err
  }
  affectedRow, err := result.RowsAffected()
  if err != nil {
    return err
  }
  if affectedRow <= 0 {
    return sql.ErrNoRows
  }
  return nil
}

type Result[T any] struct {
  Data T
  Err  error
}

func (r Result[T]) Unwrap() (T, error) {
  return r.Data, r.Err
}

func CheckSliceResult[T any](result []T, err error) Result[[]T] {
  if err != nil {
    return Result[[]T]{nil, err}
  }

  if len(result) <= 0 {
    return Result[[]T]{nil, sql.ErrNoRows}
  }
  return Result[[]T]{result, nil}
}
