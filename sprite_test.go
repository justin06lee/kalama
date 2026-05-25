package shaw

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// makePNG builds a 2x1 PNG: pixel (0,0) opaque red, pixel (1,0) transparent.
func makePNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.NRGBA{R: 255, A: 255})
	img.Set(1, 0, color.NRGBA{}) // A==0, transparent
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

func TestLoadSpriteDimensionsAndTransparency(t *testing.T) {
	s, err := LoadSprite(bytes.NewReader(makePNG(t)))
	if err != nil {
		t.Fatalf("LoadSprite: %v", err)
	}
	if s.Width() != 2 || s.Height() != 1 {
		t.Fatalf("dims = %dx%d, want 2x1", s.Width(), s.Height())
	}
	if got := s.at(0, 0); got != (Color{R: 255, A: 255}) {
		t.Errorf("opaque pixel = %+v, want red A=255", got)
	}
	if got := s.at(1, 0); got.A != 0 {
		t.Errorf("transparent pixel A = %d, want 0", got.A)
	}
}

func TestLoadSpriteRejectsNonPNG(t *testing.T) {
	if _, err := LoadSprite(bytes.NewReader([]byte("not a png"))); err == nil {
		t.Fatal("expected error decoding non-PNG, got nil")
	}
}
