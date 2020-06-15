package runtime

import (
	"testing"
)

func TestEnv(t *testing.T) {
	e := &Env{
		Vars: []string{
			"HOME=/no-home",
			"PATH=/usr/bin",
			"PATH=/usr/bin:/usr/local/bin",
			"INVALID",
		},
	}

	if got, want := e.Getenv("HOME"), "/no-home"; got != want {
		t.Errorf("e.Getenv(\"HOME\") == %q, want %q", got, want)
	}

	e.Setenv("HOME", "/home/user")
	if got, want := e.Getenv("HOME"), "/home/user"; got != want {
		t.Errorf(`e.Getenv("HOME") == %q, want %q`, got, want)
	}

	if got, want := e.Getenv("PATH"), "/usr/bin:/usr/local/bin"; got != want {
		t.Errorf(`e.Getenv("PATH") == %q, want %q`, got, want)
	}

	if got, want := e.Getenv("INVALID"), ""; got != want {
		t.Errorf(`e.Getenv("INVALID") == %q, want %q`, got, want)
	}

	for _, key := range []string{
		"",
		"=",
		"key=invalid",
	} {
		var panicValue interface{}
		func() {
			defer func() {
				panicValue = recover()
			}()
			e.Setenv(key, "")
		}()
		if panicValue == nil {
			t.Errorf("e.Setenv(%q) did not panic, want panic", key)
		}
	}
}
