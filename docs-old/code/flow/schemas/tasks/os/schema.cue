package os

Exec: {
	@task(os.Exec)
	$task: "os.Exec"

	cmd: string | [string, ...string]

	// dir specifies the working directory of the command.
	// The default is the current working directory.
	dir?: string

	// env defines the environment variables to use for this system.
	// If the value is a list, the entries mus be of the form key=value,
	// where the last value takes precendence in the case of multiple
	// occurrances of the same key.
	env: {[string]: string} | [...=~"="]

	// stdout captures the output from stdout if it is of type bytes or string.
	// The default value of null indicates it is redirected to the stdout of the
	// current process.
	stdout?: null | string | bytes

	// stderr is like stdout, but for errors.
	stderr?: null | string | bytes

	// stdin specifies the input for the process. If stdin is null, the stdin
	// of the current process is redirected to this command (the default).
	// If it is of typ bytes or string, that input will be used instead.
	stdin?: *null | string | bytes

	// success is set to true when the process terminates with with a zero exit
	// code or false otherwise. The user can explicitly specify the value
	// force a fatal error if the desired success code is not reached.
	success: bool

	// the exit code of the command
	exitcode: int

	// error from cmd.Run()
	error: string
}

// Get a filelock
FileLock: {
	@task(os.FileLock)
	$task: "os.FileLock"

	// lockfile name
	filename: string

	// read-write (true for read-write, false for read-only)
	rw: bool | *false

	// time.Duration for retries, zero means off
	retry: string | *"0"
}

// release a filelock
FileUnlock: {
	@task(os.FileUnlock)
	$task: "os.FileUnlock"

	// lockfile name
	filename: string
}

// A Value are all possible values allowed in flags.
// A null value unsets an environment variable.
Value: bool | number | *string | null

// Name indicates a valid flag name.
Name: !="" & !~"^[$]"

// Getenv gets and parses the specific command line variables.
Getenv: {
	@task(os.Getenv)
	$task: "os.Getenv"

	// if empty, get all
	{[Name]: Value}
}

Glob: {
	@task(os.Glob)
	$task: "os.Glob"

	// glob patterns to match
	globs: [...string]

	// filepaths found matching any of the globs
	filepaths: [...string]
}

// acts like 'mkdir -p' 
Mkdir: {
	@task(os.Mkdir)
	$task: "os.Mkdir"

	dir: string
}

ReadFile: {
	@task(os.ReadFile)
	$task: "os.ReadFile"

	// filename to read
	filename: string

	// filled by flow
	contents: *string | bytes
}

Sleep: {
	@task(os.Sleep)
	$task: "os.Sleep"

	// time.Duration to sleep for
	duration: string
}

// read from stdin
Stdin: {
	@task(os.Stdout)
	$task: "os.Stdout"

	// optional message to user before reading input
	prompt?: string
	// user input
	contents: string
}

// print to stdout
Stdout: {
	@task(os.Stdout)
	$task: "os.Stdout"

	// text to write
	text: string
}

Watch: {
	@task(fs.Watch)
	$task: "fs.Watch"

	// glob patterns to watch for events 
	globs: [...string]

	// todo, only handles write events
	// should add event to handler before calling

	// a flow handler to run on each event
	handler: {...}

	// (good first issue)
	// debounce?: string // time.Duration
}

WriteFile: {
	@task(os.WriteFile)
	$task: "os.WriteFile"

	filename: string
	contents: string | bytes
	mode:     int | *0o666
}
