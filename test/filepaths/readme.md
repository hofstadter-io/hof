filepaths are a pain... especitally with shadow


```
app/
	cue.mod/
	.hof/shadow
		a/out
		b/out

	/gen
	/templates

	/a
		hof.cue
		/out
	/b
		/out

```

- `hof gen -o out` from both a & b


When all gens output to same dir

- `hof gen -G a -G b`, `hof gen -G a` and not delete b


Test when there is no cue.mod
