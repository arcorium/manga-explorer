package file

import "manga-explorer/internal/util/opt"

type Name string

func (f Name) String() string {
  return string(f)
}

var NullName = opt.Null[Name]()

var NoFile = Name(" ")
