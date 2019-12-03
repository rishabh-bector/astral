package main

import (
	"github.com/faiface/pixel/pixelgl"
)

type InputControl struct {
}

func NewInputControl() InputControl {
	return InputControl{}
}

func (i *InputControl) UpdateCamera(win *pixelgl.Window) {
	if win.Pressed(pixelgl.MouseButtonLeft) {
		deltaV := win.MousePosition().Sub(win.MousePreviousPosition())
		cameraX += deltaV.X * cameraSpeed
		cameraY += deltaV.Y * cameraSpeed
	}

	deltaS := win.MouseScroll()
	if deltaS.Y != 0 {
		cameraZoomX = winX / 2
		cameraZoomY = winY / 2
	}

	cameraZoom += deltaS.Y * cameraZoomSpeed
	if cameraZoom < cameraMinZoom {
		cameraZoom = cameraMinZoom
	}
	if cameraZoom > cameraMaxZoom {
		cameraZoom = cameraMaxZoom
	}
}
