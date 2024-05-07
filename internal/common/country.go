package common

import (
  "github.com/biter777/countries"
)

type Country string

func (c Country) Validate() bool {
  return c.Code() != countries.Unknown
}

func (c Country) Code() countries.CountryCode {
  return countries.ByName(string(c))
}
