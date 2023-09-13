package eval

const EvalHelpText = `
[dodgerblue::b]Welcome to[-] [gold::bi]_[ivory]Hofstadter[-::-]

  The Hof TUI gives you access to hof's features in a
  space where you can dynamically explore and work with CUE.
  Load up data and CUE for inspection, sources can come from
  files, commands, or http requests. Use the playground to
  explore, filter, or transform them in real-time.
  Connect them to make chains like a Jupyter notebook.
  Create dashboards of these views, playgrounds, and chains
  and then save, load, or share them with others.

  Currently there is support for [violet]hof eval[-].
  New releases will expand feature support,
  notably [violet]flow[-], [violet]gen[-], and [violet]vet[-] are upcoming.

  [dodgerblue::bu]Contents:[darkgray::-]
    Legend
    App Controls
    Panel Management
    Navigation
    Components 
      Items
      Commands
      Dashbords
    Examples
    Getting Help
[-]
[dodgerblue::bu]Legend:[-::-]
  [lime]A[-] = ALT
  [lime]C[-] = CTRL
  [lime]M[-] = META
  [lime]S[-] = SHIFT

[dodgerblue::bu]App Controls:[-::-]
  [lime]C-<space>[-]    focus the command box
   [lime]C-P[-]           [darkgray](vs code like)[-]
   [lime]<esc>[-]         [darkgray](vim like)[-]
  [lime]A-<space>[-]    focus the main content
  [lime]A-/[-]          show / hide the console
  [lime]A-x[-]          clear the console logs
  [lime]A-s[-]          save console logs to file
  [lime]A-c[-]          copy console logs to clipboard
  [red]Ctrl-Alt-c[-]   close [gold::bi]_[ivory]Hofstadter[-::-] :[
   [red]:q[-]          (from the command box)

[dodgerblue::bu]Panel Management:[-::-]
  [lime]A-J[-]          create item before
  [lime]A-K[-]          create item after
  [lime]A-H[-]          move item before
  [lime]A-L[-]          move item after
  [lime]A-T[-]          split current item
  [lime]A-D[-]          delete current item
  [lime]A-F[-]          toggle flex direction
  [lime]A-P[-]          toggle panel borders
  [lime]A-O[-]          toggle item borders

[dodgerblue::bu]Navigation:[-::-]
  [lime]click[-]        to focus any box   [darkgray](with mouse)[-]
  [lime]C-<left>[-]     focus item left    [darkgray](with arrows)[-]
  [lime]C-<down>[-]     focus item down
  [lime]C-<up>[-]       focus item up
  [lime]C-<right>[-]    focus item right
  [lime]A-h[-]          focus item left    [darkgray](vim style)[-]
  [lime]A-j[-]          focus item down
  [lime]A-k[-]          focus item up
  [lime]A-l[-]          focus item right

[dodgerblue::bu]Components:[-::-]

  The main component is a Dashboard with Panels.
  Panels contain items or other panels for layout.
  Items contain widgets with data sources and are
  controlled through the command box ([lime]C-<space>[-]).

  [darkgray]The [dodgerblue]Command Box[darkgray] will modify the [gold]Item[darkgray] currently or last focused.[-]
  [darkgray]There are modifiers to specifiy non-focused items, see below.[-]

  Currently, the following items are available, many more to come.

  [gold]Viewer    [-]   explore a single CUE value as a tree or code
  [gold]Playground[-]   work on CUE or data with scope and final view

  Common commands, details below

  [violet]play[-]         open a fresh [gold]Playground[-]
  [violet]play (args)[-]  open a [gold]Playground[-] with data
  [violet]view (args)[-]  open a [gold]Viewer[-] with data
  [violet]conn (args)[-]  connect values to create a chain
  [violet]help[-]         open these help contents

[dodgerblue::bu]Items:[-::-]

  [gold]Value View[-] allows you to explore CUE with control of options.
  You will see many options at the top-middle, green is enabled
  The represent various CUE evaluation states that you can play with
  ** They toggle by using the letter without modifiers **

  [lime]TCJY[-]         Tree, Cue, Json, Yaml
  [lime]vcfr[-]         Validate, Concrete, Final, ResolveReferences
  [lime]ei[-]           Ignore Errors, Inline Output
  [lime]doh[-]          Definitions, Optional, Hidden
  [lime]DA[-]           Docs, Attributes

  [gold]Playground[-] is a multi-widget Item for working with CUE.
  You can edit CUE and see the results in real-time, with optional scope.
  The widgets in the playground are:
    [darkgray]1. [gold]viewer[-] for the [lightseagreen]scope[-] (if available)
    [darkgray]2. [gold]editor[-] for the [lightseagreen]main[-] value
    [darkgray]3. [gold]viewer[-] for the [lightseagreen]final[-] value

  [lime]A-f[-]          [lightseagreen]Rotate[-] this item  (lowercase of Panel rotate hotkey)
  [lime]A-R[-]          [lightseagreen]Reload[-] data source and refresh
  [lime]A-s[-]          [lightseagreen]Toggle[-] toggle the scope viewer
  [lime]A-S[-]          [lightseagreen]Toggle[-] toggle the scope usage

  [gold]Scope[-] is extra CUE and/or data that your edited CUE value
  will be evaluated with, as the scope or context. This is used to
  load CUE and data, and then explore or modify it with more CUE.
  The scope can come from any of the sources, including other widgets.
  The concept is tied to the CUE compiler and sets the lexical context for
  the CUE currently being processed. See https://pkg.go.dev/cuelang.org/go@v0.6.0/cue#Scope

[dodgerblue::bu]Commands:[-::-]

  Items and widgets are controlled through the [dodgerblue]Command Box[-] ([lime]C-<space>[-]).
  First, make sure the item you want to change is focused (by [lime]clicking[-] on it).
  The general format for commands is as follows.

  [violet]<command> <args> [lightseagreen]<data source>[-]

  [gold]Data Sources[-]
  [violet]<command> [lightseagreen]<eval args and flags>  [-](same as cue and hof)
  [violet]<command> [lightseagreen]https://...            [-](any http json response)
  [violet]<command> [lightseagreen]bash <args...>         [-](any bash json output)
  [violet]<command> [lightseagreen]conn <path> <expr>     [-](dot path to item, optional expr)

  [gold]Item Commands[-]         [darkgray](these create new widgets)[-]
  [violet]<item>[-]       [lightseagreen]<data>[-]   same as [violet]<item> value[-]
  [violet]<item> value[-] [lightseagreen]<data>[-]   set the main value
  [violet]<item> scope[-] [lightseagreen]<data>[-]   set the scope value
  [gold]items: [violet]play,view,help[-]

  [gold]Update Commands[-]       [darkgray](these modify the widget)[-]
  [violet]set.scope [lightseagreen]<data>[-]      set the scope for the current widget
  [violet]set.value [lightseagreen]<data>[-]      set the value for the current widget

  [gold]Other Commands[-]
  [violet]push[-]                   send editor text to cuelang.org/play
  [violet]write  <file>[-]          save editor text to file
  [violet]export <file>[-]          save final value to file
  [violet]set.panel.name <name>[-]  set the focused panel name
  [violet]set.item.name  <name>[-]  set the focused item name
  [violet]insert <loc>[-]           insert a new item at <loc>
                         (head,tail,prev,next,index)

[dodgerblue::bu]Dashboards:[-::-]
  As you layout panels and items, set their data,
  you will want to save, load, and share these.

  [violet]list[-]          list available dashboards
  [violet]save [lightseagreen]<dest>[-]   save the current view 
  [violet]show [lightseagreen]<dest>[-]   show the save file
  [violet]load [lightseagreen]<dest>[-]   load and replace current

  [gold]Destinations[-]  ([violet]<command> [lightseagreen]<dest>[-])
  [violet]<command> [lightseagreen]./<path>               [-](exact location)
  [violet]<command> [lightseagreen]<name>                 [-](module location)
  [violet]<command> [lightseagreen]@<name>                [-](global, user cache dir)
  [violet]<command> [lightseagreen]https://...            [-](http get with ?id=...)
           (contents are CUE)      (post returns the id)

[dodgerblue::bu]Examples:[-::-]

  [violet]view [lightseagreen]test.cue							        [-]create empty playG with scope from file
  [violet]view [lightseagreen]./                            [-]create empty playG with scope from dir
  [violet]view [lightseagreen]<eval args>                   [-]create empty playG with scope from eval args & flags
  [violet]play [lightseagreen]                              [-]create a new, empty playground
  [violet]play [lightseagreen]test.cue                      [-]create playG with content from a file
  [violet]play scope [lightseagreen]test.cue                [-]create empty playG with scope from file
  [violet]play scope [lightseagreen]./                      [-]create empty playG with scope from dir
  [violet]play scope [lightseagreen]<eval args>             [-]create empty playG with scope from eval args & flags

[dodgerblue::bu]Getting Help:[-::-]

  [violet]help[-]          (open this help content in an item)
  [violet]hof feedback[-]  (run this to start a new GitHub issue)

  [gold]Docs:      [deepskyblue]https://docs.hofstadter.io[-]
  [gold]GitHub:    [deepskyblue]https://github.com/hofstadter-io/hof[-]
  [gold]Slack:     [deepskyblue]https://join.slack.com/t/hofstadter-io/shared_invite/zt-e5f90lmq-u695eJur0zE~AG~njNlT1A[-]
  [gold]Discord:   [deepskyblue]https://discord.gg/6vgbKvPs[-]

`
