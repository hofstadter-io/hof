package environ

import (
	"fmt"
	"net/url"
	"strings"

	"dagger.io/dagger"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

// Create from dir, FROM (containers and EIDs)
// ... but how do we do both

type EnvironCreateOptions struct {
	Name string `json:"name"`

	// where to get a filesystem
	SrcUri  string `json:"srcUri"`
	SrcPath string `json:"srcPath"`

	// where to put the filesystem
	FromUri string `json:"fromUri"` // if not set, FROM SCRATCH
	DstPath string `json:"dstPath"` // where src is attached in FromUri
	Workdir string `json:"workdir"` // the workdir, defaults to DstPath or FromUri's value if not set

	// Env map[string]string
}

// let's make this clearer
// 1. srcUri is files (typically from the outside via a ref, extract just the files from a oci://...)
// 2. fromUri is containers (outside and everything inside)
func (le *localEnviron) Create(opts EnvironCreateOptions) (envUri string, err error) {
	fmt.Printf("fs.Create.input: %#+v\n", pretty.Formatter(opts))

	// setup container for consistency, `FROM scratch` if not set
	c := le.dag.Container()

	if opts.FromUri != "" {
		furi, err := url.Parse(opts.FromUri)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("  furi:", furi.Scheme, furi.Host, furi.Path, furi.Query())
		}

		switch furi.Scheme {
		case "veg", "oci":
			// need to strip oci:// and any query params
			// then use the query params
			from := fmt.Sprintf("%s%s", furi.Host, furi.Path)

			// dagger/buildkit doesn't like docker.io in the FROM value?
			from = strings.TrimPrefix(from, "docker.io/")
			fmt.Println("from:", from)

			// FROM that image
			c, err = c.From(from).Sync(le.ctx)
			if err != nil {
				return "", fmt.Errorf("while fetching source image: %w", err)
			}

		//
		// room for more sandboxes and generalization
		//

		default:
			if err != nil {
				return "", fmt.Errorf("unsupported from scheme: %q", furi.Scheme)
			}

		}

		// set the workdir if provided
		if opts.Workdir != "" {
			c = c.WithWorkdir(opts.Workdir)
		}
	}

	// do we have a source filesystem to attach?
	if opts.SrcUri != "" {
		// get source filesystem
		suri, err := url.Parse(opts.SrcUri)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("  suri:", suri.Scheme, suri.Host, suri.Path, suri.Query())
		}

		var d *dagger.Directory
		switch suri.Scheme {
		case "git":
			// ...
		case "https":
			r := le.dag.Git(fmt.Sprintf("https://%s%s", suri.Host, suri.Path))
			if suri.Fragment == "" {
				d = r.Head().Tree()
			} else {
				d = r.Ref(suri.Fragment).Tree()
			}
		case "file":
			d = le.dag.Host().Directory(suri.Path)

		// this should probably be oci only?
		case "veg", "oci":
			// do we need to trim prefixes, we really need to get consistent about them (probobly keep, not trim, when persisting to db or state)
			_, env, err := le.LookupEnviron(opts.SrcUri)
			if err != nil {
				return "", fmt.Errorf("while looking up environment(%s): %w", envUri, err)
			}

			// todo, this should be parameterized
			d = env.Directory(".")
			if err != nil {
				return "", fmt.Errorf("while fetching source image: %w", err)
			}
			fmt.Println("filsys from:", opts.SrcUri)

		default:
			if err != nil {
				return "", fmt.Errorf("unsupported source schem: %q", suri.Scheme)
			}

		}

		// get subpath within source filesystem
		if opts.SrcPath != "" {
			d = d.Directory(opts.SrcPath)
		}

		// attach source
		dp := opts.DstPath
		if dp == "" {
			dp = "."
		}
		c = c.WithDirectory(dp, d)
	}

	name := opts.Name
	if name == "" {
		name = opts.SrcUri
	}
	if name == "" {
		name = opts.FromUri
	}

	// create a new uri
	envUri = fmt.Sprintf("host.docker.internal:5000/%s:%s", uuid.New().String(), "0")
	tEnv := &tableEnviron{
		Name:    name,
		SrcUri:  opts.SrcUri,
		SrcPath: opts.SrcPath,
		FromUri: opts.FromUri,
		DstPath: opts.DstPath,
	}

	// persist
	fmt.Println("saving as:", envUri)
	err = le.persistEnviron(envUri, tEnv, c)
	fmt.Println("saved:", envUri, err)

	return envUri, err
}
