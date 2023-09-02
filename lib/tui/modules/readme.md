# Modules

these are the hof tui modules

- all are registered in registry.go during init
- lib/connector is used to collect instances of an interface each module cares about later (like command box)




### start a new module

from this dir

```
hof gen - -I name=<name> -T _t/+* -O <name>
```
