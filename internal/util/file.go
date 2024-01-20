package util

import (
	"errors"
	"strings"
)

var NoFormatErr = errors.New("no format found")

func GetFileFormat(filename string) (string, error) {
	split := strings.Split(filename, ".")
	if len(split) <= 1 {
		return "", NoFormatErr
	}
	return split[len(split)-1], nil
}
