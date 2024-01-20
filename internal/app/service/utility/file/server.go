package file

import (
	"io/fs"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"os"
)

func NewLocalFileService(dir string) IService {
	_ = os.Mkdir(dir, fs.ModeDir)
	return &serverFileService{
		Directory: dir,
	}
}

type serverFileService struct {
	Directory string
}

func (s serverFileService) GetURL(assetType AssetType, filename, format string) string {
	return s.Directory + "/" + assetType.String() + "/" + filename + "." + format
}

func (s serverFileService) Upload(filepath string, bytes []byte) common.Status {
	file, err := os.Create(filepath)
	if err != nil {
		return common.StatusError(status.FILE_UPLOAD_FAILED)
	}

	_, err = file.Write(bytes)
	if err != nil {
		return common.StatusError(status.FILE_UPLOAD_FAILED)
	}

	return common.StatusSuccess(status.SUCCESS_CREATED)
}

func (s serverFileService) Delete(filepath string) common.Status {
	err := os.Remove(filepath)
	return common.ConditionalStatus(err, status.INTERNAL_SERVER_ERROR)
}
