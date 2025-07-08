package main

import "os"

type Cartridge struct {
	mbc    MBC
	header CartridgeHeader
}

type CartridgeHeader struct {
	CartType uint8
	RomSize  uint8
	RamSize  uint8
}

func NewCartridge(romPath string) *Cartridge {
	data, err := os.ReadFile(romPath)
	if err != nil {
		panic(err)
	}

	header := CartridgeHeader{
		CartType: data[0x147],
		RomSize:  data[0x148],
		RamSize:  data[0x149],
	}

	return &Cartridge{
		mbc:    NewMBC(data, header),
		header: header,
	}
}

func (c *Cartridge) ReadROM(address uint16) byte {
	return c.mbc.Read(address)
}

func (c *Cartridge) WriteROM(address uint16, val byte) {
	c.mbc.Write(address, val)
}

func (c *Cartridge) ReadRAM(address uint16) byte {
	return c.mbc.Read(address)
}

func (c *Cartridge) WriteRAM(address uint16, val byte) {
	c.mbc.Write(address, val)
}
