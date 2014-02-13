package spectrogram

import (
	"bytes"
	colorful "github.com/lucasb-eyer/go-colorful"
	"image"
	"image/png"
	"math"
	"os"
)

// Draw spectrogram and save to filename
func Draw(s *Spectrogram, filename string) error {
	return DrawXY(s, filename, 100, 100)
}

// DrawXY spectrogram with certain height and width and save to filename
func DrawXY(s *Spectrogram, filename string, height, width int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	colormap := gradientTable{
		{MustParseHex("#3d1ecc"), 0.0},
		{MustParseHex("#3288bd"), 0.1},
		{MustParseHex("#66c2a5"), 0.2},
		{MustParseHex("#abdda4"), 0.3},
		{MustParseHex("#e6f598"), 0.4},
		{MustParseHex("#ffffbf"), 0.5},
		{MustParseHex("#fee090"), 0.6},
		{MustParseHex("#fdae61"), 0.7},
		{MustParseHex("#f46d43"), 0.8},
		{MustParseHex("#d53e4f"), 0.9},
		{MustParseHex("#9e0142"), 1.0},
	}

	// take the maximum only over points that are plotted
	max := -math.MaxFloat64
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// rescale number of slots to height and width
			i := int(float64(x*(*s).NumTimeSlots())/float64(width) + 0.5)
			j := int(math.Pow(2.0, (math.Log2(float64((*s).NumFreqSlots()))/float64(height)*float64(y)))+0.5) - 1
			amplitude := (*s).data[i][j]
			if amplitude > max {
				max = amplitude
			}
		}
	}

	// draw with log-scaled freqs
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// rescale number of slots to height and width
			i := int(float64(x*(*s).NumTimeSlots())/float64(width) + 0.5)
			j := int(math.Pow(2.0, (math.Log2(float64((*s).NumFreqSlots()))/float64(height)*float64(y)))+0.5) - 1
			intensity := (*s).data[i][j] / max
			img.Set(x, height-y-1, colormap.getInterpolatedColorFor(intensity))
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

// The position of each keypoint has to live in the range [0,1]
type gradientTable []struct {
	col colorful.Color
	pos float64
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
func (g gradientTable) getInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(g)-1; i++ {
		c1 := g[i]
		c2 := g[i+1]
		if c1.pos <= t && t <= c2.pos {
			t := (t - c1.pos) / (c2.pos - c1.pos)
			return c1.col.BlendHcl(c2.col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return g[0].col
}

// This is a very nice thing Golang forces you to do!
// It is necessary so that we can write out the literal of the colortable below.
func MustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return c
}
