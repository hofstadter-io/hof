# watch test

```sh
hof gen data.cue -T 'template.txt;out.txt' --watch '*.cue' --watch template.txt [--diff3]
```

1. Run the above
1. open all the files
1. edit the data or template
1. watch the out.txt update
1. `--diff3` is optional, when set you can edit the out.txt too

