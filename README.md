# hof - the high code framework

The `hof` tool tries to remove redundent development activities
by using high level designs, code generation, and diff3
while letting you write custom code directly in the output.
( low-code for developers )

- Users write Single Source of Truth (SSoT) design for data models and the application generators
- `hof` reads the SSoT, processes it through the code generators, and outputs directories and files
- Users can write custom code in the output, change their designs, and regenerate code in any order
- `hof` can be customized and extended by only editing text files and not `hof` source code.
- Use your own tools, technologies, and practices, `hof` does not make any choices for you
- `hof` is powered by Cue (https://cuelang.org & https://cuetorials.com)

## Install

You will have to download `hof` the first time.
After that `hof` will prompt you to update and
install new releases as they become available.

```shell
export HOF_VER=0.6.1

# Install (Linux, Mac, Windows)
curl -LO https://github.com/hofstadter-io/hof/releases/download/v${HOF_VER}/hof_${HOF_VER}_$(uname)_$(uname -m)
mv hof_${HOF_VER}_$(uname)_$(uname -m) /usr/local/bin/hof

# Shell Completions (bash, zsh, fish, power-shell)
echo ". <(hof completion bash)" >> $HOME/.profile
source $HOME/.profile

# Show the help text
hof --help
```

You can always find the latest version from the
[releases page](https://github.com/hofstadter-io/hof/releases)
or use `hof` to install a specific version of itself with `hof update --version vX.Y.Z`.


## Documentation

Please see __https://docs.hofstadter.io__ to learn more.

The [first-example](https://docs.hofstadter.io/first-example)
will take you through the process
of creating and using a simple generator

Join us on Slack! [https://hofstadter-io.slack.com](https://join.slack.com/t/hofstadter-io/shared_invite/zt-e5f90lmq-u695eJur0zE~AG~njNlT1A)

