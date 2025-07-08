package main

const (
	RomBankXEnd = 0x7FFF
	RamStart    = 0xA000
	RamEnd      = 0xBFFF
)

type MBC0 struct {
	rom     Memory
	ram     Memory
	ramSize int
}

func NewMBC0(romData []uint8, header CartridgeHeader) *MBC0 {
	ramSize := getRamSize(header.RamSize)
	rom := make([]byte, 0x8000)
	copy(rom, romData)

	return &MBC0{
		rom:     NewROM(rom),
		ram:     NewRAM(ramSize),
		ramSize: ramSize,
	}
}

func (m *MBC0) Read(address uint16) uint8 {
	switch {
	case address <= RomBankXEnd:
		return m.rom.Read(address)
	case address >= RamStart && address <= RamEnd:
		if m.ramSize > 0 {
			return m.ram.Read(address - RamStart)
		}
		return 0xFF
	default:
		panic("Reading from MBC0 at an invalid address")
	}
}

func (m *MBC0) Write(address uint16, val uint8) {
	switch {
	case address <= RomBankXEnd:
		return
	case address >= RamStart && address <= RamEnd:
		if m.ramSize > 0 {
			m.ram.Write(address-RamStart, val)
		}
	}
}
