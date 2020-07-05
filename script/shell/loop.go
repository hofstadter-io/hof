package shell

import (
	"strings"

	"github.com/hofstadter-io/hof/script/readline"

	"github.com/hofstadter-io/hof/script/runtime"
)

func Loop(rt *runtime.Runtime) error {
	cfg := &readline.Config{
		Prompt:                 "> ",
		HistoryFile:            "/tmp/readline-multiline",
	}

	rl, err := readline.NewEx(cfg)

	if err != nil {
		return err
	}
	defer rl.Close()

	var cmds []string
	for {
		line, lerr := rl.Readline()
		if lerr != nil {
			err = lerr
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		cmds = append(cmds, line)
		if strings.HasSuffix(line, "\\") {
			rl.SetPrompt("")
			continue
		}

		cmd := strings.Join(cmds, " ")
		cmds = cmds[:0]
		rl.SetPrompt("> ")
		rl.SaveHistory(cmd)
		println(cmd)
	}

	return err
}

