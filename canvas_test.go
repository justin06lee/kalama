package kalama

import (
	"bytes"
	"testing"
)

func TestNewCanvasForcesEvenHeight(t *testing.T) {
	c := NewCanvas(4, 3)
	if c.Width() != 4 {
		t.Errorf("width = %d, want 4", c.Width())
	}
	if c.Height() != 4 {
		t.Errorf("height = %d, want 4 (rounded up to even)", c.Height())
	}
}

func TestSetSkipsTransparentAndOutOfBounds(t *testing.T) {
	c := NewCanvas(2, 2)
	red := Color{R: 255, A: 255}
	c.Set(0, 0, red)
	c.Set(0, 1, Color{R: 9, G: 9, B: 9}) // A==0 -> ignored
	c.Set(5, 5, red)                     // out of bounds -> ignored
	if got := c.at(0, 0); got != red {
		t.Errorf("at(0,0) = %+v, want %+v", got, red)
	}
	if got := c.at(0, 1); got != (Color{}) {
		t.Errorf("at(0,1) = %+v, want zero (transparent set ignored)", got)
	}
}

func TestClearFillsEveryPixel(t *testing.T) {
	c := NewCanvas(2, 2)
	blue := Color{B: 255, A: 255}
	c.Clear(blue)
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			if got := c.at(x, y); got != blue {
				t.Errorf("at(%d,%d) = %+v, want %+v", x, y, got, blue)
			}
		}
	}
}

func TestBlitSkipsTransparentAndClips(t *testing.T) {
	s, err := LoadSprite(bytes.NewReader(makePNG(t))) // 2x1: red, transparent
	if err != nil {
		t.Fatalf("LoadSprite: %v", err)
	}
	c := NewCanvas(2, 2)
	c.Clear(Color{B: 255, A: 255}) // blue background
	c.Blit(s, 0, 0)
	if got := c.at(0, 0); got != (Color{R: 255, A: 255}) {
		t.Errorf("at(0,0) = %+v, want red (opaque sprite pixel)", got)
	}
	if got := c.at(1, 0); got != (Color{B: 255, A: 255}) {
		t.Errorf("at(1,0) = %+v, want blue (transparent sprite pixel skipped)", got)
	}
	c.Blit(s, 5, 5) // fully off-canvas: must not panic
}

func TestRenderHalfBlockOneCell(t *testing.T) {
	c := NewCanvas(1, 2)
	c.Set(0, 0, Color{R: 255, A: 255}) // top pixel red
	c.Set(0, 1, Color{B: 255, A: 255}) // bottom pixel blue
	got := c.Render()
	want := "\x1b[38;2;255;0;0m\x1b[48;2;0;0;255m▀\x1b[0m"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRenderTransparentPixelIsBlack(t *testing.T) {
	c := NewCanvas(1, 2) // nothing set -> both pixels transparent
	got := c.Render()
	want := "\x1b[38;2;0;0;0m\x1b[48;2;0;0;0m▀\x1b[0m"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRenderTwoRowsSeparatedByNewline(t *testing.T) {
	c := NewCanvas(1, 4) // 2 cell rows
	got := c.Render()
	black := "\x1b[38;2;0;0;0m\x1b[48;2;0;0;0m▀"
	want := black + "\x1b[0m\n" + black + "\x1b[0m"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}
