package environ

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"dagger.io/dagger"
)

const VEG_ENVIRONMENT_REGISTRY = "host.docker.internal:5000"

type Environ struct {
	Eid string `gorm:"primaryKey"`
	Tag string `gorm:"primaryKey"`

	Name string
	Uri  string `gorm:"index"`

	SrcUri  string
	SrcPath string
	FromUri string
	DstPath string

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`

	// Has-Many relationship: env can have many children-env.
	// Children []Environ `gorm:"foreignKey:From;references:Uri"`

	// Has-Many relationship: A session has many events.
	// Sessions []storageEvent `gorm:"foreignKey:AppName,UserID,SessionID;references:AppName,UserID,ID"`
}

func (Environ) TableName() string {
	return "environs"
}

// AutoMigrate runs the GORM auto-migration tool to ensure the database schema
// matches the internal storage models (e.g., storageSession, storageEvent).
//
// NOTE: This function relies on a type assertion to the concrete *databaseService
// implementation. It will return an error if the provided session.Service is
// a different implementation.
func (le *localEnviron) AutoMigrate() error {
	err := le.db.AutoMigrate(&Environ{})
	if err != nil {
		return fmt.Errorf("environment auto migrate failed: %w", err)
	}
	return nil
}

// LookupEnviron resolves an environment URI to its metadata and a Dagger container.
// Input: envUri (e.g. oci://host/img:tag or id:tag)
// Output: (Environ metadata, *dagger.Container, error)
func (le *localEnviron) LookupEnviron(envUri string) (Environ, *dagger.Container, error) {
	var foundEnv Environ

	if !strings.Contains(envUri, "://") {
		envUri = "oci://" + envUri
	}
	e, err := url.Parse(envUri)
	if err != nil {
		return foundEnv, nil, fmt.Errorf("lookup.parse.error: %w", err)
	}

	// fmt.Printf("LOOKUP: %s: %#+v\n", envUri, e)
	// key := fmt.Sprintf("%s%s", e.Host, e.Path)

	// if not in our registry, then it is a generic image, just return container
	if e.Host != VEG_ENVIRONMENT_REGISTRY {
		cleanUri := strings.TrimPrefix(envUri, "oci://")
		return foundEnv, le.dag.Container().From(cleanUri), nil
	}

	// extract eid and tag
	// key ~ host.docker.internal:5000/id:ver
	pathParts := strings.Split(e.Path, "/")
	if len(pathParts) < 2 {
		// shouldn't happen for our registry, but just in case
		return foundEnv, le.dag.Container().From(envUri), nil
	}

	imgParts := strings.Split(pathParts[1], ":")
	if len(imgParts) < 2 {
		return foundEnv, nil, fmt.Errorf("lookup.error: image must have a tag: %q", envUri)
	}

	eid, tag := imgParts[0], imgParts[1]

	err = le.db.WithContext(le.ctx).
		Where(&Environ{
			Eid: eid,
			Tag: tag,
		}).
		First(&foundEnv).Error

	if err != nil {
		// If not found in database, still return a container if it's in our registry?
		// Actually, if it's in our registry, it *should* be in our database.
		// But let's be safe and return a container if we can.
		cleanUri := strings.TrimPrefix(envUri, "oci://")
		return foundEnv, le.dag.Container().From(cleanUri), nil
	}

	// Use the URI from the database if available
	cleanUri := strings.TrimPrefix(foundEnv.Uri, "oci://")
	env := le.dag.Container().From(cleanUri)

	return foundEnv, env, nil
}

// persistEnviron publishes a Dagger container to the registry and updates the database.
// Input: envUri (clean OCI reference), tEnv (initial metadata), c (the Dagger container)
func (le *localEnviron) persistEnviron(envUri string, tEnv *Environ, c *dagger.Container) (err error) {

	// fmt.Println("le.persist.input", envUri)

	// publish to persist
	cleanUri := envUri
	cleanUri = strings.TrimPrefix(cleanUri, "oci://")
	_, err = c.Publish(le.ctx, cleanUri)
	if err != nil {
		return fmt.Errorf("while persisting environment to registry(%s|%s): %w", envUri, cleanUri, err)
	}

	// fmt.Println("le.persist.published", true)

	// extract eid:tag envUri
	qparts := strings.Split(cleanUri, "?")
	parts := strings.Split(qparts[0], "/")
	img := parts[len(parts)-1]
	iparts := strings.Split(img, ":")

	// fmt.Println("le.persist.vars", qparts, parts, img, iparts)

	// what about from?
	if tEnv == nil {
		tEnv = &Environ{}
	}
	tEnv.Eid = iparts[0]
	tEnv.Tag = iparts[1]
	tEnv.Uri = cleanUri

	// fmt.Printf("tEnv: %#+v\n", *tEnv)

	// save to database
	err = le.db.WithContext(le.ctx).
		Save(tEnv).Error
	if err != nil {
		return fmt.Errorf("while persisting environment to database(%s): %w", envUri, err)
	}

	// fmt.Println("le.persist.database", true)

	return nil
}

// func replaceTag(envUri, nextTag string) string {
// 	// preserve any query params
// 	qparts := strings.Split(envUri, "?")

// 	// replace tag
// 	parts := strings.Split(qparts[0], ":")
// 	parts[len(parts)-1] = nextTag

// 	// preserve any query params
// 	if len(qparts) > 1 {
// 		parts = append(parts, qparts[1:]...)
// 	}

// 	// return a new Uri
// 	return strings.Join(parts, "")
// }

// func extractPath(envUri string) (string, error) {
// 	qparts := strings.Split(envUri, "?")
// 	if len(qparts) < 2 {
// 		return "", fmt.Errorf("missing query params in Uri to extact path in: %q", envUri)
// 	}
// 	vals, err := url.ParseQuery(qparts[1])
// 	if err != nil {
// 		return "", fmt.Errorf("while parsing query params in: %q, %w", envUri, err)
// 	}
// 	path := vals.Get("path")
// 	if path == "" {
// 		return "", fmt.Errorf("empty path in: %q", envUri)
// 	}
// 	return path, nil
// }

// func extractPathEmptyOk(envUri string) (string, error) {
// 	qparts := strings.Split(envUri, "?")
// 	if len(qparts) < 2 {
// 		return "", fmt.Errorf("missing query params in Uri to extact path in: %q", envUri)
// 	}
// 	vals, err := url.ParseQuery(qparts[1])
// 	if err != nil {
// 		return "", fmt.Errorf("while parsing query params in: %q, %w", envUri, err)
// 	}
// 	path := vals.Get("path")
// 	return path, nil
// }

func (le *localEnviron) getEnvironEntry(envUri string) (Environ, error) {
	var foundEnv Environ
	e, err := url.Parse(envUri)
	if err != nil {
		return foundEnv, fmt.Errorf("database error while fetching environ: %w", err)
	}
	key := fmt.Sprintf("%s%s", e.Host, e.Path)
	// fmt.Println("le.lookup", e, key)

	err = le.db.WithContext(le.ctx).
		Where(&Environ{
			Uri: key,
		}).
		First(&foundEnv).Error

	if err != nil {
		// For any error including ErrRecordNotFound, return it as a system error.
		return foundEnv, fmt.Errorf("database error while fetching environ: %w", err)
	}

	return foundEnv, nil
}

func (le *localEnviron) listEnvironEntry() ([]Environ, error) {
	var foundEnvs []Environ

	err := le.db.WithContext(le.ctx).
		Find(&foundEnvs).Error

	if err != nil {
		// For any error including ErrRecordNotFound, return it as a system error.
		return foundEnvs, fmt.Errorf("database error while listing environs: %w", err)
	}

	return foundEnvs, nil
}
