# First command incase users forget to supply one
# cat self as help for simplicity
help:
	cat Makefile
.PHONY: help

clocall: hof.cue design schema lang lib gen ci cmd docs test cue.mods go.mod 
	cloc --read-lang-def=$$HOME/hof/jumpfiles/assets/cloc_defs.txt --exclude-dir=cue.mod,vendor $^

clocgen: gen ci cmd 
	cloc --read-lang-def=$$HOME/hof/jumpfiles/assets/cloc_defs.txt --exclude-dir=cue.mod,vendor $^

clocdesign: hof.cue design 
	cloc --read-lang-def=$$HOME/hof/jumpfiles/assets/cloc_defs.txt --exclude-dir=cue.mod,vendor $^

clochof: hof.cue design lang lib 
	cloc --read-lang-def=$$HOME/hof/jumpfiles/assets/cloc_defs.txt --exclude-dir=cue.mod,vendor $^

cloccode: cmd lang lib gen cmd
	cloc --read-lang-def=$$HOME/hof/jumpfiles/assets/cloc_defs.txt --exclude-dir=cue.mod,vendor $^

clocdev: hof.cue design schema lang lib docs test cue.mods go.mod 
	cloc --read-lang-def=$$HOME/hof/jumpfiles/assets/cloc_defs.txt --exclude-dir=cue.mod,vendor $^

WORKFLOWS = default \
	test_mod

.PHONY: workflow
workflows = $(addprefix workflow_, $(WORKFLOWS))
workflow: $(workflows)
$(workflows): workflow_%:
	@cue export --out yaml .github/workflows/$(subst workflow_,,$@).cue > .github/workflows/$(subst workflow_,,$@).yml
