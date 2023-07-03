package container

import (
	"context"
	"time"
)

func GetVersion() (RuntimeVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return rt.Version(ctx)
}
