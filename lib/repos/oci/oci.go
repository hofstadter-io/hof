package oci

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	// "path"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/google/go-containerregistry/pkg/v1/stream"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

const (
	HofstadterSchema1Beta types.MediaType = "application/vnd.hofstadter.module.v1beta1+json"
	HofstadterModuleDeps  types.MediaType = "application/vnd.hofstadter.module.deps.tar.gz"
	HofstadterModuleCode  types.MediaType = "application/vnd.hofstadter.module.code.tar.gz"
)

func IsNetworkReachable(mod string) (bool, error) {
	_, err := crane.Manifest(mod, crane.WithAuthFromKeychain(authn.DefaultKeychain))

	var terr *transport.Error
	if errors.As(err, &terr) {
		if len(terr.Errors) != 1 {
			return false, fmt.Errorf("multiple transport errors: %w", terr)
		}

		switch c := terr.Errors[0].Code; c {
		case transport.ManifestUnknownErrorCode:
			return true, nil
		case transport.NameUnknownErrorCode:
			return false, errors.New("remote repo does not exist")
		default:
			return false, fmt.Errorf("unhandled transport code: %s", c)
		}
	}

	return err == nil, err
}

func Pull(tag, path string) error {
	ref, err := name.ParseReference(tag)
	if err != nil {
		return fmt.Errorf("name parse reference: %w", err)
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return fmt.Errorf("remote image: %w", err)
	}

	r := mutate.Extract(img)
	defer r.Close()

	if err := untar(r, path); err != nil {
		return fmt.Errorf("untar: %w", err)
	}

	return nil
}

func untar(r io.Reader, target string) error {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("tar reader next: %w", err)
		}

		var (
			p = filepath.Join(target, header.Name)
			i = header.FileInfo()
		)

		if i.IsDir() {
			if err = os.MkdirAll(p, i.Mode()); err != nil {
				return fmt.Errorf("mkdir all: %w", err)
			}
			continue
		}

		f, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, i.Mode())
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}

		defer f.Close()

		if _, err = io.Copy(f, tr); err != nil {
			return fmt.Errorf("io copy: %w", err)
		}
	}

	return nil
}

func Push(tag string, img v1.Image) error {
	ref, err := name.ParseReference(tag)
	if err != nil {
		return fmt.Errorf("name parse reference: %w", err)
	}

	fmt.Println("pushing...")
	if err = remote.Write(ref, img, remote.WithAuthFromKeychain(authn.DefaultKeychain)); err != nil {
		return fmt.Errorf("remote write: %w", err)
	}

	return nil
}

func Build(workingDir string, dirs []Dir) (v1.Image, error) {
	var layers []v1.Layer

	for _, d := range dirs {
		// todo, enable printing base on verbosity
		fmt.Println("layer:", d.mediaType)
		l, err := layer(workingDir, d)
		if err != nil {
			return nil, fmt.Errorf("layer: %w", err)
		}

		layers = append(layers, l)
	}

	e := mutate.MediaType(empty.Image, types.OCIManifestSchema1)
	e = mutate.ConfigMediaType(e, HofstadterSchema1Beta)

	img, err := mutate.AppendLayers(e, layers...)
	if err != nil {
		return nil, fmt.Errorf("append layers: %w", err)
	}

	return img, nil
}

func layer(wd string, d Dir) (v1.Layer, error) {
	var (
		buf bytes.Buffer
		w   = tar.NewWriter(&buf)
	)

	err := filepath.Walk(d.relPath, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if d.Excluded(p) {
			return nil
		}
		// todo, print included filename based on verbosity
		fmt.Println(" ", p)

		h, err := tar.FileInfoHeader(i, "")
		if err != nil {
			return fmt.Errorf("tar file info header: %w", err)
		}

		h.Name = strings.ReplaceAll(p, wd, "")

		if err = w.WriteHeader(h); err != nil {
			return fmt.Errorf("tar write header: %w", err)
		}

		if i.IsDir() {
			return nil
		}

		f, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("open %s: %w", p, err)
		}

		defer f.Close()

		if _, err = io.Copy(w, f); err != nil {
			return fmt.Errorf("copy %s: %w", p, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("filepath walk: %w", err)
	}

	if err = w.Close(); err != nil {
		return nil, fmt.Errorf("tar writer close: %w", err)
	}

	var (
		rc = io.NopCloser(&buf)
		mt = stream.WithMediaType(d.mediaType)
	)

	return stream.NewLayer(rc, mt), nil
}
