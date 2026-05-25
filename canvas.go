package shaw

import (
	"fmt"
	"strings"
)

// Canvas is a width x height grid of pixels. Height is always even: each
// terminal cell row holds two vertical pixels (see Render).
type Canvas struct {
	w, h   int
	pixels []Color
}

// NewCanvas allocates a w x h canvas. An odd h is rounded up to the next even
// number so it maps cleanly onto half-block cell rows.
func NewCanvas(w, h int) *Canvas {
	if h%2 != 0 {
		h++
	}
	return &Canvas{w: w, h: h, pixels: make([]Color, w*h)}
}

// Width returns the canvas width in pixels.
func (c *Canvas) Width() int { return c.w }

// Height returns the canvas height in pixels (even).
func (c *Canvas) Height() int { return c.h }

// Set paints one pixel. Transparent colors (A == 0) and out-of-bounds
// coordinates are ignored, so callers can blit freely without clipping.
func (c *Canvas) Set(x, y int, col Color) {
	if col.A == 0 || x < 0 || y < 0 || x >= c.w || y >= c.h {
		return
	}
	c.pixels[y*c.w+x] = col
}

// Clear fills every pixel with col, bypassing the transparency skip in Set.
func (c *Canvas) Clear(col Color) {
	for i := range c.pixels {
		c.pixels[i] = col
	}
}

// at returns the stored color at (x,y). Test helper; assumes in bounds.
func (c *Canvas) at(x, y int) Color { return c.pixels[y*c.w+x] }

// Blit draws sprite s with its top-left corner at (x,y). Transparent sprite
// pixels are skipped and pixels outside the canvas are clipped, both via Set.
func (c *Canvas) Blit(s *Sprite, x, y int) {
	for sy := 0; sy < s.h; sy++ {
		for sx := 0; sx < s.w; sx++ {
			c.Set(x+sx, y+sy, s.pixels[sy*s.w+sx])
		}
	}
}

// Render draws the canvas as an ANSI truecolor string. Each terminal cell is
// the upper-half-block glyph ▀ with the top pixel as foreground and the bottom
// pixel as background, giving two vertical pixels per cell. Cell rows are joined
// by newlines with no trailing newline; each row ends with an SGR reset.
func (c *Canvas) Render() string {
	var b strings.Builder
	rows := c.h / 2
	for r := 0; r < rows; r++ {
		for x := 0; x < c.w; x++ {
			top := c.opaque(x, 2*r)
			bot := c.opaque(x, 2*r+1)
			fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀",
				top.R, top.G, top.B, bot.R, bot.G, bot.B)
		}
		b.WriteString("\x1b[0m")
		if r < rows-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// opaque returns the color at (x,y) for rendering: transparent pixels render as
// black so every cell has a concrete foreground and background color.
func (c *Canvas) opaque(x, y int) Color {
	px := c.pixels[y*c.w+x]
	if px.A == 0 {
		return Color{}
	}
	return px
}
