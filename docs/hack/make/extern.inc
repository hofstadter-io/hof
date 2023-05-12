# files copied from outside the repo
extern: cmdhelp schemas

schemas:
	# todo, get from github @ version
	cp -r ../schema/* code/hof-schemas

.PHONY: cmdhelp
cmdhelp:
	@hof         -h > code/cmd-help/hof
	@hof mod     -h > code/cmd-help/mod
	@hof flow    -h > code/cmd-help/flow
	@hof dm      -h > code/cmd-help/dm
	@hof gen     -h > code/cmd-help/gen
	@hof fmt     -h > code/cmd-help/fmt
	@hof create  -h > code/cmd-help/create
	@hof chat    -h > code/cmd-help/chat

