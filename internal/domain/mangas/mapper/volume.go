package mapper

import (
	"github.com/google/uuid"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/util/containers"
)

func ToVolumeResponse(volume *mangas.Volume) dto.VolumeResponse {
	return dto.VolumeResponse{
		Title:    volume.Title,
		Number:   volume.Number,
		Chapters: containers.CastSlicePtr(volume.Chapters, ToChapterResponse),
	}
}

func MapVolumeCreateInput(input *dto.VolumeCreateInput) mangas.Volume {
	return mangas.Volume{
		Id:      uuid.NewString(),
		MangaId: input.MangaId,
		Title:   input.Title,
		Number:  input.Number,
	}
}
