package service

import (
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/infrastructure/file"
)

type IFile interface {
	Upload(filepath string, bytes []byte) status.Object
	Delete(filepath string) status.Object
	GetURL(assetType file.AssetType, filename, format string) string
}
