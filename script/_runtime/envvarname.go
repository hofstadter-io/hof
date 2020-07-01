// +build !windows

package runtime

func envvarname(k string) string {
	return k
}
