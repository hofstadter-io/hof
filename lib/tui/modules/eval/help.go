package eval

const evalHelpText = `
Hof Tui Controls:  (scroll for more)

Legend:
  A = ALT      mouse & click
  C = CTRL     also supported
  M = META
  S = SHIFT    <esc> = unfocus

App:
  C-<space>    focus the command bar
  A-<space>    focus the main content
  A-?          show / hide help for item 
  A-/          show / hide the console
  A-z          show / hide the error log

Panels:
  // navigating between panels & items
	you can also click on a panel to focus
  (planned) A-[hjkl] (vim style movement)

  // adding & moving panels
  A-[JK]       create item           [before/after]
  A-[HL]       move item             [prev/next]
  A-T          split item or panel, making new stuff
  A-D          delete  the current item

  // display options
  A-F          flip panel flex direction
  A-P          show borders on panels
  A-O          show borders on items

Items:
  A-w          change the current widget
  A-e          edit the current settings
  A-c          connect the current item
  A-?          show / hide help for item 
`
