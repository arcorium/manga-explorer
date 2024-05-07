package containers

import (
  "errors"
)

type FilterFunc[T any] func(current *T) bool

func SliceFilter[T any](data []T, filterFunc FilterFunc[T]) []T {
  if data == nil {
    return nil
  }

  result := make([]T, 0, len(data))

  for i := range data {
    if filterFunc(&data[i]) {
      result = append(result, data[i])
    }
  }

  return result
}

func SliceContains[T comparable](data []T, expected T) bool {
  for _, v := range data {
    if v == expected {
      return true
    }
  }
  return false
}

func SliceIndexOf[T any](data []T, filterFunc FilterFunc[T]) (int, bool) {
  for i, v := range data {
    if filterFunc(&v) {
      return i, true
    }
  }
  return -1, false
}

func CombineSplices[T any](datas ...[]T) []T {
  if len(datas) == 0 {
    return nil
  }

  res := []T{}
  for i := 0; i < len(datas); i++ {
    res = append(res, datas[i]...)
  }
  return res
}

func SliceRemoveDuplicates[T comparable](slice []T) []T {
  if slice == nil {
    return nil
  }

  result := make([]T, 0, len(slice))
  uniques := make(map[T]struct{})

  for _, s := range slice {
    if _, ok := uniques[s]; !ok {
      result = append(result, s)
      uniques[s] = struct{}{}
    }
  }

  return result
}

func SliceCountDuplicates[T comparable](slice []T) map[T]uint64 {
  if slice == nil {
    return nil
  }

  result := make(map[T]uint64)

  for _, s := range slice {
    if _, ok := result[s]; !ok {
      result[s] = 1
    } else {
      result[s]++
    }
  }

  return result
}

type EqualFunc[T any] func(current T) bool
type EqualFuncMap[K comparable, V any] func(key K, val V) bool

type EqualFunc2[T any, U any] func(haystack T, needle U) bool

func MapContains[K comparable, V any](haystack map[K]V, equalFunc EqualFuncMap[K, V]) bool {
  if haystack == nil {
    return false
  }

  for k, v := range haystack {
    if equalFunc(k, v) {
      return true
    }
  }
  return false
}

func MapValues[K comparable, V any](maps map[K]V) []V {
  if maps == nil {
    return nil
  }

  arr := make([]V, 0, len(maps))
  for _, v := range maps {
    arr = append(arr, v)
  }

  return arr
}

func MapKeys[K comparable, V any](maps map[K]V) []K {
  if maps == nil {
    return nil
  }

  arr := make([]K, 0, len(maps))
  for k := range maps {
    arr = append(arr, k)
  }

  return arr
}

type SafeConvertFunc[From any, To any] func(current From) (To, error)
type CastFunc[From any, To any] func(current From) To

type ConvertFuncEnumerated[From any, To any] func(index int, current From) To

// SafeCastSlice Used to convert []From into []Recipients based on function parameter with error checking
func SafeCastSlice[From any, To any](slice []From, convertFunc SafeConvertFunc[*From, To]) ([]To, error) {
  if slice == nil {
    return nil, errors.New("parameter is nil")
  }

  result := make([]To, 0, len(slice))
  for _, val := range slice {
    cur, err := convertFunc(&val)
    if err != nil {
      return result, err
    }
    result = append(result, cur)
  }
  return result, nil
}

// CastSlice Used to convert []From into []Recipients based on function parameter without error checking, use it when the object conversion never fail and if the slice argument is nil it will return empty slice not nil
func CastSlice[From, To any](slice []From, convertFunc CastFunc[From, To]) []To {
  result := make([]To, 0, len(slice))
  if slice == nil || len(slice) == 0 {
    return result
  }

  for _, val := range slice {
    result = append(result, convertFunc(val))
  }
  return result
}

// CastSlicePtr works like CastSlice method, but instead of passing From type it will pass *From type in the function parameter
func CastSlicePtr[From, To any](slice []From, convertFunc CastFunc[*From, To]) []To {
  result := make([]To, 0, len(slice))
  if slice == nil || len(slice) == 0 {
    return result
  }

  for _, val := range slice {
    result = append(result, convertFunc(&val))
  }
  return result
}

type CastParamFunc[From, To, Param any] func(current From, param Param) To

// CastSlicePtr1 works like CastSlice but instead it will use the params for each function
func CastSlicePtr1[From any, To any, Param any](slice []From, params Param, convertFunc CastParamFunc[*From, To, Param]) []To {
  result := make([]To, 0, len(slice))
  if slice == nil || len(slice) == 0 {
    return result
  }

  for i := 0; i < len(slice); i++ {
    result = append(result, convertFunc(&slice[i], params))
  }

  return result
}
