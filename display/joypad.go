package display

import (
	"github.com/hajimehoshi/ebiten/v2"

	"garboy/utils"
)

const (
	JoypadRight  = 0
	JoypadLeft   = 1
	JoypadUp     = 2
	JoypadDown   = 3
	JoypadA      = 4
	JoypadB      = 5
	JoypadSelect = 6
	JoypadStart  = 7

	SelectDirectionKeys = 4
	SelectButtonKeys    = 5
)

type Joypad struct {
	strobe      uint8
	buttonState uint8

	joyp uint8
}

func NewJoypad() *Joypad {
	return &Joypad{
		buttonState: 0xFF,
		joyp:        0xFF,
	}
}

func (j *Joypad) Update() {
	j.buttonState = 0xFF

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		j.buttonState = utils.ResetBit(j.buttonState, JoypadRight)
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		j.buttonState = utils.ResetBit(j.buttonState, JoypadLeft)
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		j.buttonState = utils.ResetBit(j.buttonState, JoypadUp)
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		j.buttonState = utils.ResetBit(j.buttonState, JoypadDown)
	}

	if ebiten.IsKeyPressed(ebiten.KeyX) {
		j.buttonState = utils.ResetBit(j.buttonState, JoypadA)
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		j.buttonState = utils.ResetBit(j.buttonState, JoypadB)
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		j.buttonState = utils.ResetBit(j.buttonState, JoypadSelect)
	}

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		j.buttonState = utils.ResetBit(j.buttonState, JoypadStart)
	}
}

func (j *Joypad) Read() uint8 {
	if j.joyp&0x10 == 0 {
		return (j.joyp & 0xF0) | (j.buttonState & 0x0F)
	}

	if j.joyp&0x20 == 0 {
		return (j.joyp & 0xF0) | (j.buttonState >> 4)
	}

	return 0xFF
}

func (j *Joypad) Write(val uint8) {
	j.joyp = (j.joyp & 0xCF) | (val & 0x30)
}
