package environ

import (
	"fmt"
	"net/url"
	"strings"

	"dagger.io/dagger"
	"github.com/google/uuid"
)

// Create from dir, FROM (containers and EIDs)
// ... but how do we do both
func (le *localEnviron) Create(srcUri, fromUri string) (envUri string, err error) {
	uri, err := url.Parse(fromUri)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("url:", uri.Scheme, uri.Host, uri.Path, uri.Query())
	}

	fmt.Println("fs.Create:", srcUri, fromUri, uri)

	c := le.dag.Container()

	switch uri.Scheme {
	case "git":
		// ...
	case "https":
		r := le.dag.Git(fmt.Sprintf("https://%s%s", uri.Host, uri.Path))
		var d *dagger.Directory
		if uri.Fragment == "" {
			d = r.Head().Tree()
		} else {
			d = r.Ref(uri.Fragment).Tree()
		}
		c = c.WithDirectory(uri.Path, d).WithWorkdir(uri.Path)
	case "file":
		d := le.dag.Host().Directory(uri.Path)
		c = c.WithDirectory(uri.Path, d).WithWorkdir(uri.Path)

	case "veg":
		_, env, err := le.lookupEnviron(fromUri)
		if err != nil {
			return "", fmt.Errorf("while looking up environment(%s): %w", envUri, err)
		}

		c = env
		if err != nil {
			return "", fmt.Errorf("while fetching source image: %w", err)
		}
		fmt.Println("started from:", fromUri)

	case "oci":
		// need to strip oci:// and any query params
		// then use the query params
		from := fmt.Sprintf("%s%s", uri.Host, uri.Path)
		from = strings.TrimPrefix(from, "docker.io/")
		fmt.Println("from:", from)
		c, err = c.From(from).Sync(le.ctx)
		if err != nil {
			return "", fmt.Errorf("while fetching source image: %w", err)
		}
		fmt.Println("pulled:", from)
		// default:
		// 	if strings.HasPrefix(uri.Path, "/") {
		// 	} else {
		// 		d := le.dag.Host().Directory(uri.Path)
		// 		c = c.WithDirectory(uri.Path, d).WithWorkdir(uri.Path)
		// 	}
	}

	name := srcUri
	if name == "" {
		name = fromUri
	}

	// create a new uri
	envUri = fmt.Sprintf("host.docker.internal:5000/%s:%s", uuid.New().String(), "genesis")
	tEnv := &tableEnviron{
		Name: name,
		From: fromUri,
		Src:  srcUri,
	}

	// persist
	fmt.Println("saving as:", envUri)
	err = le.persistEnviron(envUri, tEnv, c)
	fmt.Println("saved:", envUri, err)

	return envUri, err
}

func (le *localEnviron) WriteFile(envUri string, nextTag, content string) (err error) {
	_, env, err := le.lookupEnviron(envUri)
	if err != nil {
		return fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	path, err := extractPath(envUri)
	if err != nil {
		return fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	f := le.dag.File(path, content)
	c, err := env.WithFile(path, f).Sync(le.ctx)
	if err != nil {
		return fmt.Errorf("while writing file(%s %s): %w", envUri, path, err)
	}

	// persist the change
	nextUri := replaceTag(envUri, nextTag)
	return le.persistEnviron(nextUri, nil, c)
}

func (le *localEnviron) EditFile(eid, path string, edits []any) (err error) {

	return nil
}

func (le *localEnviron) Delete(envUri, nextTag string, recursive bool) (err error) {
	_, env, err := le.lookupEnviron(envUri)
	if err != nil {
		return fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	path, err := extractPath(envUri)
	if err != nil {
		return fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	c, err := env.WithoutFile(path).Sync(le.ctx)
	if err != nil {
		return fmt.Errorf("while writing file(%s): %w", envUri, err)
	}

	// persist
	nextUri := replaceTag(envUri, nextTag)
	return le.persistEnviron(nextUri, nil, c)
}

func (le *localEnviron) Copy(envUri, nextTag, source, destination string, overwrite bool) error {
	_, env, err := le.lookupEnviron(envUri)
	if err != nil {
		return fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	// is it a file or dir?
	s, err := le.Stat(envUri, source)
	if err != nil {
		return fmt.Errorf("while finding source(%s): %w", envUri, err)
	}

	var c *dagger.Container
	if s.Dir {
		d := env.Directory(source)
		c, err = env.WithDirectory(destination, d).Sync(le.ctx)
		if err != nil {
			return fmt.Errorf("while copying directory(%s)[%s,%s]: %w", envUri, source, destination, err)
		}
	} else {
		f := env.File(source)
		c, err = env.WithFile(destination, f).Sync(le.ctx)
		if err != nil {
			return fmt.Errorf("while copying file(%s)[%s,%s]: %w", envUri, source, destination, err)
		}
	}

	// persist
	nextUri := replaceTag(envUri, nextTag)
	return le.persistEnviron(nextUri, nil, c)
}

func (le *localEnviron) Move(envUri, nextTag, source, destination string) error {
	_, env, err := le.lookupEnviron(envUri)
	if err != nil {
		return fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	// is it a file or dir?
	s, err := le.Stat(envUri, source)
	if err != nil {
		return fmt.Errorf("while finding source(%s): %w", envUri, err)
	}

	var c *dagger.Container
	if s.Dir {
		d := env.Directory(source)
		c, err = env.WithoutDirectory(source).WithDirectory(destination, d).Sync(le.ctx)
		if err != nil {
			return fmt.Errorf("while moving directory(%s)[%s,%s]: %w", envUri, source, destination, err)
		}
	} else {
		f := env.File(source)
		c, err = env.WithoutFile(source).WithFile(destination, f).Sync(le.ctx)
		if err != nil {
			return fmt.Errorf("while moving file(%s)[%s,%s]: %w", envUri, source, destination, err)
		}
	}

	// persist
	nextUri := replaceTag(envUri, nextTag)
	return le.persistEnviron(nextUri, nil, c)
}
