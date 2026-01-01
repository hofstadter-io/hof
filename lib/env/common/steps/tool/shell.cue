package tool

import (
	"github.com/hofstadter-io/hof/lib/env/common/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

zsh: {
	// hmm, this is tied to app, but we need general 'utils.pkg'
	install: utils.apt.install & {#pkgs: ["zsh"]}

	customize: [
		omz,
		makeDefault,
	]

	omz: env.Exec & {
		args: ["sh", "-c", _script]
		_script: """
			sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
			sed -i 's/^ZSH_THEME=.*/ZSH_THEME="frisk"/' /root/.zshrc
			"""
	}
	makeDefault: [
		env.Args & {args: ["zsh"]},
		env.Term & {args: ["zsh"]},
		env.Entrypoint & {args: ["zsh"]},
	]
}
