package environ

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
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

func (le *localEnviron) ListEnvirons() ([]Environ, error) {

	var envs []Environ
	err := le.db.WithContext(le.ctx).
		Model(&Environ{}).
		Preload("Children").
		Where("tag =?", "0").
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

func IncrementTag(envUri string) (nextUri, nextTag string, err error) {
	// preserve any query params
	qparts := strings.Split(envUri, "?")

	// handle potential encoding in the OCI reference
	// if we don't have enough colons, try unescaping
	if strings.Count(qparts[0], ":") < 3 {
		decoded, err := url.PathUnescape(qparts[0])
		if err == nil {
			qparts[0] = decoded
		}
	}

	// replace tag
	parts := strings.Split(qparts[0], ":")
	currTag := parts[len(parts)-1]
	currInt, err := strconv.ParseInt(currTag, 10, 64)
	if err != nil {
		return "", "", err
	}
	nextTag = fmt.Sprintf("%d", currInt+1)
	parts[len(parts)-1] = nextTag
	nextUri = strings.Join(parts, ":")

	// preserve any query params
	if len(qparts) > 1 {
		newParts := append([]string{nextUri}, qparts[1:]...)
		nextUri = strings.Join(newParts, "?")
	}

	return nextUri, nextTag, nil
}

func ReplaceTag(envUri, nextTag string) string {
	// preserve any query params
	qparts := strings.Split(envUri, "?")

	// handle potential encoding in the OCI reference
	if strings.Count(qparts[0], ":") < 3 {
		decoded, err := url.PathUnescape(qparts[0])
		if err == nil {
			qparts[0] = decoded
		}
	}

	// replace tag
	parts := strings.Split(qparts[0], ":")
	parts[len(parts)-1] = nextTag
	newUri := strings.Join(parts, ":")

	// preserve any query params
	if len(qparts) > 1 {
		newParts := append([]string{newUri}, qparts[1:]...)
		newUri = strings.Join(newParts, "?")
	}

	return newUri
}
