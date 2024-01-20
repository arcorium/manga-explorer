package fixtures

import "github.com/google/uuid"

type Type int

const (
	UserType Type = iota
	MangaType
)

var uuids = []string{
	uuid.NewString(),
	uuid.NewString(),
	uuid.NewString(),
	uuid.NewString(),
}

func GetConstantUUID(index uint8) string {
	if index > 3 {
		panic("You MORON!")
	}

	return uuids[index]
}
