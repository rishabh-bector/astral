package main

import (
	"math/rand"

	"github.com/faiface/pixel/pixelgl"
)

type SolarSystem struct {
	seed    int64
	planets []Planet
}

func NewRandomSolarSystem(seed int64) SolarSystem {
	np := int(float64(numPlanets) + (float64(numPlanetsR) * (rand.Float64()*2 - 1)))
	pls := []Planet{}

	d := solarOrbitMinDistance
	r := d

	for i := 0; i < np; i++ {
		pls = append(pls, NewRandomizedPlanet(seed+int64(i*1000), r, r))
		d *= solarOrbitStepDistance
		r = r + solarOrbitStepDistance + (solarOrbitStepDistanceR * (RandomFloat(seed+int64(i*10))*2 - 1))
	}

	return SolarSystem{
		seed:    seed,
		planets: pls,
	}
}

func (s *SolarSystem) Build() {
	for i := 0; i < len(s.planets); i++ {
		s.planets[i].BuildMini()
	}
}

func (s *SolarSystem) Update() {
	for i := 0; i < len(s.planets); i++ {
		s.planets[i].UpdateOrbit(0.001)
	}
}

func (s *SolarSystem) Draw(win *pixelgl.Window) {
	for i := 0; i < len(s.planets); i++ {
		s.planets[i].DrawMini(win)
	}
}
