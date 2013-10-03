package spectrogram

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Draw spectrogram and save to filename
func Draw(s *Spectrogram, filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, s.Width(), s.Height()))
	for y := range *s {
		for x := range (*s)[y] {
			img.Set(x, y, color.RGBA{0, 0, 255, (*s)[y][x]})
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()
	f.Write(buf.Bytes())

	return nil
}
