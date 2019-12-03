package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/ojrac/opensimplex-go"
)

type SimplexGenerator struct {
	frequency float64
	magnitude float64

	octaves     int
	persistence float64

	gen opensimplex.Noise
}

func NewSimplexGenerator(seed int64, frequency, magnitude float64, octaves int, persistence float64) SimplexGenerator {
	gen := opensimplex.New(int64(seed))
	return SimplexGenerator{
		frequency:   frequency,
		magnitude:   magnitude,
		octaves:     octaves,
		persistence: persistence,
		gen:         gen,
	}
}

func (s *SimplexGenerator) D2(x, y float64) float64 {
	total := float64(0)
	tfreq := s.frequency
	tmag := s.magnitude
	tmax := float64(0)
	for i := 0; i < s.octaves; i++ {
		total += s.gen.Eval2(x*tfreq, y*tfreq) * tmag
		tmax += tmag
		tmag *= s.persistence
		tfreq *= 2
	}
	return total
}

func GenerateSimplexTest(size float64, frequency float64, magnitude float64) *imdraw.IMDraw {
	noise := NewSimplexGenerator(10, frequency, magnitude, 6, 0.5)
	imd := imdraw.New(nil)
	for x := float64(0); x < winX; x += size {
		for y := float64(0); y < winY; y += size {
			brightness := noise.D2(x, y)
			imd.Color = pixel.RGB(brightness, brightness, brightness)
			imd.Push(pixel.V(x, y))
			imd.Push(pixel.V(x+size, y+size))
			imd.Rectangle(0)
		}
	}
	return imd
}

func RandomFloat(seed int64) float64 {
	rand.Seed(seed)
	return rand.Float64()
}
