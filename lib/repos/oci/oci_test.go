package oci

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPushAndPull(t *testing.T) {
	// const ociRegistryAddr = "localhost:1111"
	const ociRegistryAddr = "us-central1-docker.pkg.dev/hof-io--develop/testing"

	// Start OCI registry
	//	go func() {
	//		r := registry.New()
	//		http.ListenAndServe(ociRegistryAddr, r)
	//	}()

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
		NewDir(HofstadterModuleDeps, "cue.mod", []string{
			"*",
			"!module.cue",
			"!sums.cue",
		}),
		NewDir(HofstadterModuleCode, "", []string{
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
