package environ

import (
	"context"
	"fmt"
	"os"
	"sync"

	"dagger.io/dagger"
	"gorm.io/gorm"
)

type localEnviron struct {
	mx  sync.RWMutex
	ctx context.Context
	dag *dagger.Client
	db  *gorm.DB
}

var LE *localEnviron

const DAGGER_HOST = "container://veg-dagger-engine"

func Initialize(ctx context.Context, db *gorm.DB) (err error) {
	os.Setenv("_EXPERIMENTAL_DAGGER_RUNNER_HOST", DAGGER_HOST)
	LE = &localEnviron{
		ctx: ctx,
		db:  db,
	}
	err = LE.AutoMigrate()
	if err != nil {
		return err
	}
	LE.dag, err = dagger.Connect(ctx)
	return err
}

func Client() *localEnviron {
	if LE == nil {
		panic("EnvironmentService(dagger) has not been initialized")
	}
	return LE
}

func (le *localEnviron) ListEnvirons() ([]tableEnviron, error) {

	var envs []tableEnviron
	err := le.db.WithContext(le.ctx).
		Model(&tableEnviron{}).
		Preload("Children").
		Where("tag =?", "genesis").
		Find(&envs).Error
	if err != nil {
		return nil, fmt.Errorf("while fetching environs from database: %w", err)
	}

	return envs, nil
}

type listEnvironTagsResponse struct {
}

func (le *localEnviron) ListEnvironTags(envUri string) ([]listEnvironTagsResponse, error) {
	return nil, nil
}
