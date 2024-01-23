package service

import (
	"fmt"
	"io"
	"io/fs"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/infrastructure/file"
	"manga-explorer/internal/util"
	"mime/multipart"
	"os"
	"path/filepath"
)

func NewLocalFileService(dir string, localhost string) IFile {
	_ = os.MkdirAll(filepath.Dir(dir), fs.ModeDir)
	return &serverFileService{
		Directory: dir,
		LocalHost: localhost,
	}
}

type serverFileService struct {
	Directory string
	LocalHost string
}

func (s serverFileService) getLocalPath(types file.AssetType, filename file.Name) string {
	return fmt.Sprintf("/%s/%s/%s", s.Directory, types.String(), filename)
}

func (s serverFileService) Upload(types file.AssetType, fileHeader *multipart.FileHeader) (file.Name, status.Object) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", status.InternalError()
	}
	defer src.Close()

	// Get and validate format
	format, err := file.ParseFileFormat(fileHeader.Filename)
	if err != nil {
		return "", status.Error(status.BAD_BODY_REQUEST_ERROR)
	}
	// TODO: Validate size each image

	// Make new filename and append the format
	filename := format.Filename(util.GenerateRandomString(30))
	localPath := s.getLocalPath(types, filename)
	// Save file
	dst, err := os.Create(localPath)
	if err != nil {
		return "", status.InternalError()
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", status.InternalError()
	}

	return filename, status.Created()
}

func (s serverFileService) Uploads(types file.AssetType, files []multipart.FileHeader) ([]file.Name, status.Object) {
	filenames := []file.Name{}
	for _, fl := range files {
		filename, stat := s.Upload(types, &fl)
		if stat.IsError() {
			return nil, stat
		}
		filenames = append(filenames, filename)
	}

	return filenames, status.Success()
}

func (s serverFileService) Delete(types file.AssetType, filename file.Name) status.Object {
	// Get local path
	localPath := s.getLocalPath(types, filename)
	err := os.Remove(localPath)
	if err != nil {
		return status.InternalError()
	}
	return status.Success(status.DELETED)

}
