package main

type PPU struct {
	vram Memory
	oam  Memory
}

func NewPPU() *PPU {
	return &PPU{
		vram: NewRAM(0x2000),
		oam:  NewRAM(0xA0),
	}
}

func (p *PPU) ReadVRAM(offset uint16) byte {
	return p.vram.Read(offset)
}

func (p *PPU) WriteVRAM(offset uint16, val byte) {
	p.vram.Write(offset, val)
}

func (p *PPU) ReadOAM(offset uint16) byte {
	return p.oam.Read(offset)
}

func (p *PPU) WriteOAM(offset uint16, val byte) {
	p.oam.Write(offset, val)
}
