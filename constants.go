package main

import "time"

// EPIC
var universeSeed int64 = time.Now().UnixNano()

// --------------------------------------------------
// Input parameters
// --------------------------------------------------

const cameraSpeed float64 = 0.5

const cameraZoomSpeed float64 = 0.04
const cameraMaxZoom float64 = 4.0
const cameraMinZoom float64 = 0.2

// --------------------------------------------------
// Base solar system parameters
// --------------------------------------------------

const solarPixelSize float64 = 2

const numPlanets int = 6
const numPlanetsR int = 3

const solarOrbitMinDistance float64 = 400
const solarOrbitStepDistance float64 = 300
const solarOrbitStepDistanceR float64 = 150

const solarOrbitSpeed float64 = 0.5
const solarOrbitSpeedR float64 = 0.3

const solarOrbitStarting float64 = 0
const solarOrbitStartingR float64 = 10000

// --------------------------------------------------
// Base planet parameters
// --------------------------------------------------

// Planet orbit
const planetOrbitX float64 = 1000
const planetOrbitY float64 = 1000

// Planet geometry
const planetRad float64 = 400
const planetMiniSize float64 = 0.05
const planetTol float64 = 75

// Land generation
const planetFreq float64 = 0.008
const planetMag float64 = 1.0
const planetOct int = 6
const planetPers float64 = 0.5

// Altitude generation
var planetAltitudes = []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
var planetColors [][]int = [][]int{
	{15, 56, 104},
	{15, 96, 200},
	{211, 190, 153},
	{68, 116, 84},
	{143, 68, 38},
	{240, 240, 240},
}
var planetColorInterpolation float64 = 0.6

// Randomized parameter ranges
const planetMiniSizeR float64 = 0.04
const planetFreqR float64 = 0.0035
const planetMagR float64 = 0.2
const planetOctR float64 = 3
const planetPersR float64 = 0.25
