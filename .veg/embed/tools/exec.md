`exec` is for executing one or more commands in a bash script.
All arguements are required for this tool.
Prefer calling this tool once with multiple commands or a script
instead of many times with one command.

1. `key` - the cache key for storing output and exit code
2. `script` - the command, sequence of commands, or script to run

The Results will be stored into the key/value cache and made available to you in your next turn.

1. `ExitCode: #`
2. `Stdout: ...`
2. `Stderr: ...`

Under-the-hood, the `script` arg is run inside a <container>, wrapped in a script seen below.
Thus, there is no need to provide the header lines for a script.
The `script` is run from the `basedir` in your <env>.
You have large freedom to run any commands from inside the isolated <container>.

```sh
#!/bin/bash
set -euo pipefail

<any_command_or_script_contents>
```