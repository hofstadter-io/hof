# RENAME|REWRITE hof -> veg

also moving to verdverm account


## rewrite

- dynamic gen'd from CUE cli, shouldn't need to recompile, available to users too
    - more minimal CI, veg <user-cmd> should work, probably other env/flow like commands too
- include config in this dynamic system, make it work with multi-tier coming up, this system should be WATCHED
- consolidate core subsystems: veg's, libs, runtime, args/flags
- build multi-tier lookup / unify 
- more contextual awareness, especially as it relates to commands being invoked via vscode/remote/rest/websocket
- offline mode, if have remote settings and not connecting, create local db/reg, sync when connected again later
    - or have this system also aware and multi-db/reg capable
