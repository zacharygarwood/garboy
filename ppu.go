package main

import "fmt"

const (
	// PPU constants
	ScreenWidth    = 160
	ScreenHeight   = 144
	TileSize       = 8
	TileMapSize    = 32
	OamSize        = 160
	MaxSprites     = 40
	SpritesPerLine = 10

	// LCD Control Register bits
	LcdEnable        = 7
	WindowTileMap    = 6
	WindowEnable     = 5
	BgWindowTileData = 4
	BgTileMap        = 3
	SpriteSize       = 2
	SpriteEnable     = 1
	BgWindowEnable   = 0

	// LCD Status Register bits
	LycInterruptBit    = 6
	OamInterruptBit    = 5
	VBlankInterruptBit = 4
	HBlankInterruptBit = 3
	LycFlagBit         = 2
	ModeFlag           = 0x03

	// PPU Modes
	HBlankMode = 0
	VBlankMode = 1
	OamMode    = 2
	VramMode   = 3

	// PPU timing in T-cycles
	OamScanCycles  = 80
	VramScanCycles = 172
	HBlankCycles   = 204
	ScanlineCycles = 456
	VBlankLines    = 10

	// Sprite flags
	SpritePriority = 7
	SpriteYFlip    = 6
	SpriteXFlip    = 5
	SpritePalette  = 4
)

type Color uint8

type Sprite struct {
	y         uint8
	x         uint8
	tileIndex uint8
	flags     uint8
}

type PPU struct {
	vram Memory
	oam  Memory

	lcdc uint8
	stat uint8
	scy  uint8
	scx  uint8
	ly   uint8
	lyc  uint8
	bgp  uint8
	obp0 uint8
	obp1 uint8
	wy   uint8
	wx   uint8

	mode        uint8
	cycles      int
	framebuffer [ScreenHeight][ScreenWidth]Color

	interrupts *Interrupts
}

func NewPPU(interrupts *Interrupts) *PPU {
	return &PPU{
		vram:       NewRAM(0x2000),
		oam:        NewRAM(0xA0),
		lcdc:       0x91,
		stat:       0x85,
		bgp:        0xFC,
		obp0:       0xFF,
		obp1:       0xFF,
		interrupts: interrupts,
	}
}

func (p *PPU) Step() {
	if !p.isLcdEnabled() {
		return
	}

	p.cycles += 4

	switch p.mode {
	case OamMode:
		p.handleOamMode()
	case VramMode:
		p.handleVramMode()
	case HBlankMode:
		p.handleHBlankMode()
	case VBlankMode:
		p.handleVBlankMode()
	}
}

func (p *PPU) handleOamMode() {
	if p.cycles >= OamScanCycles {
		p.enterMode(VramMode)
	}
}

func (p *PPU) handleVramMode() {
	if p.cycles >= VramScanCycles {
		p.enterMode(HBlankMode)
		p.renderScanline()
		p.checkHBlankInterrupt()
	}
}

func (p *PPU) handleHBlankMode() {
	if p.cycles >= HBlankCycles {
		p.moveToNextScanline()

		if p.ly == ScreenHeight {
			p.enterMode(VBlankMode)
			p.checkVBlankInterrupt()
		} else {
			p.enterMode(OamMode)
			p.checkOamInterrupt()
		}
	}
}

func (p *PPU) handleVBlankMode() {
	if p.cycles >= ScanlineCycles {
		p.moveToNextScanline()

		if p.ly > 153 {
			p.ly = 0
			p.enterMode(OamMode)
			p.updateLyc()
			p.checkOamInterrupt()
		}
	}
}

func (p *PPU) enterMode(mode uint8) {
	p.mode = mode
	p.stat = (p.stat & ^uint8(ModeFlag)) | mode
	p.cycles = 0
}

func (p *PPU) moveToNextScanline() {
	p.ly++
	p.updateLyc()
}

func (p *PPU) updateLyc() {
	if p.ly == p.lyc {
		p.stat |= LycFlagBit
		if IsBitSet(p.stat, LycInterruptBit) {
			p.interrupts.Request(LcdInterrupt)
		}
	} else {
		p.stat &= ^uint8(1 << LycFlagBit)
	}
}

