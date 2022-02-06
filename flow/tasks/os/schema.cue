package os

Exec: {
  @task(os.Exec)

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
}

ReadFile: {
  @task(os.ReadFile)
  filename: string
  f: filename

  // filled by hof
  contents: string | bytes
}

WriteFile: {
  @task(os.WriteFile)
  filename: string
  f: filename

  // filled by hof
  contents: string | bytes
  mode: int | *0o666
}

Stdin: {
  @task(os.Stdout)
  msg?: string
  contents: string
}

Stdout: {
  @task(os.Stdout)
  text: string
}

// A Value are all possible values allowed in flags.
// A null value unsets an environment variable.
Value: bool | number | *string | null

// Name indicates a valid flag name.
Name: !="" & !~"^[$]"

// Getenv gets and parses the specific command line variables.
Getenv: {
	@task(os.Getenv)

  // if empty, get all
	{[Name]: Value}
}
