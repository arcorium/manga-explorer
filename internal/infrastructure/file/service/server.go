package service

import (
	"io/fs"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/infrastructure/file"
	"os"
)

func NewLocalFileService(dir string) IFile {
	_ = os.Mkdir(dir, fs.ModeDir)
	return &serverFileService{
		Directory: dir,
	}
}

type serverFileService struct {
	Directory string
}

func (s serverFileService) GetURL(assetType file.AssetType, filename, format string) string {
	return s.Directory + "/" + assetType.String() + "/" + filename + "." + format
}

func (s serverFileService) Upload(filepath string, bytes []byte) status.Object {
	file, err := os.Create(filepath)
	if err != nil {
		return status.Error(status.FILE_UPLOAD_FAILED)
	}

	_, err = file.Write(bytes)
	if err != nil {
		return status.Error(status.FILE_UPLOAD_FAILED)
	}

	return status.Success(status.CREATED)
}

func (s serverFileService) Delete(filepath string) status.Object {
	err := os.Remove(filepath)
	return status.Conditional(err, status.INTERNAL_SERVER_ERROR)

}
