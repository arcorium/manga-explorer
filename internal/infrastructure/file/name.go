package file

import (
	"fmt"
)

type Name string

func (f Name) Fullpath(webUrl string, assetType AssetType) string {
	if len(f) == 0 {
		return ""
	}
	return fmt.Sprintf("%s%s", webUrl, string(f))
}

func (f Name) HostnameFullpath(assetType AssetType) string {
	return f.Fullpath(hostName, assetType)
}

func (f Name) String() string {
	return string(f)
}

var hostName string

func GetHostName() string {
	return hostName
}

func SetHostName(hostname string) {
	hostName = hostname
}
