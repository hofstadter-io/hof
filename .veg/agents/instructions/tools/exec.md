
`exec` is for executing one or more commands in a bash script.

Under-the-hood, the `script` arg is run inside a shell, wrapped in a script seen below.
Thus, there is no need to provide the header lines for a script.

```sh
#!/bin/bash
set -euo pipefail

<command_or_script_contents>

```

Output and exit code will be returned after the tool has completed.
