package runtime

// scriptCmds are the script command implementations.
// Keep list and the implementations below sorted by name.
//
// NOTE: If you make changes here, update doc.go.
//
var scriptCmds = map[string]func(*Script, int, []string){
	"call":    (*Script).CmdCall,
	"cd":      (*Script).CmdCd,
	"chmod":   (*Script).CmdChmod,
	"cmp":     (*Script).CmdCmp,
	"cmpenv":  (*Script).CmdCmpenv,
	"cp":      (*Script).CmdCp,
	"env":     (*Script).CmdEnv,
	"exec":    (*Script).CmdExec,
	"exists":  (*Script).CmdExists,
	"grep":    (*Script).CmdGrep,
	"http":    (*Script).CmdHttp,
	"mkdir":   (*Script).CmdMkdir,
	"regexp":  (*Script).CmdRegexp,
	"rm":      (*Script).CmdRm,
	"unquote": (*Script).CmdUnquote,
	"sed":     (*Script).CmdSed,
	"skip":    (*Script).CmdSkip,
	"stdin":   (*Script).CmdStdin,
	"stderr":  (*Script).CmdStderr,
	"stdout":  (*Script).CmdStdout,
	"status":  (*Script).CmdStatus,
	"stop":    (*Script).CmdStop,
	"symlink": (*Script).CmdSymlink,
	"wait":    (*Script).CmdWait,
}


