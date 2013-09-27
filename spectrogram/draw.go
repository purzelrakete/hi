package spectrogram

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

// Draw a spectrogram and save to filename
func Draw(spectrogram [][]complex128, width int, filename string) error {
	height := len(spectrogram)
	max := 0.0
	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			pixel := spectrogram[x][y]
			if cmplx.Abs(pixel) > max {
				max = cmplx.Abs(pixel)
			}
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			intensity := uint8(255 * cmplx.Abs(spectrogram[x][y]) / max)
			if y == 600 {
				fmt.Println(intensity)
			}
			img.Set(x, y, color.RGBA{0, 0, 255, intensity})
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
