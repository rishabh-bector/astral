package main

import (
	"image/color"
	"math"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

type Planet struct {
	name string
	seed int64

	// sizing
	minisize  float64
	radius    float64
	tolerance float64

	// altitudes
	cols [][]int

	// positions (solar view)
	x float64
	y float64

	orbitX        float64
	orbitY        float64
	orbitVelocity float64
	orbitStarting float64

	t   float64
	mat pixel.Matrix

	// planetary noise generator
	gen SimplexGenerator

	// planet view target
	imd *imdraw.IMDraw

	// solar view target
	mini  *imdraw.IMDraw
	orbit *imdraw.IMDraw

	// solar view memory
	miniMem [][]color.Color
}

func NewStandardPlanet(seed float64) Planet {
	return Planet{
		name:      "standard",
		radius:    400,
		tolerance: 75,
		gen:       NewSimplexGenerator(universeSeed, planetFreq, planetMag, planetOct, planetPers),
	}
}

func NewRandomizedPlanet(seed int64, ox, oy float64) Planet {
	f := planetFreq + (planetFreqR * (RandomFloat(seed+10)*2 - 1))
	m := planetMag + (planetMagR * (RandomFloat(seed+20)*2 - 1))
	p := planetPers + (planetPersR * (RandomFloat(seed+30)*2 - 1))
	s := planetMiniSize + (planetMiniSizeR * (RandomFloat(seed+40)*2 - 1))
	o := int(float64(planetOct) + (planetOctR * (RandomFloat(seed+50)*2 - 1)))

	//println(f, m, o, p, s)

	cols := planetColors

	ncols := make([][]int, len(planetColors))
	for i := 0; i < len(planetColors); i++ {
		ncols[i] = make([]int, len(planetColors[i]))
	}

	for i := 0; i < len(cols); i++ {
		ncols[i] = distColors(seed+int64(i), cols[i], planetColorInterpolation)
	}

	ov := solarOrbitSpeed + (solarOrbitSpeedR * (RandomFloat(seed+60) - 1))
	os := solarOrbitStarting + (solarOrbitStartingR * (RandomFloat(seed+70) - 1))

	return Planet{
		name:          "standard",
		seed:          seed,
		radius:        planetRad,
		tolerance:     planetTol,
		minisize:      s,
		gen:           NewSimplexGenerator(seed, f, m, o, p),
		cols:          ncols,
		orbitX:        ox,
		orbitY:        oy,
		orbitVelocity: ov,
		t:             os,
	}
}

func (p *Planet) Build() {
	p.imd = imdraw.New(nil)
	p.imd.SetMatrix(pixel.IM)

	largeradius := (p.radius * p.radius) + (p.tolerance * p.tolerance)
	smallradius := (p.radius * p.radius) - (p.tolerance * p.tolerance)

	cx := float64(winX / 2)
	cy := float64(winY / 2)

	// atmostphere
	p.imd.Color = pixel.ToRGBA(scaleRGB(p.cols[0][0], p.cols[0][1], p.cols[0][2], 1.0))
	for x := cx - p.radius - p.tolerance + 10; x < cx+p.radius+p.tolerance; x += 10 {
		for y := cy - p.radius - p.tolerance; y < cx+p.radius+p.tolerance; y += 10 {
			radial := distsquare(x, y, cx, cy)
			if radial <= largeradius && radial >= smallradius-100 {
				p.imd.Push(pixel.V(x, y))
				p.imd.Push(pixel.V(x+10, y+10))
				p.imd.Rectangle(0)
			}
		}
	}

	// Color
	for x := cx - p.radius - p.tolerance; x < cx+p.radius+p.tolerance; x += 10 {
		for y := cy - p.radius - p.tolerance; y < cy+p.radius+p.tolerance; y += 10 {
			radial := distsquare(x, y, cx, cy)
			if radial < smallradius {
				n := p.gen.D2(x, y)

				col := p.genColor(n, p.cols)
				if len(col) < 3 {
					println(p.cols)
				}
				p.imd.Color = pixel.ToRGBA(scaleRGB(col[0], col[1], col[2], 1.0))

				p.imd.Push(pixel.V(x, y))
				p.imd.Push(pixel.V(x+10, y+10))
				p.imd.Rectangle(0)
			}
		}
	}
}

func (p *Planet) UpdateOrbit(timestep float64) {
	p.t += timestep*p.orbitVelocity + p.orbitStarting
	p.x = p.orbitX * math.Sin(p.t)
	p.y = p.orbitY * math.Cos(p.t)
}

func (p *Planet) BuildMini() {
	p.orbit = imdraw.New(nil)

	//largeradius := (p.radius * p.radius * p.minisize) + (p.tolerance * p.tolerance * p.minisize)
	smallradius := (p.radius * p.radius * p.minisize) - (p.tolerance * p.tolerance * p.minisize)

	// slice center
	tc := (p.radius + p.tolerance) / solarPixelSize

	// slice dimensions
	sd := (2*p.radius + 2*p.tolerance) / solarPixelSize
	p.miniMem = make([][]color.Color, int(sd))
	for i := 0; i < int(sd); i++ {
		p.miniMem[i] = make([]color.Color, int(sd))
	}

	for x := float64(0); x < sd; x++ {
		for y := float64(0); y < sd; y++ {
			radial := distsquare(x*solarPixelSize, y*solarPixelSize, tc*solarPixelSize, tc*solarPixelSize)
			if radial < smallradius {
				n := p.gen.D2(x*5, y*5) // this is meant to be 5.0 for calibration
				gc := p.genColor(n, p.cols)
				col := pixel.ToRGBA(scaleRGB(gc[0], gc[1], gc[2], 1.0))
				p.miniMem[int(x)][int(y)] = col
			} else {
				p.miniMem[int(x)][int(y)] = nil
			}
		}
	}
}

func (p *Planet) Draw(win *pixelgl.Window) {
	p.imd.Draw(win)
}

func (p *Planet) DrawMini(win *pixelgl.Window) {
	p.mini = imdraw.New(nil)

	sd := int((2*p.radius + 2*p.tolerance) / solarPixelSize)

	mat := pixel.IM
	mat = mat.Moved(pixel.V(p.x-(p.radius+p.tolerance), p.y-(p.radius+p.tolerance)))
	mat = mat.Moved(pixel.V(cameraX, cameraY))
	mat = mat.Scaled(pixel.V(cameraZoomX, cameraZoomY), cameraZoom)
	p.mini.SetMatrix(mat)

	p.mini.Color = colornames.White

	for x := 0; x < sd; x++ {
		for y := 0; y < sd; y++ {
			if p.miniMem[x][y] != nil {
				p.mini.Color = p.miniMem[x][y]

				p.mini.Push(pixel.V(float64(x)*solarPixelSize, float64(y)*solarPixelSize))
				p.mini.Push(pixel.V(float64(x)*solarPixelSize+solarPixelSize, float64(y)*solarPixelSize+solarPixelSize))
				p.mini.Rectangle(0)
			}
		}
	}

	p.DrawOrbit(win)
	p.mini.Draw(win)
}

func (p *Planet) DrawOrbit(win *pixelgl.Window) {
	// orbital ellipse
	p.orbit = imdraw.New(nil)

	mat := pixel.IM
	mat = mat.Moved(pixel.V(cameraX, cameraY))
	mat = mat.Scaled(pixel.V(cameraZoomX, cameraZoomY), cameraZoom)
	p.orbit.SetMatrix(mat)

	c := pixel.ToRGBA(colornames.Darkslategrey)
	c.A = 0.1
	p.orbit.Color = c
	p.orbit.Push(pixel.V(0, 0))
	p.orbit.Ellipse(pixel.V(p.orbitX, p.orbitY), 2)

	p.orbit.Draw(win)
}

func (p *Planet) genColor(a float64, cols [][]int) []int {
	t := float64(0)
	for i := 0; i < len(planetAltitudes); i++ {
		t += planetAltitudes[i]
	}

	ind := math.Abs(((a + 0.8) / 2) * t)
	if ind < 0 {
		ind = 0.01
	}
	if ind >= t {
		ind = t - 0.01
	}

	t2 := float64(0)
	for i := 0; i < len(planetAltitudes); i++ {
		t2 += planetAltitudes[i]
		if ind < t2 {
			return cols[i]
		}
	}

	return nil
}

func distColors(seed int64, s []int, fac float64) []int {

	f1 := RandomFloat(seed + 80)
	f2 := RandomFloat(seed + 90)
	f3 := RandomFloat(seed + 100)

	rc := []int{
		int((f1 * 255.0)),
		int((f2 * 255.0)),
		int((f3 * 255.0)),
	}

	rcd := []int{
		rc[0] - s[0],
		rc[1] - s[1],
		rc[2] - s[2],
	}

	return []int{
		s[0] + int(fac*float64(rcd[0])),
		s[1] + int(fac*float64(rcd[1])),
		s[2] + int(fac*float64(rcd[2])),
	}
}

func absInt(n int) int {
	return int(math.Abs(float64(n)))
}
