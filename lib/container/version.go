package container

import (
	"context"

	"github.com/docker/docker/api/types"
)

func GetVersion() (string, types.Version, error) {
	clientVer := cClient.ClientVersion()
	serverVer, err := cClient.ServerVersion(context.Background())
	return clientVer, serverVer, err
}
