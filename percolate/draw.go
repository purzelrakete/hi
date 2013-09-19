package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	Black = color.Gray16{0}
	Blue  = color.RGBA{0, 0, 255, 255}
	Grey  = color.Gray16{0x1111}
)

func Draw(m *Material, open []bool, filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, m.W, m.H))
	for y := 0; y < m.W; y++ {
		for x := 0; x < m.H; x++ {
			i := y*m.W + x
			switch {
			case open[i]:
				img.Set(x, y, Blue)
			case m.Vacant(i):
				img.Set(x, y, Grey)
			default:
				img.Set(x, y, Black)
			}
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
