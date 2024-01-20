package common

import "github.com/biter777/countries"

type Country string

func (c Country) Validate() bool {
	return countries.ByName(string(c)) == countries.Unknown
}

func NewCountry(code countries.CountryCode) Country {
	return Country(code.Alpha2())
}

type Language string

func (l Language) Validate() bool {
	return countries.ByName(string(l)) != countries.Unknown
}

func NewLanguage(code countries.CountryCode) Language { return Language(code.Alpha3()) }