func (p *PPU) checkVBlankInterrupt() {
	p.interrupts.Request(VBlankInterrupt)

	if IsBitSet(p.stat, VBlankInterruptBit) {
		p.interrupts.Request(LcdInterrupt)
	}
}

func (p *PPU) checkHBlankInterrupt() {
	if IsBitSet(p.stat, HBlankInterruptBit) {
		p.interrupts.Request(LcdInterrupt)
	}
}

func (p *PPU) checkOamInterrupt() {
	if IsBitSet(p.stat, OamInterruptBit) {
		p.interrupts.Request(LcdInterrupt)
	}
}

func (p *PPU) renderScanline() {
	if p.ly >= ScreenHeight {
		return
	}

	p.clearScanline()

	if IsBitSet(p.lcdc, BgWindowEnable) {
		p.renderBackground()
	}

	if IsBitSet(p.lcdc, WindowEnable) && p.shouldRenderWindow() {
		p.renderWindow()
	}

	if IsBitSet(p.lcdc, SpriteEnable) {
		p.renderSprites()
	}
}

func (p *PPU) clearScanline() {
	for x := 0; x < ScreenWidth; x++ {
		p.framebuffer[p.ly][x] = 0
	}
}

func (p *PPU) shouldRenderWindow() bool {
	return p.ly >= p.wy && p.wx <= 166
}

func (p *PPU) renderBackground() {
	tileMapBase := p.getBackgroundTileMapBase()

	scrolledY := int(p.ly) + int(p.scy)
	tileRow := (scrolledY / TileSize) % TileMapSize
	pixelRow := scrolledY % TileSize

	for screenX := 0; screenX < ScreenWidth; screenX++ {
		scrolledX := screenX + int(p.scx)
		tileCol := (scrolledX / TileSize) % TileMapSize
		pixelCol := scrolledX % TileSize

		tileMapAddress := tileMapBase + uint16(tileRow*TileMapSize+tileCol)
		tileIndex := p.readVram(tileMapAddress)

		color := p.getTilePixel(tileIndex, pixelCol, pixelRow)
		p.framebuffer[p.ly][screenX] = p.applyBackgroundPalette(color)
	}
}

func (p *PPU) getBackgroundTileMapBase() uint16 {
	if IsBitSet(p.lcdc, BgTileMap) {
		return 0x9C00
	}
	return 0x9800
}

func (p *PPU) getWindowTileMapBase() uint16 {
	if IsBitSet(p.lcdc, WindowTileMap) {
		return 0x9C00
	}
	return 0x9800
}

func (p *PPU) renderWindow() {
	windowY := int(p.ly) - int(p.wy)
	if windowY < 0 {
		return
	}

	tileMapBase := p.getWindowTileMapBase()

	tileRow := windowY / TileSize
	pixelRow := windowY % TileSize

	windowStartX := int(p.wx) - 7
	if windowStartX < 0 {
		windowStartX = 0
	}

	for screenX := windowStartX; screenX < ScreenWidth; screenX++ {
		windowX := screenX - windowStartX
		tileCol := windowX / TileSize
		pixelCol := windowX % TileSize

		tileMapAddress := tileMapBase + uint16(tileRow*TileMapSize+tileCol)
		tileIndex := p.readVram(tileMapAddress)

		color := p.getTilePixel(tileIndex, pixelCol, pixelRow)
		p.framebuffer[p.ly][screenX] = p.applyBackgroundPalette(color)
	}
}

func (p *PPU) applyBackgroundPalette(color Color) Color {
	shift := color * 2
	return Color((p.bgp >> shift) & 3)
}

func (p *PPU) renderSprites() {
	spriteHeight := p.getSpriteHeight()
	spritesRendered := 0

	for spriteIndex := MaxSprites - 1; spriteIndex >= 0 && spritesRendered < SpritesPerLine; spriteIndex-- {
		sprite := p.getSprite(spriteIndex)

		if p.isSpriteOnCurrentScanline(sprite, spriteHeight) {
			p.renderSprite(sprite, spriteHeight)
			spritesRendered++
		}
	}
}

