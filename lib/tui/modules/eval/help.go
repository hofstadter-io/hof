package eval

const EvalHelpText = `[blue::b]Welcome[-] to [gold]_[ivory]Hofstadter[-::-]

  The Hof TUI gives you access to hof's features.
  The main area is a dashboard where you can add
  items or widgets, organize them as you wish,
  and save/load them and their sources. 
  There are commands and hotkeys for most actions.

  [darkgray](scroll for more details)[-]

[blue::b]Legend:[-::-]
  [lime]A[-] = ALT
  [lime]C[-] = CTRL
  [lime]M[-] = META
  [lime]S[-] = SHIFT
  [lime]click[-] to focus a box

[blue::b]App Controls:[-::-]
  [lime]C-<space>[-]    focus the command box
  [lime]A-<space>[-]    focus the main content
  [lime]A-/[-]          show / hide the console
  [lime]A-x[-]          clear the console logs
  [red]Ctrl-Alt-c[-]   close [gold]_[ivory]Hofstadter[-]
  [red]:q[-]           (from the command box)

[blue::b]Panel Management:[-::-]
  [lime]A-J[-]          new item before
  [lime]A-K[-]          new item after
  [lime]A-H[-]          move item before
  [lime]A-L[-]          move item after
  [lime]A-T[-]          split panel into two
  [lime]A-D[-]          delete current item
  [lime]A-F[-]          flip panel flex direction
  [lime]A-P[-]          show borders on panels
  [lime]A-O[-]          show borders on items

[blue::b]Items and Commands:[-::-]
  Panels contain items or other panels for layout.
  Items contain widgets and data sources, and are
  controlled through the command box ([lime]C-<space>[-]).

  Currently, the following items are available, many more to come.

  [gold]Value View[-]   explore a single CUE value as a tree or code
  [gold]Playground[-]   work on CUE or data with scope and final view

  Common commands, details below

  [violet]play[-]         open a fresh [gold]Playground[-]
  [violet]play (args)[-]  open a [gold]Playground[-] with data
  [violet]view (args)[-]  open a [gold]Value View[-] with data
  [violet]help[-]         open these help contents

[blue::b]Items:[-::-]

  [gold]Value View[-] allows you to explore CUE with control of options.
  You will see many options at the top-middle, green is enabled
  The represent various CUE evaluation states that you can play with
  ** They toggle by using the letter without modifiers **

  [lime]TCJY S[-]       Tree, Cue, Json, Yaml, Scope(toggle usage)
  [lime]vcfr[-]         Validate, Concrete, Final, ResolveReferences
  [lime]ei[-]           Ignore Errors, Inline Output
  [lime]doh[-]          Definitions, Optional, Hidden
  [lime]DA[-]           Docs, Attributes
  [lime]A-R[-]          Reload data source and refresh

  [gold]Playground[-] is a multi-widget Item for working with CUE.
  You can edit CUE and see the results in real-time, with optional scope.
  The widgets in the playground are:
    1. the scope (if available)
    2. a value editor
    3. a browswer for the result

  [lime]A-f[-]          Rotate this item  (lowercase of Panel hotkey)
  [lime]A-R[-]          Reload data source and refresh
  [lime]A-S[-]          Toggle scope usage

  [gold]Scope[-] is extra CUE and/or data that your edited CUE value
  will be evaluated with, as the scope or context. This is
  often used to load your CUE and data, and then explore or
  modify it with more CUE. This is especially helpful for
  building CUE that will transform or collect larger values
  into smaller or new values. You have full access to CUE.

[blue::b]Commands:[-::-]

  Items and widgets are controlled through the command box ([lime]C-<space>[-]).
  First, make sure the item you want to change is focused (by [lime]clicking[-] on it).
  The general format for commands is as follows.

  [violet]<item> <modifiers>[-] <data source>

  [violet]<item>[-]       <data>    same as [violet]<item> value[-]
  [violet]<item> value[-] <data>    set the main value
  [violet]<item> scope[-] <data>    set the scope value

  [gold]<data sources>[-]

  <command>  [aqua]eval args and flags[-]
  <command>  [aqua]https://... (any http json response)[-]
  <command>  [aqua]bash <args...>(any bash json output)[-]

  [gold]<other commands>[-]
  [violet]push[-]          playground, push text value to CUE playground
                playground links load like any https://...

[blue::b]Dashboards:[-::-]
  As you layout panels and items, set their data,
  it is likely you will want to save, load, and share these.

  [violet]list[-]          list available dashboards
  [violet]save[-] <name>   save the current view 
  [violet]show[-] <name>   show the save file
  [violet]load[-] <name>   load and replace current

  You can also set panel & item names

  [violet]set.item.name[-] <name>
  [violet]set.panel.name[-] <name>

[blue::b]Getting help:[-::-]

  [violet]help[-]          (open this help content in an item)
  [violet]hof feedback[-]  (run this to start a new GitHub issue)

  Docs:     [dogerblue]https://docs.hofstadter.io[-]
  GitHub:   [dogerblue]https://github.com/hofstadter-io/hof[-]
  Slack:    [dogerblue]https://join.slack.com/t/hofstadter-io/shared_invite/zt-e5f90lmq-u695eJur0zE~AG~njNlT1A[-]
  Discord:  [dogerblue]https://discord.gg/6vgbKvPs[-]
`
