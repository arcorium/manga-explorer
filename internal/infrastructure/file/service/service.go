package service

import (
  "manga-explorer/internal/common/status"
  "manga-explorer/internal/infrastructure/file"
  "mime/multipart"
)

type IFile interface {
  Upload(types file.AssetType, header *multipart.FileHeader) (file.Name, status.Object)
  Uploads(types file.AssetType, header []multipart.FileHeader) ([]file.Name, status.Object) // TODO: Handle when there is an error in the middle of uploading
  Delete(types file.AssetType, filename file.Name) status.Object
  Endpoint(assetType file.AssetType) string
  GetFullpath(assetType file.AssetType, filename file.Name) string
}
