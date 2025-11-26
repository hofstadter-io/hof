package dagger

import (
	"context"

	"dagger.io/dagger"
)

var client *dagger.Client

func Get(ctx context.Context) (*dagger.Client, error) {
	if client != nil {
		return client, nil
	}
	// try connecting to dagger
	var err error
	client, err = dagger.Connect(ctx)
	return client, err
}
