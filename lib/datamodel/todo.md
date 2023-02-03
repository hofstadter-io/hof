# dm-revamp branch

MORE TESTS!

- [x] history | log commands, do we really need both?
- [x] some way to print nesting for user, info detail level & format
- [x] some way to query for nested diff/hist where appropriate
- [x] cleanup flags
- [x] stepwise diff (between snapshots)

- [x] checkpoint message -> description
- [x] checkpoint only things with diffs

- [x] print log by timestamp, displaying what changed at each

- work on CLI help message, mainly needs examples

- work on hof docs
  - datamodel docs
	(also today)
	- datafiles in code gen docs
	- module docs


# next (friday)

- load into generators with history

# later

- filtering by temporal flags

- tags & appropriate bumping
  - should major/minor bumps (backwards incompatiable changes) force a full checkpoing & versioning?

- be able to print the stepwise diff (via log command?)
  we do this a bit already, but there may be other commands where something like this would be helpful

- support for custom lenses (in diff format for now?), what about references?

- MULTI ~limited~ diff if there is a nested diff
  - local & full diffs, based on @node()? (or @history()?)

- diff with color (--no-color flag too, can we detect TTY)
  we can probably do this with a stack based parser (by {}) on the string
	we only need to recognize "+" & "-" and then find newline or matching curlies
