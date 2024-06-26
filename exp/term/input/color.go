package input

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// ForegroundColorEvent represents a foreground color change event.
type ForegroundColorEvent struct{ color.Color }

// String implements fmt.Stringer.
func (e ForegroundColorEvent) String() string {
	return colorToHex(e)
}

// BackgroundColorEvent represents a background color change event.
type BackgroundColorEvent struct{ color.Color }

// String implements fmt.Stringer.
func (e BackgroundColorEvent) String() string {
	return colorToHex(e)
}

// CursorColorEvent represents a cursor color change event.
type CursorColorEvent struct{ color.Color }

// String implements fmt.Stringer.
func (e CursorColorEvent) String() string {
	return colorToHex(e)
}

func colorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func xParseColor(s string) color.Color {
	switch {
	case strings.HasPrefix(s, "rgb:"):
		parts := strings.Split(s[4:], "/")
		if len(parts) != 3 {
			return color.Black
		}

		r, _ := strconv.ParseUint(parts[0], 16, 32)
		g, _ := strconv.ParseUint(parts[1], 16, 32)
		b, _ := strconv.ParseUint(parts[2], 16, 32)

		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	case strings.HasPrefix(s, "rgba:"):
		parts := strings.Split(s[5:], "/")
		if len(parts) != 4 {
			return color.Black
		}

		r, _ := strconv.ParseUint(parts[0], 16, 32)
		g, _ := strconv.ParseUint(parts[1], 16, 32)
		b, _ := strconv.ParseUint(parts[2], 16, 32)
		a, _ := strconv.ParseUint(parts[3], 16, 32)

		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	return color.Black
}
