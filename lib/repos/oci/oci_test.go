package oci

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	workingDir := t.TempDir()

	one := path.Join(workingDir, "one")
	require.NoError(t, os.Mkdir(one, 0o700))
	two := path.Join(workingDir, "two")
	require.NoError(t, os.Mkdir(two, 0o700))
	three := path.Join(workingDir, "three")
	require.NoError(t, os.Mkdir(three, 0o700))

	dirs := []Dir{
		NewDir("one", nil),
		NewDir("two", nil),
		NewDir("three", nil),
	}

	var (
		output  = t.TempDir()
		content = []byte("aaa111")
	)

	for _, d := range dirs {
		f, err := os.Create(path.Join(workingDir, d.Path, "file.txt"))
		require.NoError(t, err)

		_, err = f.Write(content)
		require.NoError(t, err)
		require.NoError(t, f.Close())
	}

	img, err := Build(workingDir, dirs)
	require.NoError(t, err)

	f, err := os.Create(path.Join(output, "img.tar"))
	require.NoError(t, err)
	defer f.Close()

	tag, err := name.NewTag("test:1111")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tarball.Write(tag, img, &buf)
	require.NoError(t, err)

	untar(&buf, output)

	var (
		manifest bool
		digest   bool
	)

	err = filepath.Walk(output, func(_ string, i os.FileInfo, err error) error {
		require.NoError(t, err)

		if i.IsDir() {
			return nil
		}

		switch n := i.Name(); {
		case n == "manifest.json":
			manifest = true
		case strings.HasPrefix(n, "sha256:"):
			digest = true
		}

		return nil
	})
	require.NoError(t, err)

	require.True(t, manifest)
	require.True(t, digest)
}

func TestPushAndPull(t *testing.T) {
	const ociRegistryAddr = "localhost:1111"

	// Start OCI registry
	go func() {
		r := registry.New()
		http.ListenAndServe(ociRegistryAddr, r)
	}()

	var (
		rootOrig = t.TempDir()
		rootNew  = t.TempDir()
	)

	// Simulate user activity
	require.NoError(t, os.WriteFile(path.Join(rootOrig, "one.txt"), []byte("111"), 0o700))
	require.NoError(t, os.WriteFile(path.Join(rootOrig, "two.txt"), []byte("222"), 0o700))

	subDir := path.Join(rootOrig, "subdir")
	require.NoError(t, os.Mkdir(subDir, 0o700))
	require.NoError(t, os.WriteFile(path.Join(subDir, "three.txt"), []byte("333"), 0o700))

	run(t, rootOrig, "git", "init")
	run(t, rootOrig, "git", "add", ".")
	run(t, rootOrig, "git", "commit", "-m", "'First commit'")
	run(t, rootOrig, "go", "mod", "init", "github.com/hofstadter-io/oci-test")
	run(t, rootOrig, "hof", "mod", "init", "github.com/hofstadter-io/oci-test")
	run(t, rootOrig, "hof", "mod", "get", "github.com/hofstadter-io/hof@next")
	run(t, rootOrig, "tree", "-a")

	img, err := Build(rootOrig, []Dir{
		NewDir("cue.mod", []string{
			"*",
			"!module.cue",
			"!sums.cue",
		}),
		NewDir("", []string{
			"cue.mod/pkg",
			".git",
		}),
	})
	require.NoError(t, err)

	layers, err := img.Layers()
	require.NoError(t, err)
	require.Equal(t, 2, len(layers))

	tag := fmt.Sprintf("%s/test:latest", ociRegistryAddr)

	err = Push(tag, img)
	require.NoError(t, err)

	// Pull the image in a new directory
	err = Pull(tag, rootNew)
	require.NoError(t, err)

	run(t, rootNew, "tree", "-a")
}

func run(t *testing.T, dir, name string, args ...string) {
	t.Helper()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	t.Log(string(out))
	require.NoError(t, err)
}
