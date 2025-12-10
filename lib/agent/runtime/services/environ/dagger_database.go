package environ

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"dagger.io/dagger"
)

const VEG_ENVIRONMENT_REGISTRY = "host.docker.internal:5000"

type tableEnviron struct {
	Eid string `gorm:"primaryKey"`
	Tag string `gorm:"primaryKey"`

	Name string
	Uri  string `gorm:"index"`
	Src  string `gorm:"index"`
	From string `gorm:"index"`
	Base string

	CreateAt time.Time `gorm:"index"`
	UpdateAt time.Time `gorm:"index"`

	// Has-Many relationship: env can have many children-env.
	Children []tableEnviron `gorm:"foreignKey:From;references:Uri"`

	// Has-Many relationship: A session has many events.
	// Sessions []storageEvent `gorm:"foreignKey:AppName,UserID,SessionID;references:AppName,UserID,ID"`
}

func (tableEnviron) TableName() string {
	return "environs"
}

// AutoMigrate runs the GORM auto-migration tool to ensure the database schema
// matches the internal storage models (e.g., storageSession, storageEvent).
//
// NOTE: This function relies on a type assertion to the concrete *databaseService
// implementation. It will return an error if the provided session.Service is
// a different implementation.
func (le *localEnviron) AutoMigrate() error {
	err := le.db.AutoMigrate(&tableEnviron{})
	if err != nil {
		return fmt.Errorf("environment auto migrate failed: %w", err)
	}
	return nil
}

func (le *localEnviron) lookupEnviron(envUri string) (tableEnviron, *dagger.Container, error) {
	var foundEnv tableEnviron
	// fucking more hacks because our paths / URIs are a mess...
	// we are seeing veg://... here, which is not correct, we should never see that in the server, it is a vscode thing only!
	if !strings.Contains(envUri, "://") {
		envUri = "oci://" + envUri
	}
	e, err := url.Parse(envUri)
	if err != nil {
		return foundEnv, nil, fmt.Errorf("database error while fetching environ: %w", err)
	}
	// fmt.Printf("LOOKUP: %s: %#+v\n", envUri, e)
	key := fmt.Sprintf("%s%s", e.Host, e.Path)
	// fmt.Println("le.lookupEnviron.key", key)
	// hacky, error potential, but we should only be getting internal reps here anyway
	parts := strings.Split(strings.Split(key, "/")[1], ":")
	eid, tag := parts[0], parts[1]

	err = le.db.WithContext(le.ctx).
		Where(&tableEnviron{
			Eid: eid,
			Tag: tag,
		}).
		First(&foundEnv).Error

	if err != nil {
		// For any error including ErrRecordNotFound, return it as a system error.
		return foundEnv, nil, fmt.Errorf("database error while fetching environ: %w", err)
	}

	// fmt.Printf("le.lookupEnviron.table %v\n", foundEnv)

	env := le.dag.Container().From(foundEnv.Uri)

	return foundEnv, env, nil
}

func (le *localEnviron) persistEnviron(envUri string, tEnv *tableEnviron, c *dagger.Container) (err error) {

	// fmt.Println("le.persist.input", envUri)

	// publish to persist
	_, err = c.Publish(le.ctx, envUri)
	if err != nil {
		return fmt.Errorf("while persisting environment to registry(%s): %w", envUri, err)
	}

	// fmt.Println("le.persist.published", true)

	// extract eid:tag envUri
	qparts := strings.Split(strings.TrimPrefix(envUri, "oci://"), "?")
	parts := strings.Split(qparts[0], "/")
	img := parts[len(parts)-1]
	iparts := strings.Split(img, ":")

	// fmt.Println("le.persist.vars", qparts, parts, img, iparts)

	// what about from?
	if tEnv == nil {
		tEnv = &tableEnviron{}
	}
	tEnv.Eid = iparts[0]
	tEnv.Tag = iparts[1]
	tEnv.Uri = envUri

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

func replaceTag(envUri, nextTag string) string {
	// preserve any query params
	qparts := strings.Split(envUri, "?")

	// replace tag
	parts := strings.Split(qparts[0], ":")
	parts[len(parts)-1] = nextTag

	// preserve any query params
	if len(qparts) > 1 {
		parts = append(parts, qparts[1:]...)
	}

	// return a new Uri
	return strings.Join(parts, "")
}

func extractPath(envUri string) (string, error) {
	qparts := strings.Split(envUri, "?")
	if len(qparts) < 2 {
		return "", fmt.Errorf("missing query params in Uri to extact path in: %q", envUri)
	}
	vals, err := url.ParseQuery(qparts[1])
	if err != nil {
		return "", fmt.Errorf("while parsing query params in: %q, %w", envUri, err)
	}
	path := vals.Get("path")
	if path == "" {
		return "", fmt.Errorf("empty path in: %q", envUri)
	}
	return path, nil
}

func extractPathEmptyOk(envUri string) (string, error) {
	qparts := strings.Split(envUri, "?")
	if len(qparts) < 2 {
		return "", fmt.Errorf("missing query params in Uri to extact path in: %q", envUri)
	}
	vals, err := url.ParseQuery(qparts[1])
	if err != nil {
		return "", fmt.Errorf("while parsing query params in: %q, %w", envUri, err)
	}
	path := vals.Get("path")
	return path, nil
}

func (le *localEnviron) getEnvironEntry(envUri string) (tableEnviron, error) {
	var foundEnv tableEnviron
	e, err := url.Parse(envUri)
	if err != nil {
		return foundEnv, fmt.Errorf("database error while fetching environ: %w", err)
	}
	key := fmt.Sprintf("%s%s", e.Host, e.Path)
	// fmt.Println("le.lookup", e, key)

	err = le.db.WithContext(le.ctx).
		Where(&tableEnviron{
			Uri: key,
		}).
		First(&foundEnv).Error

	if err != nil {
		// For any error including ErrRecordNotFound, return it as a system error.
		return foundEnv, fmt.Errorf("database error while fetching environ: %w", err)
	}

	return foundEnv, nil
}

func (le *localEnviron) listEnvironEntry() ([]tableEnviron, error) {
	var foundEnvs []tableEnviron

	err := le.db.WithContext(le.ctx).
		Find(&foundEnvs).Error

	if err != nil {
		// For any error including ErrRecordNotFound, return it as a system error.
		return foundEnvs, fmt.Errorf("database error while listing environs: %w", err)
	}

	return foundEnvs, nil
}
