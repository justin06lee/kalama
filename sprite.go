package shaw

import (
	"image/color"
	"image/png"
	"io"
)

// Sprite is a fixed pixel grid loaded from a PNG. Pixels with alpha 0 are
// transparent and are skipped when blitted onto a canvas.
type Sprite struct {
	w, h   int
	pixels []Color
}

// Width returns the sprite width in pixels.
func (s *Sprite) Width() int { return s.w }

// Height returns the sprite height in pixels.
func (s *Sprite) Height() int { return s.h }

// at returns the color at (x,y). Test helper; assumes in bounds.
func (s *Sprite) at(x, y int) Color { return s.pixels[y*s.w+x] }

// LoadSprite decodes a PNG into a Sprite. Fully transparent source pixels
// (alpha 0) become transparent Colors; all others become opaque (A = 255).
func LoadSprite(r io.Reader) (*Sprite, error) {
	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	pixels := make([]Color, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			n := color.NRGBAModel.Convert(img.At(b.Min.X+x, b.Min.Y+y)).(color.NRGBA)
			if n.A == 0 {
				pixels[y*w+x] = Color{} // transparent
				continue
			}
			pixels[y*w+x] = Color{R: n.R, G: n.G, B: n.B, A: 255}
		}
	}
	return &Sprite{w: w, h: h, pixels: pixels}, nil
}
