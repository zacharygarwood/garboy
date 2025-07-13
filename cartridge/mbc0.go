package cartridge

import (
	"garboy/addresses"
	"garboy/memory"
)

type MBC0 struct {
	rom     memory.Memory
	ram     memory.Memory
	ramSize int
}

func NewMBC0(romData []uint8, header CartridgeHeader) *MBC0 {
	ramSize := getRamSize(header.RamSize)
	rom := make([]byte, 0x8000)
	copy(rom, romData)

	return &MBC0{
		rom:     memory.NewROM(rom),
		ram:     memory.NewRAM(ramSize),
		ramSize: ramSize,
	}
}

func (m *MBC0) Read(address uint16) uint8 {
	switch {
	case address <= addresses.RomBankXEnd:
		return m.rom.Read(address)
	case address >= addresses.RamStart && address <= addresses.RamEnd:
		if m.ramSize > 0 {
			return m.ram.Read(address - addresses.RamStart)
		}
		return 0xFF
	default:
		panic("Reading from MBC0 at an invalid address")
	}
}

func (m *MBC0) Write(address uint16, val uint8) {
	switch {
	case address <= addresses.RomBankXEnd:
		return
	case address >= addresses.RamStart && address <= addresses.RamEnd:
		if m.ramSize > 0 {
			m.ram.Write(address-addresses.RamStart, val)
		}
	}
}
