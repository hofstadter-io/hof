package docker

import (
	"context"
	"github.com/docker/docker/api/types"
)

func GetVersion() (string, types.Version, error) {
	clientVer := dockerClient.ClientVersion()	
	serverVer, err := dockerClient.ServerVersion(context.Background())
	return clientVer, serverVer, err
}
