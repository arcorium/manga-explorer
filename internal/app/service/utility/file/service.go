package file

import (
	"manga-explorer/internal/app/common"
)

type IService interface {
	Upload(filepath string, bytes []byte) common.Status
	Delete(filepath string) common.Status
	GetURL(assetType AssetType, filename, format string) string
}
