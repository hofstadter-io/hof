// +build !windows,!darwin

package configdir_test

import (
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu/configdir"
)

func reset() {
	// Unset environment variables to make the tests deterministic and choose
	// their own default values. The individual test runners may override these
	// variables for their own use case.
	os.Setenv("HOME", "/home/user")
	os.Setenv("XDG_CONFIG_DIRS", "")
	os.Setenv("XDG_CONFIG_HOME", "")
	os.Setenv("XDG_CACHE_HOME", "")
	configdir.Refresh()
}

func TestSystemConfig(t *testing.T) {
	reset()

	// Cases to test.
	var tests = []TestCase{
		TestCase{
			Paths:  []string{},
			Values: []string{"/etc/xdg"},
		},
		TestCase{
			Paths:  []string{"vendor-name"},
			Values: []string{"/etc/xdg/vendor-name"},
		},
		TestCase{
			Paths:  []string{"vendor-name", "app-name"},
			Values: []string{"/etc/xdg/vendor-name/app-name"},
		},

		// With custom XDG paths...
		TestCase{
			Env:     "/etc/xdg:/opt/global/conf",
			Refresh: true,
			Paths:   []string{},
			Values:  []string{"/etc/xdg", "/opt/global/conf"},
		},
		TestCase{
			Paths:  []string{"vendor-name"},
			Values: []string{"/etc/xdg/vendor-name", "/opt/global/conf/vendor-name"},
		},
		TestCase{
			Paths:  []string{"vendor-name", "app-name"},
			Values: []string{"/etc/xdg/vendor-name/app-name", "/opt/global/conf/vendor-name/app-name"},
		},
	}

	for _, test := range tests {
		if test.Refresh {
			os.Setenv("XDG_CONFIG_DIRS", test.Env)
			configdir.Refresh()
		}

		result := configdir.SystemConfig(test.Paths...)

		// Make sure we got the expected result back.
		if len(result) != len(test.Values) {
			t.Errorf("SystemConfig didn't give the expected number of results. "+
				"Expected %d, got %d (env: %s, input paths: %v, result paths: %v)\n",
				len(test.Values),
				len(result),
				test.Env,
				test.Paths,
				result,
			)
			continue
		}

		// Make sure each result is what we expect.
		for i, path := range result {
			if path != test.Values[i] {
				t.Errorf("Got wrong path result on index %d. "+
					"Expected %s, got %s\n",
					i,
					test.Values[i],
					path,
				)
			}
		}
	}
}

func TestLocalConfig(t *testing.T) {
	testLocalCommon(t, "config", "/home/user/.config", "/opt/local")
}

func TestLocalCache(t *testing.T) {
	testLocalCommon(t, "cache", "/home/user/.cache", "/tmp/cache")
}
