package spectrogram

import (
	"bytes"
	colorful "github.com/lucasb-eyer/go-colorful"
	"image"
	"image/png"
	"math"
	"os"
)

var (
	// colormap for spectrogram. lowest fqs are blue, highest are red.
	colormap = gradientTable{
		{mustParseHex("#3d1ecc"), 0.0},
		{mustParseHex("#3288bd"), 0.1},
		{mustParseHex("#66c2a5"), 0.2},
		{mustParseHex("#abdda4"), 0.3},
		{mustParseHex("#e6f598"), 0.4},
		{mustParseHex("#ffffbf"), 0.5},
		{mustParseHex("#fee090"), 0.6},
		{mustParseHex("#fdae61"), 0.7},
		{mustParseHex("#f46d43"), 0.8},
		{mustParseHex("#d53e4f"), 0.9},
		{mustParseHex("#9e0142"), 1.0},
	}
)

// Draw spectrogram and save to filename
func Draw(s *Spectrogram, filename string) error {
	return DrawXY(s, filename, 100, 100)
}

// TODO(cs): take maximum / average from range
// DrawXY spectrogram with certain height and width and save to filename
func DrawXY(s *Spectrogram, filename string, height, width int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	nfreqs, ntimes := float64((*s).NumFreqSlots()), float64((*s).NumTimeSlots())

	// move to log scale for frequencies
	nfreqsLog := math.Log(nfreqs)

	// take the maximum only over points that are plotted
	max := -math.MaxFloat64
	for x := 1; x <= width; x++ {
		// use rule of three to project (linear) time slots to width
		i := ceil(float64(x)*ntimes/(float64(width))) - 1
		for y := 1; y <= height; y++ {
			// use rule of three to project (log) frequency slots to height
			jLog := float64(y) * nfreqsLog / (float64(height))
			j := ceil(math.Exp(jLog)) - 1

			amplitude := (*s).data[i][j]
			if amplitude > max {
				max = amplitude
			}
		}
	}

	// draw log-scaled freqs
	for x := 1; x <= width; x++ {
		// use rule of three to project (linear) time slots to width
		i := ceil(float64(x)*ntimes/(float64(width))) - 1
		for y := 1; y <= height; y++ {
			// use rule of three to project (log) frequency slots to height
			jLog := float64(y) * nfreqsLog / (float64(height))
			j := ceil(math.Exp(jLog)) - 1

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

func mustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}

	return c
}

// round an float following the mathematical laws
func ceil(f float64) int {
	return int(math.Ceil(f))
}
