package runtime

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/hofstadter-io/hof/script/ast"
)

func (RT *Runtime) expand(s string) string {
	return os.Expand(s, func(key string) string {
		if key1 := strings.TrimSuffix(key, "@R"); len(key1) != len(key) {
			return regexp.QuoteMeta(RT.Getenv(key1))
		}
		return RT.Getenv(key)
	})
}

// Setenv sets the value of the environment variable named by the key.
func (RT *Runtime) Setenv(key, value string) {
	RT.envMap[envvarname(key)] = value
}

// Getenv gets the value of the environment variable named by the key.
func (RT *Runtime) Getenv(key string) string {
	return RT.envMap[envvarname(key)]
}

func (RT *Runtime) Cmd_env(cmd *ast.Cmd, r *ast.Result) (err error) {

	if cmd.Exp != ast.Pass {
		r.Status = 1
		RT.status = r.Status
		return fmt.Errorf("unsupported: !? env")
	}

	r.Status = 0
	RT.status = r.Status

	// no args, rint all env vars
	if len(cmd.Args) == 0 {
		for k, v := range RT.envMap {
			k = envvarname(k)
			fmt.Fprintf(cmd.Result().Stdout, "%s=%s\n", k, v)
		}
		return nil
	}

	// loop over args
	for i := 0; i < len(cmd.Args); i++ {
		env := cmd.Args[i]
		RT.logger.Debug(env)

		if i + 2 < len(cmd.Args) && cmd.Args[i+1] == "=" {
			k, v := env, cmd.Args[i+2]
			i += 2
			if v[0] == '@' {
				fname := v[1:] // for error messages
				if fname == "stdout" {
					v = RT.GetStdout()
				} else if fname == "stderr" {
					v = RT.GetStderr()
				} else {
					data, err := ioutil.ReadFile(RT.MkAbs(fname))
					if err != nil {
						r.Status = 1
						RT.status = r.Status
						return err
					}
					v = string(data)
				}
			}

			// set the env var
			RT.Setenv(k,v)
		} else {
			fmt.Fprintf(cmd.Result().Stdout, "%s=%s\n", env, RT.Getenv(env))
		}

	}

	r.Status = 0
	RT.status = r.Status
	return nil
}

