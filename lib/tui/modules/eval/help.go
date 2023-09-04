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
	// (TBD) navigating between panels & items
	// (click on the panel you want for now)
	A-[hjkl]     move left,down,up,right (vim)
	C-<arrow>    same as previus

	// adding & moving panels
	A-[JK]       create item           [before/after]
	A-[HL]       move item             [prev/next]
	C-S-<arrow>  same as previus two
	A-T          split item or panel, making new stuff

	// display options
	A-[PO]       show borders on [panels/items]
	A-F          flip panel flex direction

Items: (WIP)
	A-E          edit    the current item
	A-D          delete  the current item
	A-C          connect the current item
`
