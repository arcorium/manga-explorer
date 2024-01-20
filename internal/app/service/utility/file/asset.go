package file

type AssetType string

const (
	ImageAsset AssetType = "images"
	AudioAsset           = "audios"
)

func (a AssetType) String() string {
	return string(a)
}
