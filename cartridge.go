package main

import "os"

type Cartridge struct {
	rom Memory
	ram Memory
}

func NewCartridge(romPath string, ramSize int) *Cartridge {
	data, err := os.ReadFile(romPath)
	if err != nil {
		panic(err)
	}
	return &Cartridge{
		rom: NewROM(data),
		ram: NewRAM(ramSize),
	}
}

func (c *Cartridge) ReadROM(offset uint16) byte {
	return c.rom.Read(offset)
}

func (c *Cartridge) WriteROM(offset uint16, val byte) {
	c.rom.Write(offset, val)
}

func (c *Cartridge) ReadRAM(offset uint16) byte {
	return c.ram.Read(offset)
}

func (c *Cartridge) WriteRAM(offset uint16, val byte) {
	c.ram.Write(offset, val)
}