func (p *PPU) getSpriteHeight() int {
	if IsBitSet(p.lcdc, SpriteSize) {
		return 16
	}
	return 8
}

func (p *PPU) getSprite(index int) Sprite {
	base := uint16(OamAddress) + uint16(index*4)
	return Sprite{
		y:         p.readOam(base),
		x:         p.readOam(base + 1),
		tileIndex: p.readOam(base + 2),
		flags:     p.readOam(base + 3),
	}
}

func (p *PPU) isSpriteOnCurrentScanline(sprite Sprite, spriteHeight int) bool {
	if sprite.y == 0 || sprite.y >= 160 {
		return false
	}

	spriteY := int(sprite.y) - 16
	return int(p.ly) >= spriteY && int(p.ly) < spriteY+spriteHeight
}

func (p *PPU) renderSprite(sprite Sprite, spriteHeight int) {
	spriteY := int(sprite.y) - 16
	spriteLine := int(p.ly) - spriteY

	if IsBitSet(sprite.flags, SpriteYFlip) {
		spriteLine = spriteHeight - 1 - spriteLine
	}

	tileIndex := sprite.tileIndex
	if spriteHeight == 16 {
		tileIndex &= 0xFE

		if spriteLine >= 8 {
			tileIndex |= 0x01
			spriteLine -= 8
		}
	}

	for pixelX := 0; pixelX < 8; pixelX++ {
		screenX := int(sprite.x) - 8 + pixelX

		if screenX < 0 || screenX >= ScreenWidth {
			continue
		}

		spritePixelX := pixelX
		if IsBitSet(sprite.flags, SpriteXFlip) {
			spritePixelX = 7 - pixelX
		}

		color := p.getSpriteTilePixel(tileIndex, spritePixelX, spriteLine)

		if color.isTransparent() {
			continue
		}

		if IsBitSet(sprite.flags, SpritePriority) && p.framebuffer[p.ly][screenX] != 0 {
			continue
		}

		paletteColor := p.applySpritePalette(sprite, color)
		p.framebuffer[p.ly][screenX] = paletteColor
	}
}

func (c Color) isTransparent() bool {
	return c == 0
}

func (p *PPU) applySpritePalette(sprite Sprite, color Color) Color {
	palette := p.obp0
	if IsBitSet(sprite.flags, SpritePalette) {
		palette = p.obp1
	}

	shift := color * 2
	return Color((palette >> shift) & 3)
}

func (p *PPU) getTilePixel(tileIndex uint8, pixelX int, pixelY int) Color {
	tileDataAddress := p.getTileDataAddress(tileIndex)
	return p.getPixelFromTileData(tileDataAddress, pixelX, pixelY)
}

func (p *PPU) getSpriteTilePixel(tileIndex uint8, pixelX int, pixelY int) Color {
	tileDataAddress := uint16(0x8000) + uint16(tileIndex)*16
	return p.getPixelFromTileData(tileDataAddress, pixelX, pixelY)
}

func (p *PPU) getPixelFromTileData(tileDataAdddress uint16, pixelX int, pixelY int) Color {
	lineOffset := pixelY * 2
	lowByte := p.readVram(tileDataAdddress + uint16(lineOffset))
	highByte := p.readVram(tileDataAdddress + uint16(lineOffset) + 1)

	bitPosition := 7 - pixelX
	lowBit := (lowByte >> bitPosition) & 1
	highBit := (highByte >> bitPosition) & 1

	return Color((highBit << 1) | lowBit)
}

func (p *PPU) getTileDataAddress(tileIndex uint8) uint16 {
	if IsBitSet(p.lcdc, BgWindowTileData) {
		return uint16(0x8000) + uint16(tileIndex)*16
	} else {
		signedIndex := int8(tileIndex)
		return uint16(int32(0x9000) + int32(signedIndex)*16)
	}
}

func (p *PPU) isLcdEnabled() bool {
	return IsBitSet(p.lcdc, LcdEnable)
}

