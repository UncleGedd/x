package input

import (
	"regexp"
	"strconv"
)

// MouseButton represents the button that was pressed during a mouse event.
type MouseButton int

// Mouse event buttons
//
// This is based on X11 mouse button codes.
//
//	1 = left button
//	2 = middle button (pressing the scroll wheel)
//	3 = right button
//	4 = turn scroll wheel up
//	5 = turn scroll wheel down
//	6 = push scroll wheel left
//	7 = push scroll wheel right
//	8 = 4th button (aka browser backward button)
//	9 = 5th button (aka browser forward button)
//	10
//	11
//
// Other buttons are not supported.
const (
	MouseButtonNone MouseButton = iota
	MouseButtonLeft
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
	MouseButtonWheelLeft
	MouseButtonWheelRight
	MouseButtonBackward
	MouseButtonForward
	MouseButton10
	MouseButton11
)

var mouseButtons = map[MouseButton]string{
	MouseButtonNone:       "none",
	MouseButtonLeft:       "left",
	MouseButtonMiddle:     "middle",
	MouseButtonRight:      "right",
	MouseButtonWheelUp:    "wheel up",
	MouseButtonWheelDown:  "wheel down",
	MouseButtonWheelLeft:  "wheel left",
	MouseButtonWheelRight: "wheel right",
	MouseButtonBackward:   "backward",
	MouseButtonForward:    "forward",
	MouseButton10:         "button 10",
	MouseButton11:         "button 11",
}

// mouse represents a mouse event.
type mouse struct {
	X, Y   int
	Button MouseButton
	Mod
}

// IsWheel returns true if the mouse event is a wheel event.
func (m mouse) IsWheel() bool {
	return isWheel(m.Button)
}

func isWheel(b MouseButton) bool {
	return b >= MouseButtonWheelUp && b <= MouseButtonWheelRight
}

// String implements fmt.Stringer.
func (m mouse) String() (s string) {
	if m.Mod.IsCtrl() {
		s += "ctrl+"
	}
	if m.Mod.IsAlt() {
		s += "alt+"
	}
	if m.Mod.IsShift() {
		s += "shift+"
	}

	str, ok := mouseButtons[m.Button]
	if !ok {
		s += "unknown"
	} else if str != "none" { // motion events don't have a button
		s += str
	}

	return s
}

// MouseDownEvent represents a mouse button press event.
type MouseDownEvent mouse

// IsWheel returns true if the mouse event is a wheel event.
func (d MouseDownEvent) IsWheel() bool {
	m := mouse(d)
	return m.IsWheel()
}

// String implements fmt.Stringer.
func (d MouseDownEvent) String() (s string) {
	m := mouse(d)
	return m.String()
}

// MouseUpEvent represents a mouse button release event.
type MouseUpEvent mouse

// IsWheel returns true if the mouse event is a wheel event.
func (u MouseUpEvent) IsWheel() bool {
	m := mouse(u)
	return m.IsWheel()
}

// String implements fmt.Stringer.
func (u MouseUpEvent) String() (s string) {
	m := mouse(u)
	return m.String()
}

// MouseMoveEvent represents a mouse motion event.
type MouseMoveEvent mouse

// IsWheel returns true if the mouse event is a wheel event.
func (m MouseMoveEvent) IsWheel() bool {
	mm := mouse(m)
	return mm.IsWheel()
}

// String implements fmt.Stringer.
func (m MouseMoveEvent) String() (s string) {
	mm := mouse(m)
	return mm.String()
}

var mouseSGRRegex = regexp.MustCompile(`(\d+);(\d+);(\d+)([Mm])`)

// Parse SGR-encoded mouse events; SGR extended mouse events. SGR mouse events
// look like:
//
//	ESC [ < Cb ; Cx ; Cy (M or m)
//
// where:
//
//	Cb is the encoded button code
//	Cx is the x-coordinate of the mouse
//	Cy is the y-coordinate of the mouse
//	M is for button press, m is for button release
//
// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Extended-coordinates
func parseSGRMouseEvent(buf []byte) Event {
	str := string(buf[3:])
	matches := mouseSGRRegex.FindStringSubmatch(str)
	if len(matches) != 5 {
		return UnknownCsiEvent(buf)
	}

	b, _ := strconv.Atoi(matches[1])
	px := matches[2]
	py := matches[3]
	release := matches[4] == "m"
	mod, btn, _, isMotion := parseMouseButton(b)
	x, _ := strconv.Atoi(px)
	y, _ := strconv.Atoi(py)

	// (1,1) is the upper left. We subtract 1 to normalize it to (0,0).
	x--
	y--

	// Wheel buttons don't have release events
	// Motion can be reported as a release event in some terminals (Windows Terminal)
	if !isMotion && !isWheel(btn) && release {
		return MouseUpEvent{X: x, Y: y, Button: btn, Mod: mod}
	} else if isMotion {
		return MouseMoveEvent{X: x, Y: y, Button: btn, Mod: mod}
	}
	return MouseDownEvent{X: x, Y: y, Button: btn, Mod: mod}
}

const x10MouseByteOffset = 32

// Parse X10-encoded mouse events; the simplest kind. The last release of X10
// was December 1986, by the way. The original X10 mouse protocol limits the Cx
// and Cy coordinates to 223 (=255-032).
//
// X10 mouse events look like:
//
//	ESC [M Cb Cx Cy
//
// See: http://www.xfree86.org/current/ctlseqs.html#Mouse%20Tracking
func parseX10MouseEvent(buf []byte) Event {
	v := buf[3:6]
	b := int(v[0])
	if b >= x10MouseByteOffset {
		// XXX: b < 32 should be impossible, but we're being defensive.
		b -= x10MouseByteOffset
	}

	mod, btn, isRelease, isMotion := parseMouseButton(b)

	// (1,1) is the upper left. We subtract 1 to normalize it to (0,0).
	x := int(v[1]) - x10MouseByteOffset - 1
	y := int(v[2]) - x10MouseByteOffset - 1

	if isMotion {
		return MouseMoveEvent{X: x, Y: y, Button: btn, Mod: mod}
	} else if isRelease {
		return MouseUpEvent{X: x, Y: y, Button: btn, Mod: mod}
	}
	return MouseDownEvent{X: x, Y: y, Button: btn, Mod: mod}
}

// See: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Extended-coordinates
func parseMouseButton(b int) (mod Mod, btn MouseButton, isRelease bool, isMotion bool) {
	// mouse bit shifts
	const (
		bitShift  = 0b0000_0100
		bitAlt    = 0b0000_1000
		bitCtrl   = 0b0001_0000
		bitMotion = 0b0010_0000
		bitWheel  = 0b0100_0000
		bitAdd    = 0b1000_0000 // additional buttons 8-11

		bitsMask = 0b0000_0011
	)

	// Modifiers
	if b&bitAlt != 0 {
		mod |= Alt
	}
	if b&bitCtrl != 0 {
		mod |= Ctrl
	}
	if b&bitShift != 0 {
		mod |= Shift
	}

	if b&bitAdd != 0 {
		btn = MouseButtonBackward + MouseButton(b&bitsMask)
	} else if b&bitWheel != 0 {
		btn = MouseButtonWheelUp + MouseButton(b&bitsMask)
	} else {
		btn = MouseButtonLeft + MouseButton(b&bitsMask)
		// X10 reports a button release as 0b0000_0011 (3)
		if b&bitsMask == bitsMask {
			btn = MouseButtonNone
			isRelease = true
		}
	}

	// Motion bit doesn't get reported for wheel events.
	if b&bitMotion != 0 && !isWheel(btn) {
		isMotion = true
	}

	return
}
