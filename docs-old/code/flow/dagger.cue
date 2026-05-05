package flow

import (
	"dagger.io/dagger"
	"dagger.io/dagger/core"
	"universe.dagger.io/docker"
	"universe.dagger.io/git"

	"github.com/hofstadter-io/harmony-cue/testers"
)

dagger.#Plan

client: commands: clean: {_done: _, name: "bash", args: ["-c", "rm -rf examples/ schemas/"]}
client: filesystem: ".": {
	read: contents:  dagger.#FS
	write: contents: actions.all.output
}

actions: {
	versions: testers.Versions
	versions: hof: "v0.6.3"

	code: git.#Pull & {
		remote:     "https://github.com/hofstadter-io/hof"
		ref:        versions.hof
		keepGitDir: true
	}

	image: testers.Build & {"versions": versions}

	copy: docker.#Run & {
		input: image.output
		command: {
			name: "bash"
			args: ["-c", _script]
			_script: #"""
				echo "hello"
				mkdir -p /final/{examples,schemas}

				# copy schemas
				pushd flow 
				for file in `find ./ -type f -name schema.cue`; do
				  echo "$file"
				  DIR=$(dirname $file)
				  echo "$file - $DIR"
				  mkdir -p /final/schemas/$DIR 
				  cp $file /final/schemas/$file
				done
				popd

				# copy testdata txtar files
				cp -r flow/testdata/* /final/examples/

				"""#
		}
		mounts: source: {
			dest:     "/src"
			contents: code.output
		}
		workdir: "/src"
	}

	all: core.#Copy & {
		_dummy:   client.commands.clean._done
		input:    client.filesystem.".".read.contents
		contents: copy.output.rootfs
		source:   "/final/"
	}

}
