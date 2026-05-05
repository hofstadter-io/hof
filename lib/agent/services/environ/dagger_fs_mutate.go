package environ

import (
	"fmt"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

func (le *localEnviron) WriteFile(envUri, path, content string) (nextUri string, err error) {
	_, env, err := le.LookupEnviron(envUri)
	if err != nil {
		return "", fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	_, fname := filepath.Split(path)
	f := le.dag.File(fname, content)
	nextLayer, err := env.WithFile(path, f).Sync(le.ctx)
	if err != nil {
		return "", fmt.Errorf("while writing file(%s %s): %w", envUri, path, err)
	}

	nextUri, _, err = IncrementTag(envUri)
	if err != nil {
		return "", fmt.Errorf("while incrementing tag(%s %s): %w", envUri, path, err)
	}

	err = le.persistEnviron(nextUri, nil, nextLayer)
	if err != nil {
		return "", fmt.Errorf("while persisting environ(%s %s): %w", envUri, path, err)
	}

	return nextUri, nil
}

func (le *localEnviron) CreateDirectory(envUri, path string) (nextUri string, err error) {
	_, env, err := le.LookupEnviron(envUri)
	if err != nil {
		return "", fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	nextLayer, err := env.WithDirectory(path, le.dag.Directory()).Sync(le.ctx)
	if err != nil {
		return "", fmt.Errorf("while creating directory(%s %s): %w", envUri, path, err)
	}

	nextUri, _, err = IncrementTag(envUri)
	if err != nil {
		return "", fmt.Errorf("while incrementing tag(%s %s): %w", envUri, path, err)
	}

	err = le.persistEnviron(nextUri, nil, nextLayer)
	if err != nil {
		return "", fmt.Errorf("while persisting environ(%s %s): %w", envUri, path, err)
	}

	return nextUri, nil
}

type EditOp struct {
	Old   string `json:"old"`
	New   string `json:"new"`
	Count int    `json:"count"`
}

func (le *localEnviron) EditFile(envUri, path string, edits []EditOp) (nextUri, nextContent string, err error) {
	_, env, err := le.LookupEnviron(envUri)
	if err != nil {
		return "", "", fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	content, err := env.File(path).Contents(le.ctx)
	if err != nil {
		return "", "", fmt.Errorf("while getting current contents for exit(%s -> %s): %w", envUri, path, err)
	}

	//
	// Check and Replace content
	//
	for _, edit := range edits {
		count := edit.Count
		if count == 0 {
			count = 1
		}
		found := strings.Count(content, edit.Old)
		if found != count {
			fmt.Println("fsEdit.count.error", err)
			err = fmt.Errorf("while editing %q, expected %d matches, but found %d", path, count, found)
			// HMM, should we build up errors so we can be resilient with multiple ops?
			// perhaps we do this at the tool layer? (across files rather than within)
			// we are already segmented here by path, so failing within a single file makes the most sense
			return "", "", err
		}
		content = strings.Replace(content, edit.Old, edit.New, count)
	}

	// Update in Dagger
	_, fname := filepath.Split(path)
	f := le.dag.File(fname, content)
	nextLayer, err := env.WithFile(path, f).Sync(le.ctx)
	if err != nil {
		return "", "", fmt.Errorf("while writing file(%s %s): %w", envUri, path, err)
	}

	// update internal tag
	nextUri, _, err = IncrementTag(envUri)
	if err != nil {
		return "", "", fmt.Errorf("while incrementing tag(%s %s): %w", envUri, path, err)
	}

	// persist to DB & OCI
	err = le.persistEnviron(nextUri, nil, nextLayer)
	if err != nil {
		return "", "", fmt.Errorf("while persisting environ(%s %s): %w", envUri, path, err)
	}

	return nextUri, content, nil
}

func (le *localEnviron) Delete(envUri, path string, recursive bool) (nextUri string, err error) {
	_, env, err := le.LookupEnviron(envUri)
	if err != nil {
		return "", fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	// is it a file or dir?
	s, err := le.Stat(envUri, path, false)
	if err != nil {
		return "", fmt.Errorf("while stat'n path(%s): %w", envUri, err)
	}

	var nextLayer *dagger.Container
	if s.Dir {
		nextLayer = env.WithoutDirectory(path)
	} else {
		nextLayer = env.WithoutFile(path)
	}

	nextLayer, err = nextLayer.Sync(le.ctx)
	if err != nil {
		return "", fmt.Errorf("while deleting path(%s|%s): %w", envUri, path, err)
	}

	nextUri, _, err = IncrementTag(envUri)
	if err != nil {
		return "", fmt.Errorf("while incrementing tag(%s %s): %w", envUri, path, err)
	}

	err = le.persistEnviron(nextUri, nil, nextLayer)
	if err != nil {
		return "", fmt.Errorf("while persisting environ(%s %s): %w", envUri, path, err)
	}

	return nextUri, nil
}

func (le *localEnviron) Copy(envUri, source, destination string, overwrite bool) (nextId string, err error) {
	_, env, err := le.LookupEnviron(envUri)
	if err != nil {
		return "", fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	// is it a file or dir?
	s, err := le.Stat(envUri, source, false)
	if err != nil {
		return "", fmt.Errorf("while finding source(%s): %w", envUri, err)
	}

	var nextLayer *dagger.Container
	if s.Dir {
		d := env.Directory(source)
		nextLayer = env.WithDirectory(destination, d)
	} else {
		f := env.File(source)
		nextLayer = env.WithFile(destination, f)
	}

	nextLayer, err = nextLayer.Sync(le.ctx)
	if err != nil {
		return "", fmt.Errorf("while copying path(%s)[%s,%s]: %w", envUri, source, destination, err)
	}

	nextUri, _, err := IncrementTag(envUri)
	if err != nil {
		return "", fmt.Errorf("while incrementing tag(%s)[%s,%s]: %w", envUri, source, destination, err)
	}

	err = le.persistEnviron(nextUri, nil, nextLayer)
	if err != nil {
		return "", fmt.Errorf("while persisting environ(%s)[%s,%s]: %w", envUri, source, destination, err)
	}

	return nextUri, nil
}

func (le *localEnviron) Move(envUri, source, destination string, overwrite bool) (newId string, err error) {
	_, env, err := le.LookupEnviron(envUri)
	if err != nil {
		return "", fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	// is it a file or dir?
	s, err := le.Stat(envUri, source, false)
	if err != nil {
		return "", fmt.Errorf("while finding source(%s): %w", envUri, err)
	}

	var nextLayer *dagger.Container
	if s.Dir {
		d := env.Directory(source)
		nextLayer = env.WithoutDirectory(source).WithDirectory(destination, d)
	} else {
		f := env.File(source)
		nextLayer = env.WithoutFile(source).WithFile(destination, f)
	}

	nextLayer, err = nextLayer.Sync(le.ctx)
	if err != nil {
		return "", fmt.Errorf("while moving path(%s)[%s,%s]: %w", envUri, source, destination, err)
	}

	nextUri, _, err := IncrementTag(envUri)
	if err != nil {
		return "", fmt.Errorf("while incrementing tag(%s)[%s,%s]: %w", envUri, source, destination, err)
	}

	err = le.persistEnviron(nextUri, nil, nextLayer)
	if err != nil {
		return "", fmt.Errorf("while persisting environ(%s)[%s,%s]: %w", envUri, source, destination, err)
	}

	return nextUri, nil
}
