package common

import (
	"golang.org/x/text/language"
)

type Language string

func (l Language) Validate() bool {
	_, err := language.Parse(string(l))
	if err != nil {
		return false
	}
	return true
}

// Parse convert the language into ISO 3 code as a string, it will return empty string ("") when the language
// is unknown
func (l Language) Parse() string {
	tag, err := language.Parse(string(l))
	if err != nil {
		return ""
	}
	base, confidence := tag.Base()
	if confidence == language.No {
		return ""
	}

	return base.ISO3()
}

// ParseLang wrapper for Parse but returning Language instead of string
func (l Language) ParseLang() Language {
	return Language(l.Parse())
}

func NewLanguage(lang string) Language {
	return Language(lang).ParseLang()
}
