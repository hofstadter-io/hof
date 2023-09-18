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

  Currently there is support for [blue]hof/eval[-] and [blue]hof/flow[-].
  New releases will expand core feature support in this tui,
  notably [blue]datamodel[-], [blue]gen[-], and [blue]vet[-] are upcoming.


[dodgerblue::bu]Contents:[darkgray::-]
  App Controls
  Overview 
  Items
  Commands
  Dashbords
  Getting Help
[-]

[dodgerblue::bu]App Controls:[-::-]
  [lime]Ctrl-<space>[-]   focus the command box
   [lime]Ctrl-P[-]          [darkgray](vs code like)[-]
   [lime]<esc>[-]           [darkgray](vim like)[-]
  [lime]Alt-<space>[-]    focus the main content
  [lime]Alt-/[-]          show / hide the console
  [lime]Alt-x[-]          clear the console logs
  [lime]Alt-s[-]          save console logs to file
  [lime]Alt-c[-]          copy console logs to clipboard
  [red]Ctrl-Alt-c[-]      close [gold::bi]_[ivory]Hofstadter[-::-] :[
   [red]:q[-]            [darkgray](from the command box)[-]


[dodgerblue::bu]Overview:[-::-]

  The main component is a Dashboard with Panels.
  Panels contain items or other panels for layout.
  Items contain widgets with data sources and are
  controlled through the command box ([lime]Ctrl-<space>[-],[lime]Ctrl-P[-],[lime]<esc>[-]).

  [darkgray]The [dodgerblue]Command Box[darkgray] will modify the [gold]Item[darkgray] currently or last focused.[-]
  [darkgray]There are modifiers to specifiy non-focused items, see below.[-]

  Currently, the following items are available, many more to come.

  [gold]Viewer    [-]   explore a single CUE value as a tree or code
  [gold]Playground[-]   work on CUE or data with scope and final view

  Common commands, details below

  [violet]play          [-]open a fresh [gold]Playground[-]
  [violet]play (args)   [-]open a [gold]Playground[-] with data
  [violet]view (args)   [-]open a [gold]Viewer[-] with data
  [violet]conn (args)   [-]connect values to from a chain or DAG
  [violet]watch         [-]watch inputs and refresh contents
  [violet]help          [-]open these help contents


[dodgerblue::bu]Items:[-::-]

  [gold::b]Viewer[-::-] allows you to explore CUE with control of options.
  You will see many options at the top-middle, green is enabled.
  These represent various CUE evaluation states that you can play with.
  ** [::i]They toggle by using the letter without modifiers[-::-] **

    [blue::bu]Hotkeys[-::-]
    [lime]TCJY[-]         Tree, Cue, Json, Yaml
    [lime]vcfr[-]         Validate, Concrete, Final, ResolveReferences
    [lime]ei[-]           Ignore Errors, Inline Output
    [lime]doh[-]          Definitions, Optional, Hidden
    [lime]DA[-]           Docs, Attributes

    [blue::bu]Commands[-::-]
    [violet]view [lightseagreen]<data>            [-]open a new [gold]Viewer[-] with [lightseagreen]<data-source>[-]
    [violet]view [lightseagreen]test.cue          [-]create empty playG with scope from file
    [violet]view [lightseagreen]./                [-]create empty playG with scope from dir
    [violet]view [lightseagreen]<eval args>       [-]create empty playG with scope from eval args & flags

  [gold::b]Playground[-::-] is a multi-widget Item for working with CUE.
  You can edit CUE and see the results in real-time, with optional scope.
  There are also two modes: [blue]hof/eval[-] (default) and [blue]hof/flow[-].

  The widgets in the playground are:
    [darkgray]1. [gold]viewer[-] for the [lightseagreen]scope[-]
    [darkgray]2. [gold]editor[-] for the [lightseagreen]main[-]
    [darkgray]3. [gold]viewer[-] for the [lightseagreen]final[-]

    [blue::bu]Hotkeys[-::-]
    [lime]Alt-f[-]        [lightseagreen]Rotate[-] this item  (lowercase of Panel rotate hotkey)
    [lime]Alt-R[-]        [lightseagreen]Reload[-] data source and refresh
    [lime]Alt-s[-]        [lightseagreen]Toggle[-] toggle the scope viewer
    [lime]Alt-S[-]        [lightseagreen]Toggle[-] toggle the scope usage
    [lime]Alt-E[-]        [lightseagreen]Toggle[-] set to eval mode ([blue]hof/eval[-],default)
    [lime]Alt-W[-]        [lightseagreen]Toggle[-] set to flow mode ([blue]hof/flow[-])

    [blue::bu]Commands[-::-]
    [violet]play [lightseagreen]                              [-]create a new, empty playground
    [violet]play [lightseagreen]test.cue                      [-]create play with content from a file
    [violet]play [lightseagreen]bash kubectl ...              [-]create play with json output of a command
    [violet]play [lightseagreen]https://cuelang.org/play...   [-]create play with content from a cue/play link
    [violet]push [lightseagreen]                              [-]push the current edit text to cue/play, link returned
    [violet]play scope [lightseagreen]test.cue                [-]create new play with scope from file
    [violet]play scope [lightseagreen]./                      [-]create new play with scope from dir
    [violet]play scope [lightseagreen]bash or https://...     [-]create new play with any of the data sources
    [violet]play scope [lightseagreen]<eval args & flags>     [-]create new play with scope from eval args & flags

  There are many more commands for working with widgets, primarily [gold]play[-].
  The [dodgerblue]Commands[-] section below has more examples and advanced options.

  [gold]Scope[-] is extra CUE and/or data that your edited CUE value
  will be evaluated with, as the scope or context. This is used to
  load CUE and data, and then explore or modify it with more CUE.
  The scope can come from any of the sources, including other widgets.
  The concept is tied to the CUE compiler and sets the lexical context for
  the CUE currently being processed. See https://pkg.go.dev/cuelang.org/go@v0.6.0/cue#Scope


[dodgerblue::bu]Commands:[-::-]

  Items and widgets are controlled through the [dodgerblue]Command Box[-] ([lime]Ctrl-<space>[-],[lime]Ctrl-P[-],[lime]<esc>[-]).
  First, make sure the item you want to change is focused (by [lime]clicking[-] on it).
  Then run any command that is understood by the widget, some are common, others unique.
  The general format for commands is as follows:

  [violet]<command> [lightseagreen]<args> <data source>[-]

  [blue::bu]Data Sources[-::-]
  [violet]<command> [lightseagreen]<eval args and flags>  [-](same as cue and hof)
  [violet]<command> [lightseagreen]https://...            [-](any http json response)
  [violet]<command> [lightseagreen]bash <args...>         [-](any bash json output)
  [violet][lightseagreen]conn <path> <expr>               [-](dot path to item, optional expr)

  [blue::bu]Setting Sources[-::-]          [darkgray](these modify the widget)[-]
  [violet]set [lightseagreen]<data>               [-]sets the scope for the current play
  [violet]add [lightseagreen]<data>               [-]adds another source to the current scope
  [violet]conn [lightseagreen]<path> <expr>       [-]adds a value connection to the current scope which
                           comes from the result of another play or view item
       [lightseagreen]<path>[-]              id/name dotpath to the source item
       [lightseagreen]<expr>[-]              a CUE dotpath to select a subvalue

  [blue::bu]Watching Inputs[-::-]          [darkgray](so you can still edit files your IDE)[-]
  [violet]refresh [lightseagreen]<delay>          [-]change the default refresh time for play edit
  [violet]watch [lightseagreen]                   [-]watch the current item with default delay
  [violet]watch [lightseagreen]<delay>            [-]watch the current item with custom delay
  [violet]globs [lightseagreen]<globs>            [-]manually set watched globs that trigger
                           [darkgray](if watch seems stuck, refresh with [lime]Alt-R[darkgray])[-]
        [lightseagreen]<delay>[-]            a value like 1s or 42ms
        [lightseagreen]<globs>[-]            one or more *.yaml, */*.json, or **/* like args

  [blue::bu]Push and Save[-::-]
  [violet]push                     [-]send editor text to cuelang.org/play
  [violet]write  <file>            [-]save editor text to file
  [violet]export <file>            [-]save final value to file


[dodgerblue::bu]Dashboards:[-::-]

  Dashboards are sets of panels and items along with their
  configuration and data sources. There are hotkeys and commands to
  layout panels and items, as well as save, load, and share them.

  [blue::bu]Panel Management[-::-]
  [lime]Alt-J[-]                    create a new item before
  [lime]Alt-K[-]                    create a new item after
  [lime]Alt-H[-]                    move (swap) item with prev
  [lime]Alt-L[-]                    move (swap) item with next
  [lime]Alt-T[-]                    split the current item
  [lime]Alt-D[-]                    delete the current item
  [lime]Alt-F[-]                    toggle flex direction
  [lime]Alt-P[-]                    toggle panel borders
  [lime]Alt-O[-]                    toggle item borders

  [violet]insert <loc>             [-]insert a new item at <loc>
         <loc>             [-]from: {head,tail,prev,next,index}
  [violet]set.panel.name <name>    [-]set the focused panel name
  [violet]set.item.name  <name>    [-]set the focused item name
  [violet]set.size <int>           [-]set the flex fixed size for an item
  [violet]set.ratio <int>          [-]set the flex ratio for an item

  [violet]<cmd> @<path> <args>[-]     most panel & item commands support @<path>
        @<path>            use this to target a specific panel or item
                           by id/name dotpaths, it must be before args
                           ex: [darkgray]set.size @p0.b2.name 42[-] 


  [blue::bu]Navigation[-::-]
  [lime]click[-]          to focus any box   [darkgray](with mouse)[-]
  [lime]Ctrl-<left>[-]    focus item left    [darkgray](with arrows)[-]
  [lime]Ctrl-<down>[-]    focus item down
  [lime]Ctrl-<up>[-]      focus item up
  [lime]Ctrl-<right>[-]   focus item right
  [lime]Alt-h[-]          focus item left    [darkgray](vim style)[-]
  [lime]Alt-j[-]          focus item down
  [lime]Alt-k[-]          focus item up
  [lime]Alt-l[-]          focus item right

  [blue::bu]Saving and Loading[-::-]
  [violet]list[-]          list available dashboards
  [violet]save [lightseagreen]<dest>[-]   save the current view 
  [violet]show [lightseagreen]<dest>[-]   show the save file
  [violet]load [lightseagreen]<dest>[-]   load and replace current

  [blue::bu]Destinations[-::-]
  [violet]<command> [lightseagreen]./dir/file.cue         [-]path to exact location
  [violet]<command> [lightseagreen]<name>                 [-]module location, in .hof/tui
  [violet]<command> [lightseagreen]@<name>                [-]global, user cache dir, in hof/tui
  [violet]<command> [lightseagreen]https://...            [-]http get with ?id=...
            contents are CUE       post returns the id

  Note, that saving dashboards does not include the data.
  The config and args are saved so that it can be recreated.
  This means any files will have to be present when loaded.


[dodgerblue::bu]Getting Help:[-::-]

  [violet]help[-]          (open this help content in an item)
  [violet]hof feedback[-]  (run this to start a new GitHub issue)

  [gold]Docs:      [deepskyblue]https://docs.hofstadter.io[-]
  [gold]GitHub:    [deepskyblue]https://github.com/hofstadter-io/hof[-]
  [gold]Discord:   [deepskyblue]https://discord.gg/BXwX7n6B8w[-]
  [gold]Slack:     [deepskyblue]https://join.slack.com/t/hofstadter-io/shared_invite/zt-e5f90lmq-u695eJur0zE~AG~njNlT1A[-]

`
