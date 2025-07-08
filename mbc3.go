package main

const (
	MBC3RamEnableEnd = 0x1FFF
	MBC3RomBankStart = 0x2000
	MBC3RomBankEnd   = 0x3FFF
	MBC3RamBankStart = 0x4000
	MBC3RamBankEnd   = 0x5FFF
	MBC3LatchStart   = 0x6000
	MBC3LatchEnd     = 0x7FFF
)

type MBC3 struct {
	rom            Memory
	ram            Memory
	romBankCount   int
	ramBankCount   int
	currentRomBank int
	currentRamBank int
	ramEnabled     bool
	rtcEnabled     bool
	rtcRegister    int
	rtcData        [5]uint8
	rtcLatched     [5]uint8
}

func NewMBC3(rom []uint8, header CartridgeHeader) *MBC3 {
	romSize := getRomSize(header.RomSize)
	ramSize := getRamSize(header.RamSize)

	return &MBC3{
		rom:            NewROM(rom),
		ram:            NewRAM(ramSize),
		romBankCount:   romSize / RomBankSize,
		ramBankCount:   ramSize / RamBankSize,
		currentRomBank: 1,
		currentRamBank: 0,
		ramEnabled:     false,
		rtcEnabled:     false,
		rtcRegister:    0,
		rtcData:        [5]uint8{0, 0, 0, 0, 0}, // Seconds, Minutes, Hours, Days low, Days high
		rtcLatched:     [5]uint8{0, 0, 0, 0, 0}, // ^ but latched
	}
}

func (m *MBC3) Read(address uint16) uint8 {
	switch {
	case address <= RomBank0End:
		return m.rom.Read(address)
	case address >= RomBankXStart && address <= RomBankXEnd:
		bankOffset := uint16(m.currentRomBank) * RomBankSize
		return m.rom.Read(bankOffset + (address - RomBankXStart))
	case address >= RamStart && address <= RamEnd:
		if !m.ramEnabled {
			return 0xFF
		}

		if m.currentRamBank >= 0x08 && m.currentRamBank <= 0x0C {
			rtcIndex := m.currentRamBank - 0x08
			return m.rtcLatched[rtcIndex]
		}

		if m.ramBankCount > 0 && m.currentRamBank < m.ramBankCount {
			bankOffset := uint16(m.currentRamBank) * RamBankSize
			return m.ram.Read(bankOffset + (address - RamStart))
		}

		return 0xFF
	default:
		panic("Reading from MBC3 at invalid address")
	}
}

func (m *MBC3) Write(address uint16, val uint8) {
	switch {
	case address <= MBC3RamEnableEnd:
		m.ramEnabled = (val & 0x0F) == 0x0A
	case address >= MBC3RomBankStart && address <= MBC3RomBankEnd:
		bank := int(val & 0x7F)
		if bank == 0 {
			bank = 1
		}
		if bank < m.romBankCount {
			m.currentRomBank = bank
		}
	case address >= MBC3RamBankStart && address <= MBC3RamBankEnd:
		if val <= 0x03 {
			m.currentRamBank = int(val)
		} else if val >= 0x08 && val <= 0x0C {
			m.currentRamBank = int(val)
		}
	case address >= MBC3LatchStart && address <= MBC3LatchEnd:
		if val == 0x01 {
			m.rtcLatched = m.rtcData
		}
	case address >= RamStart && address <= RamEnd:
		if !m.ramEnabled {
			return
		}

		if m.currentRamBank >= 0x08 && m.currentRamBank <= 0x0C {
			rtcIndex := m.currentRamBank - 0x08
			m.rtcData[rtcIndex] = val
			return
		}

		if m.ramBankCount > 0 && m.currentRamBank < m.ramBankCount {
			bankOffset := uint16(m.currentRamBank) * RamBankSize
			m.ram.Write(bankOffset+(address-RamStart), val)
		}
	}
}
