package display

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
	ppu    *PPU
	joypad *Joypad
	screen *ebiten.Image
}

func NewDisplay(ppu *PPU, joypad *Joypad) *Display {
	return &Display{
		ppu:    ppu,
		joypad: joypad,
		screen: ebiten.NewImage(ScreenWidth, ScreenHeight),
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
	d.joypad.Update()
	return nil
}

func (d *Display) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return DisplayWidth, DisplayHeight
}

func (d *Display) updateScreen() {
	d.ppu.mu.Lock()
	defer d.ppu.mu.Unlock()

	framebuffer := d.ppu.GetFrameBuffer()
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			gbColor := framebuffer[y][x]
			rgbaColor := gameBoyColorToRgba(gbColor)
			d.screen.Set(x, y, rgbaColor)
		}
	}
}
func gameBoyColorToRgba(gbColor Color) color.RGBA {
	switch gbColor {
	case 0:
		return color.RGBA{197, 219, 212, 255}
	case 1:
		return color.RGBA{119, 142, 152, 255}
	case 2:
		return color.RGBA{65, 72, 93, 255}
	case 3:
		return color.RGBA{34, 30, 49, 255}
	default:
		return color.RGBA{0, 0, 0, 255}
	}
}
