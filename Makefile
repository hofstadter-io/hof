# First command incase users forget to supply one
# cat self as help for simplicity
help:
	cat Makefile
.PHONY: help

cloc: hof.cue design schema lang lib gen docs test cue.mods go.mod 
	cloc --read-lang-def=$$HOME/jumpfiles/assets/cloc_defs.txt --exclude-dir=cue.mod,vendor $^
