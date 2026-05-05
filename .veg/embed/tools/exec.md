`exec` is for executing one or more commands in a bash script.
Prefer calling this tool once with multiple commands or a script
instead of many times with one command as necessary.

`script` - the command, sequence of commands, or script to run

Under-the-hood, the `script` arg is run inside a <container>, wrapped in a script seen below.
Thus, there is no need to provide the header lines for a script.
The `script` is run from the `basedir` in your <env>.
You have large freedom to run any commands from inside the isolated <container>.

```sh
#!/bin/sh
set -euo pipefail

<any_command_or_script_contents>
```