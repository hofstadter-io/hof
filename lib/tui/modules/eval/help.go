package eval

const EvalHelpText = `
[gold]Hof Tui Controls:       [red](scroll for more)[-]

[blue]Legend:[-]
  [lime]A[-] = ALT      mouse & click
  [lime]C[-] = CTRL     also supported
  [lime]M[-] = META
  [lime]S[-] = SHIFT    <esc> = unfocus

[blue]App:[-]
  [lime]C-<space>[-]    focus the command bar
  [lime]A-<space>[-]    focus the main content
  [lime]A-/[-]          show / hide the console
  [lime]A-z[-]          show / hide the error log

[blue]Exit:[-]          [red]Ctrl-Alt-c[-]    close [gold]_[ivory]Hofstadter[-]

[blue]Panels:[-]
  // navigating between panels & items
  you can also click on a panel to focus
  (planned) A-[hjkl[] (vim style movement)

  // adding & moving panels
  [lime]A-[ JK ][-]     create item, before/after
  [lime]A-[ HL ][-]     move item, prev/next
  [lime]A-T[-]          split item or panel, making new stuff
  [lime]A-D[-]          delete  the current item

  // display options
  [lime]A-F[-]          flip panel flex direction
  [lime]A-P[-]          show borders on panels
  [lime]A-O[-]          show borders on items

[blue]Items:[-]
  currently, you can open a browser or evaluator
  value scope still needs to be reconnected, but will be better than before

  [lime]to work with values...[-]
  1. select an item and CTRL-<space> to get to the command box
  2. eval ....
    a. [hotpink]eval[-] [...cue eval args and flags as normal]
    b. [hotpink]cue/play links[-]: https://cuelang.org/play/...
    c. [hotpink]any http(s)://[-] ... that returns json or CUE
    e. [hotpink]bash[-] [...any command the outputs json or CUE]
    

  tbd...
  A-w          change the current widget
  A-e          edit the current settings
  A-c          connect the current item
  A-?          show / hide help for item 


[blue]The Value Browser:[-]
  The browser is the typical box with colored code or a tree
  You will see many options at the top-middle, green is enabled
  The represent various CUE evaluation states that you can play with
  ** They toggle by using the letter without modifiers **

  [lime]TCJY S[-]       Tree, Cue, Json, Yaml, Scope
  [lime]vcfr[-]         Validate, Concrete, Final, ResolveReferences
  [lime]ei[-]           Ignore Errors, Inline Output
  [lime]doh[-]          Definitions, Optional, Hidden
  [lime]DA[-]           Docs, Attributes

[blue]The Value Playground:[-]
  This is a two-panel Item
  1. the value 2. a browswer for the result
  (3) it has a scope, which will be hooked up again very soon
  You can edit the code and the result will live-reload
  Use [lime]A-f[-] to rotate this Item  (lowercase of the Panel)


[blue]Other:[-]
  There is a [lime]ls[-] command that opens a filebrowser
    (hof tui ls  or   C-<space> ... ls<enter>)
  double clicking a file CUE understands will open
  it back in the evaluator, which is preserved

`

