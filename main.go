package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// --------------------------------------------------
// Hyperparameters
// --------------------------------------------------

// Window
const winX = 3840
const winY = 2160

// Stars
const starsSeed float64 = 10.0
const starsFreq float64 = 0.02
const starsMag float64 = 1.0
const starsOct int = 6
const starsPers float64 = 0.5

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

// Camera
var cameraX float64 = winX / 2
var cameraY float64 = winY / 2
var cameraZoom float64 = 0.3
var cameraZoomX float64 = 0
var cameraZoomY float64 = 0

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Galactica",
		Bounds: pixel.R(0, 0, winX, winY),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	musicControl := NewMusicControl()
	musicControl.Init()
	//musicControl.Play()
	defer musicControl.Close()

	inputControl := NewInputControl()

	rand.Seed(time.Now().UnixNano())

	win.Clear(colornames.Black)
	stars := generateStars(1000)
	background := generateBackground()
	solar := NewRandomSolarSystem(12345)
	solar.Build()

	//test := GenerateSimplexTest(10, 0.015, 1.0)
	//test.Draw(win)

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape))

		inputControl.UpdateCamera(win)

		// Clear window
		win.Clear(colornames.Black)

		// Draw background and stars
		background.Draw(win)
		stars.Draw(win)

		// Update and draw planet
		solar.Update()
		solar.Draw(win)

		if win.JustPressed(pixelgl.KeyEnter) {
			win.Clear(colornames.Black)
			stars = generateStars(1000)
			background = generateBackground()
			solar = NewRandomSolarSystem(rand.Int63())
			solar.Build()
		}

		win.Update()
	}
}

func main() {
	println("Starting...")
	pixelgl.Run(run)
}

func generateStars(num int) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	noise := NewSimplexGenerator(universeSeed, starsFreq, starsMag, starsOct, starsPers)
	for i := 0; i < num; i++ {
		x := rand.Float64() * float64(winX)
		y := rand.Float64() * float64(winY)

		brightness := noise.D2(x, y)
		radius := solarPixelSize //rand.Float64() * 5

		imd.Color = pixel.RGB(brightness, brightness, brightness)
		imd.Push(pixel.V(x, y))
		imd.Push(pixel.V(x+radius, y+radius))
		imd.Rectangle(0)
	}

	return imd
}

func generateBackground() *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(0, 0, 0)
	imd.Push(pixel.V(0, 0))
	imd.Push(pixel.V(winX, winY))
	imd.Rectangle(0)
	return imd
}

func distsquare(x1, y1, x2, y2 float64) float64 {
	return (x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)
}

func scaleRGB(r, g, b int, factor float64) color.Color {
	return pixel.RGB(factor*float64(r)/255.0, factor*float64(g)/255.0, factor*float64(b)/255.0)
}