func (p *PPU) readVram(address uint16) uint8 {
	if address >= VramAddress && address <= VramEndAddress {
		return p.vram.Read(address - VramAddress)
	}
	return 0xFF
}

func (p *PPU) writeVram(address uint16, val uint8) {
	if address >= VramAddress && address <= VramEndAddress {
		p.vram.Write(address-VramAddress, val)
	}
}

func (p *PPU) readOam(address uint16) uint8 {
	//fmt.Printf("[DEBUG] Reading from OAM. address: %x\n", address)
	return p.oam.Read(address - OamAddress)
}

func (p *PPU) writeOam(address uint16, val uint8) {
	//fmt.Printf("[DEBUG] Writing to OAM. address: %x, val: %x\n", address, val)
	p.oam.Write(address-OamAddress, val)
}

func (p *PPU) Read(address uint16) uint8 {
	//fmt.Printf("[DEBUG] In Read(). address: %x\n", address)
	switch address {
	case LcdControlAddress:
		return p.lcdc
	case LcdStatusAddress:
		return p.stat | 0x80 // 7th bit always 1
	case ScrollYAddress:
		return p.scy
	case ScrollXAddress:
		return p.scx
	case LyAddress:
		return p.ly
	case LycAddress:
		return p.lyc
	case WindowYAddress:
		return p.wy
	case WindowXAddress:
		return p.wx
	case BgPaletteAddress:
		return p.bgp
	case ObP0PaletteAddress:
		return p.obp0
	case ObP1PaletteAddress:
		return p.obp1
	default:
		if address >= VramAddress && address <= VramEndAddress {
			return p.readVram(address)
		}
		if address >= OamAddress && address <= OamEndAddress {
			return p.readOam(address)
		}
		return 0xFF
	}
}

func (p *PPU) Write(address uint16, val uint8) {
	//fmt.Printf("[DEBUG] In Write(), address: %x, val: %x\n", address, val)
	switch address {
	case LcdControlAddress:
		prevEnabled := p.isLcdEnabled()
		p.lcdc = val

		if prevEnabled && !p.isLcdEnabled() {
			p.ly = 0
			p.enterMode(HBlankMode)
		}
	case LcdStatusAddress:
		p.stat = (p.stat & 0x87) | (val & 0x78) // Only bits 6-3 are writable
	case ScrollYAddress:
		p.scy = val
	case ScrollXAddress:
		p.scx = val
	case LycAddress:
		p.lyc = val
		p.updateLyc()
	case WindowYAddress:
		p.wy = val
	case WindowXAddress:
		p.wx = val
	case BgPaletteAddress:
		p.bgp = val
	case ObP0PaletteAddress:
		p.obp0 = val
	case ObP1PaletteAddress:
		p.obp1 = val
	default:
		if address >= VramAddress && address <= VramEndAddress {
			p.writeVram(address, val)
			return
		}
		if address >= OamAddress && address <= OamEndAddress {
			p.writeOam(address, val)
			return
		}
	}
}

func (p *PPU) GetFrameBuffer() *[ScreenHeight][ScreenWidth]Color {
	return &p.framebuffer
}

func (p *PPU) Reset() {
	p.lcdc = 0x91
	p.stat = 0x85
	p.scy = 0
	p.scx = 0
	p.ly = 0
	p.lyc = 0
	p.bgp = 0xFC
	p.obp0 = 0xFF
	p.obp1 = 0xFF
	p.wy = 0
	p.wx = 0
	p.mode = 0
	p.cycles = 0

	p.vram = NewRAM(0x2000)
	p.oam = NewRAM(0xA0)

	for y := range p.framebuffer {
		for x := range p.framebuffer[y] {
			p.framebuffer[y][x] = 0
		}
	}
}

func (p *PPU) PrintState() {
	fmt.Printf("PPU: LCDC=%x LY:%x LYC:%x MODE:%x OBP0:%x OBP1:%x SCX:%x SCY:%x STAT:%x WX:%x WY:%x\n",
		p.lcdc, p.ly, p.lyc, p.mode, p.obp0, p.obp1, p.scx, p.scy, p.stat, p.wx, p.wy)
}
