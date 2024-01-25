package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"manga-explorer/internal/common"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/infrastructure/file"
	"manga-explorer/internal/util"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func NewLocalFileService(config *common.Config, host, endpoint, dir string, routes gin.IRouter) IFile {
	dir = filepath.Dir(dir + "/") // Append '/' so /file would do
	err := os.MkdirAll(dir, fs.ModePerm)
	util.DoNothing(err)

	for _, asset := range util.SliceWrap(file.MangaAsset, file.ProfileAsset, file.CoverAsset) {
		path := filepath.Join(dir, asset.String())
		err = os.MkdirAll(path, fs.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("Failed to create directory: %s", err))
		}
	}
	// Serve static
	routes.Static(endpoint, dir)

	return &serverFileService{
		Directory: dir,
		endpoint:  fmt.Sprintf("%s/%s", host, strings.TrimPrefix(endpoint, "/")),
	}
}

type serverFileService struct {
	endpoint  string
	Directory string
}

func (s serverFileService) getLocalPath(types file.AssetType, filename file.Name) string {
	return fmt.Sprintf("./%s/%s/%s", s.Directory, types.String(), filename)
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
		return "", status.Error(status.BAD_REQUEST_ERROR)
	}
	// TODO: Response size each image

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

func (s serverFileService) Endpoint(types file.AssetType) string {
	return fmt.Sprintf("%s/%s", s.endpoint, types.String())
}
func (s serverFileService) GetFullpath(assetType file.AssetType, filename file.Name) string {
	if len(filename) == 0 || filename == file.NoFile {
		return ""
	}
	return fmt.Sprintf("%s/%s", s.Endpoint(assetType), filename)
}
