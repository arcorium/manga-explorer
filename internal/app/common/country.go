package common

import "github.com/biter777/countries"

type Country string

func (c Country) Validate() bool {
	return c.Code() != countries.Unknown
}

func (c Country) Code() countries.CountryCode {
	return countries.ByName(string(c))
}

func NewCountry(code countries.CountryCode) Country {
	return Country(code.Alpha2())
}

type Language string

func (l Language) Validate() bool {
	return l.Code() != countries.Unknown
}

func (l Language) Code() countries.CountryCode {
	return countries.ByName(string(l))
}

func NewLanguage(code countries.CountryCode) Language { return Language(code.Alpha3()) }
