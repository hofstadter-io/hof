package tview

import "github.com/gdamore/tcell/v2"

// Primitive is the top-most interface for all graphical primitives.
type Primitive interface {
	// return the unique id for the instance
	Id() string

	// Draw draws this primitive onto the screen. Implementers can call the
	// screen's ShowCursor() function but should only do so when they have focus.
	// (They will need to keep track of this themselves.)
	Draw(screen tcell.Screen)

	// GetRect returns the current position of the primitive, x, y, width, and
	// height.
	GetRect() (int, int, int, int)

	// SetRect sets a new position of the primitive.
	SetRect(x, y, width, height int)

	// InputHandler returns a handler which receives key events when it has focus.
	// It is called by the Application class.
	//
	// A value of nil may also be returned, in which case this primitive cannot
	// receive focus and will not process any key events.
	//
	// The handler will receive the key event and a function that allows it to
	// set the focus to a different primitive, so that future key events are sent
	// to that primitive.
	//
	// The Application's Draw() function will be called automatically after the
	// handler returns.
	//
	// The Box class provides functionality to intercept keyboard input. If you
	// subclass from Box, it is recommended that you wrap your handler using
	// Box.WrapInputHandler() so you inherit that functionality.
	InputHandler() func(event *tcell.EventKey, setFocus func(p Primitive))

	EventHandler() func(event tcell.Event, setFocus func(p Primitive))

	// MouseHandler returns a handler which receives mouse events.
	// It is called by the Application class.
	//
	// A value of nil may also be returned to stop the downward propagation of
	// mouse events.
	//
	// The Box class provides functionality to intercept mouse events. If you
	// subclass from Box, it is recommended that you wrap your handler using
	// Box.WrapMouseHandler() so you inherit that functionality.
	MouseHandler() func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive)

	// Focus is called by the application when the primitive receives focus.
	// Implementers may call delegate() to pass the focus on to another primitive.
	Focus(delegate func(p Primitive))

	// HasFocus determines if the primitive has focus. This function must return
	// true also if one of this primitive's child elements has focus.
	HasFocus() bool

	// Blur is called by the application when the primitive loses focus.
	Blur()

	// Mount is a longer term context for bringing a widget into scope
	Mount(context map[string]any) error

	// Refresh(context map[string]any) error

	// Unmount is the opposite of mount
	Unmount() error

	// IsMounted returns true if the primitive is mounted
	IsMounted() bool

	// Render is called when something in the future does so.
	Render() error
	
	// GetProps returns the primitive's prop
	GetProp(prop string) (any, bool)

	// GetProps returns the primitive's props
	GetProps() map[string]any

	// SetProp sets a primitive's props
	// It is up to the implementor to ensure correctness.
	SetProp(props string, value any) error

	// SetProps replaces the primitive's props
	// It is up to the implementor to ensure correctness.
	SetProps(newProps map[string]any) error

	// GetLabels returns the primitive's label
	GetLabel(label string) (string, bool)

	// GetLabels returns the primitive's labels
	GetLabels() map[string]string

	// SetLabel sets a primitive's labels
	// It is up to the implementor to ensure correctness.
	SetLabel(labels string, value string) error

	// SetLabels replaces the primitive's labels
	// It is up to the implementor to ensure correctness.
	SetLabels(newLabels map[string]string) error
}
