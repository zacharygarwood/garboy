package main

const (
	RomBankingMode = 0
	RamBankingMode = 1

	RomBank0End       = 0x3FFF
	RomBankXStart     = 0x4000
	RomBankSize       = 0x4000
	RamBankSize       = 0x2000
	MBC1RamEnableEnd  = 0x1FFF
	MBC1RomBankStart  = 0x2000
	MBC1RomBankEnd    = 0x3FFF
	MBC1RamBankStart  = 0x4000
	MBC1RamBankEnd    = 0x5FFF
	MBC1BankModeStart = 0x6000
	MBC1BankModeEnd   = 0x7FFF
)

type MBC1 struct {
	rom            Memory
	ram            Memory
	romBankCount   int
	ramBankCount   int
	ramSize        int
	currentRomBank int
	currentRamBank int
	ramEnabled     bool
	bankingMode    int
}

func NewMBC1(rom []uint8, header CartridgeHeader) *MBC1 {
	romSize := getRomSize(header.RomSize)
	ramSize := getRamSize(header.RamSize)

	return &MBC1{
		rom:            NewROM(rom),
		ram:            NewRAM(ramSize),
		romBankCount:   romSize / RomBankSize,
		ramBankCount:   ramSize / RamBankSize,
		ramSize:        ramSize,
		currentRomBank: 1,
		currentRamBank: 0,
		ramEnabled:     false,
		bankingMode:    0,
	}
}

func (m *MBC1) Read(address uint16) uint8 {
	switch {
	case address <= RomBank0End:
		return m.rom.Read(address)
	case address >= RomBankXStart && address <= RomBankXEnd:
		bankOffset := uint16(m.currentRomBank) * RomBankSize
		return m.rom.Read(bankOffset + (address - RomBankXStart))
	case address >= RamStart && address <= RamEnd:
		if !m.ramEnabled || m.ramSize == 0 {
			return 0xFF
		}
		bankOffset := uint16(m.currentRamBank) * RamBankSize
		return m.ram.Read(bankOffset + (address - RamStart))
	default:
		panic("Reading from MBC1 at invalid address")
	}
}

func (m *MBC1) Write(address uint16, val uint8) {
	switch {
	case address <= MBC1RamEnableEnd:
		m.ramEnabled = (val & 0x0F) == 0x0A
	case address >= MBC1RomBankStart && address <= MBC1RomBankEnd:
		bank := int(val & 0x1F)
		if bank == 0 {
			bank = 1
		}

		if m.bankingMode == RomBankingMode {
			bank = (m.currentRomBank & 0x60) | bank
		}

		if bank >= m.romBankCount {
			bank = bank % m.romBankCount
		}

		m.currentRomBank = bank
	case address >= MBC1RamBankStart && address <= MBC1RamBankEnd:
		upperBits := int(val & 0x03)

		if m.bankingMode == RomBankingMode {
			m.currentRomBank = (m.currentRomBank & 0x1F) | (upperBits << 5)
			if m.currentRomBank >= m.romBankCount {
				m.currentRomBank = m.currentRomBank % m.romBankCount
			}
		} else {
			if upperBits < m.ramBankCount {
				m.currentRamBank = upperBits
			}
		}
	case address >= MBC1BankModeStart && address <= MBC1BankModeEnd:
		m.bankingMode = int(val & 0x01)
	case address >= RamStart && address <= RamEnd:
		if !m.ramEnabled || m.ramSize == 0 {
			return
		}
		bankOffset := uint16(m.currentRamBank) * RamBankSize
		m.ram.Write(bankOffset+(address-RamStart), val)
	}
}
