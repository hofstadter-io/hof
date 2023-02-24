package runtime

import (
	"io/ioutil"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/hofstadter-io/hof/lib/gotils/txtar"
)

// setupTest sets up the test execution temporary directory and environment.
// It returns the script content section of the txtar archive.
func (ts *Script) setupTest() string {
	ts.workdir = filepath.Join(ts.testTempDir, "script-"+ts.name)

	ts.resetDirs()

	env := &Env{
		Vars: []string{
			"WORK=" + ts.workdir, // must be first for ts.abbrev
			"PATH=" + os.Getenv("PATH"),
			"USER=" + os.Getenv("USER"),
			homeEnvName() + "=" + ts.workdir + "/home",
			tempEnvName() + "=" + filepath.Join(ts.workdir, "tmp"),
			"devnull=" + os.DevNull,
			"/=" + string(os.PathSeparator),
			":=" + string(os.PathListSeparator),
		},
		WorkDir: ts.workdir,
		Values:  make(map[interface{}]interface{}),
		Cd:      ts.workdir,
		ts:      ts,
	}

	return ts.setupFromEnv(env)
}

func (ts *Script) resetDirs() {
	dirs := []string{
		"tmp",
		"home",
	}
	for _, d := range dirs {
		f := filepath.Join(ts.workdir, d)
		ts.Check(os.RemoveAll(f))
		ts.Check(os.MkdirAll(f, 0777))
	}
}

// setupRun sets up the script execution for working in the current directory.
// the current environment will be exposed to the script
// It returns the script content section of the txtar archive.
func (ts *Script) setupRun() string {

	// expose external ENV here
	env := &Env{
		Vars:    os.Environ(),
		WorkDir: ts.workdir,
		Values:  make(map[interface{}]interface{}),
		Cd:      ts.workdir,
		ts:      ts,
	}

	return ts.setupFromEnv(env)
}

func (ts *Script) setupFromEnv(env *Env) string {
	// Must preserve SYSTEMROOT on Windows: https://github.com/golang/go/issues/25513 et al
	if goruntime.GOOS == "windows" {
		env.Vars = append(env.Vars,
			"SYSTEMROOT="+os.Getenv("SYSTEMROOT"),
			"exe=.exe",
		)
	} else {
		env.Vars = append(env.Vars,
			"exe=",
		)
	}
	ts.cd = env.Cd

	ts.unpackArchive()

	// Run any user-defined setup.
	if ts.params.Setup != nil {
		ts.Check(ts.params.Setup(env))
	}

	ts.cd = env.Cd
	ts.env = env.Vars
	ts.values = env.Values

	ts.setupEnvMap()

	return ts.orig
}

func (ts *Script) unpackArchiveOrig() {
	// Unpack archive.
	a, err := txtar.ParseFile(ts.file)
	ts.Check(err)
	ts.archive = a
	for _, f := range a.Files {
		name := ts.MkAbs(ts.expand(f.Name))
		ts.scriptFiles[name] = f.Name
		ts.Check(os.MkdirAll(filepath.Dir(name), 0777))
		ts.Check(ioutil.WriteFile(name, f.Data, 0666))
	}
}

// Unpack archive.
func (ts *Script) unpackArchive() {
	var err error
	a, err := txtar.ParseFile(ts.file)
	ts.Check(err)
	ts.archive = a
	for _, f := range a.Files {

		// sub name
		subd := ts.expand(f.Name)
		name := subd
		var fmode os.FileMode
		fmode = 0644
		chown := ""

		// check for filemode/chown
		flds := strings.Split(subd, ";")
		if len(flds) > 1 {
			name = flds[0]
			if len(flds) == 3 {
				if strings.Contains(flds[2], ":") {
					chown = flds[2]
				} else {
					fm, err := strconv.ParseUint(flds[2], 8, 32)
					ts.Check(err)
					fmode = os.FileMode(fm)
				}
			}
			if strings.Contains(flds[1], ":") {
				chown = flds[1]
			} else {
				fm, err := strconv.ParseUint(flds[1], 8, 32)
				ts.Check(err)
				fmode = os.FileMode(fm)
			}
		}

		// make the file
		fname := ts.MkAbs(name)
		ts.scriptFiles[fname] = f.Name
		os.RemoveAll(fname)
		ts.Check(os.MkdirAll(filepath.Dir(fname), 0755))
		ts.Check(ioutil.WriteFile(fname, f.Data, fmode))
		if chown != "" {
			flds := strings.Split(chown, ":")
			uid, err := strconv.Atoi(flds[0])
			ts.Check(err)
			gid, err := strconv.Atoi(flds[1])
			ts.Check(err)
			ts.Check(os.Chown(fname, uid, gid))
		}
	}

	ts.orig = string(a.Comment)
}

func (ts *Script) setupEnvMap() {
	ts.envMap = make(map[string]string)
	for _, kv := range ts.env {
		if i := strings.Index(kv, "="); i >= 0 {
			ts.envMap[envvarname(kv[:i])] = kv[i+1:]
		}
	}
}

func (ts *Script) setupZap() {

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// setup our config and console encoder
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core)
	ts.Logger = logger.Sugar()
}
