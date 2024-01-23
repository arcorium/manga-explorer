package file

import (
	"errors"
	"fmt"
	"manga-explorer/internal/util"
	"strings"
)

type AssetType string

const (
	MangaAsset   AssetType = "mangas"
	ProfileAsset           = "profiles"
	UnknownAsset           = ""
)

var ErrNoFormat = errors.New("file has no format")

func ParseFileFormat(filename string) (Format, error) {
	split := strings.Split(filename, ".")
	if len(split) <= 1 {
		return FormatUnknown, ErrNoFormat
	}
	return Format(split[len(split)-1]), nil
}

type Format string

const (
	FormatJPG     Format = "jpg"
	FormatJPEG           = "jpeg"
	FormatPNG            = "png"
	FormatUnknown        = ""
)

func (f Format) String() string {
	return string(f)
}

func (f Format) Filename(name string) Name {
	return Name(fmt.Sprintf("%s.%s", name, f.String()))
}

func (f Format) Validate() bool {
	return util.IsOneOf(f, FormatJPG, FormatJPEG, FormatPNG)
}

func (a AssetType) String() string {
	return string(a)
}
