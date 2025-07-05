package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Scale         = 3
	DisplayWidth  = ScreenWidth * Scale
	DisplayHeight = ScreenHeight * Scale
)

type Display struct {
	framebuffer *[ScreenHeight][ScreenWidth]Color
	screen      *ebiten.Image
}

func NewDisplay(framebuffer *[ScreenHeight][ScreenWidth]Color) *Display {
	return &Display{
		framebuffer: framebuffer,
		screen:      ebiten.NewImage(ScreenWidth, ScreenHeight),
	}
}

func RunDisplay(display *Display) {
	ebiten.SetWindowSize(DisplayWidth, DisplayHeight)
	ebiten.SetWindowTitle("Garboy")

	if err := ebiten.RunGame(display); err != nil {
		panic("Error when running display")
	}
}

func (d *Display) Draw(screen *ebiten.Image) {
	d.updateScreen()

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(Scale, Scale)
	screen.DrawImage(d.screen, options)
}

func (d *Display) Update() error {
	return nil // Step is handled by the Scheduler
}

func (d *Display) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return DisplayWidth, DisplayHeight
}

func (d *Display) updateScreen() {
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			gbColor := d.framebuffer[y][x]
			rgbaColor := gameBoyColorToRgba(gbColor)
			d.screen.Set(x, y, rgbaColor)
		}
	}
}
func gameBoyColorToRgba(gbColor Color) color.RGBA {
	switch gbColor {
	case 0:
		return color.RGBA{155, 188, 15, 255}
	case 1:
		return color.RGBA{139, 172, 15, 255}
	case 2:
		return color.RGBA{48, 98, 48, 255}
	case 3:
		return color.RGBA{15, 56, 15, 255}
	default:
		return color.RGBA{0, 0, 0, 255}
	}
}
