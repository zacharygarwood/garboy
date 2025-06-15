package main

type MMU struct {
	cartridge *Cartridge
	ppu       *PPU
	wram      Memory
	hram      Memory
	io        Memory
	ie        Memory
}

func NewMMU(cart *Cartridge, ppu *PPU) *MMU {
	return &MMU{
		cartridge: cart,
		ppu:       ppu,
		wram:      NewRAM(0x2000),
		hram:      NewRAM(0x7F),
		io:        NewIORegisters(),
		ie:        NewInterruptRegister(),
	}
}

func (m *MMU) Read(address uint16) byte {
	switch {
	case address < 0x4000:
		return m.cartridge.ReadROM(address)
	case address < 0x8000:
		return m.cartridge.ReadROM(address - 0x4000)
	case address < 0xA000:
		return m.ppu.ReadVRAM(address - 0x8000)
	case address < 0xC000:
		return m.cartridge.ReadRAM(address - 0xA000)
	case address < 0xE000:
		return m.wram.Read(address - 0xC000)
	case address < 0xFE00:
		return m.wram.Read(address - 0xE000) // Echo RAM
	case address < 0xFEA0:
		return m.ppu.ReadOAM(address - 0xFE00)
	case address < 0xFF00:
		return 0xFF // Not usable
	case address < 0xFF80:
		return m.io.Read(address - 0xFF00)
	case address < 0xFFFF:
		return m.hram.Read(address - 0xFF80)
	case address == 0xFFFF:
		return m.ie.Read(0)
	default:
		panic("Should not be reading past 0xFFFF")
	}
}

func (m *MMU) Write(address uint16, val byte) {
	switch {
	case address < 0x4000:
		m.cartridge.WriteROM(address, val)
	case address < 0x8000:
		m.cartridge.WriteROM(address-0x4000, val)
	case address < 0xA000:
		m.ppu.WriteVRAM(address-0x8000, val)
	case address < 0xC000:
		m.cartridge.WriteRAM(address-0xA000, val)
	case address < 0xE000:
		m.wram.Write(address-0xC000, val)
	case address < 0xFE00:
		m.wram.Write(address-0xE000, val) // Echo RAM
	case address < 0xFEA0:
		m.ppu.WriteOAM(address-0xFE00, val)
	case address < 0xFF00:
		return // Not usable
	case address < 0xFF80:
		m.io.Write(address-0xFF00, val)
	case address < 0xFFFF:
		m.hram.Write(address-0xFF80, val)
	case address == 0xFFFF:
		m.ie.Write(0, val)
	default:
		panic("Should not be reading past 0xFFFF")
	}
}
