package main

const (
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
	rom        []byte
	ram        []byte
	romBank    byte
	ramBank    byte
	ramEnabled bool
	bankMode   byte
	hasRam     bool
}

func NewMBC1(data []byte, header CartridgeHeader) *MBC1 {
	mbc := &MBC1{
		rom:        data,
		romBank:    1,
		ramBank:    0,
		ramEnabled: false,
		bankMode:   0,
		hasRam:     header.CartType == 0x02 || header.CartType == 0x03,
	}

	if mbc.hasRam {
		ramSize := getRamSize(header.RamSize)
		mbc.ram = make([]byte, ramSize)
	}

	return mbc
}

func (m *MBC1) Read(address uint16) byte {
	switch {
	case address < MBC1RamBankStart:
		bankOffset := 0
		if m.bankMode == 1 {
			bankOffset = int(m.ramBank) << 5
		}
		return m.rom[bankOffset*RomBankSize+int(address)]
	case address < VramAddress:
		actualBank := m.romBank
		if m.bankMode == 0 {
			actualBank |= m.ramBank << 5
		}
		bankOffset := int(actualBank) * RomBankSize
		romAddress := bankOffset + int(address-MBC1RamBankStart)
		if romAddress < len(m.rom) {
			return m.rom[romAddress]
		}
		return 0xFF
	case address >= ExternalRamAddress && address < WramAddress:
		if !m.ramEnabled || !m.hasRam {
			return 0xFF
		}

		actualRamBank := byte(0)
		if m.bankMode == 1 {
			actualRamBank = m.ramBank
		}

		ramAddress := int(actualRamBank)*RamBankSize + int(address-ExternalRamAddress)
		if ramAddress < len(m.ram) {
			return m.ram[ramAddress]
		}
		return 0xFF
	default:
		return 0xFF
	}
}

func (m *MBC1) Write(address uint16, val byte) {
	switch {
	case address < MBC1RomBankStart:
		m.ramEnabled = (val & 0x0F) == 0x0A
	case address < MBC1RamBankStart:
		bank := val & 0x1F
		if bank == 0 {
			bank = 1
		}
		m.romBank = bank
	case address < MBC1BankModeStart:
		m.ramBank = val & 0x03
	case address < VramAddress:
		m.bankMode = val & 0x01
	case address >= ExternalRamAddress && address < WramAddress:
		if !m.ramEnabled || !m.hasRam {
			return
		}

		actualRamBank := byte(0)
		if m.bankMode == 1 {
			actualRamBank = m.ramBank
		}

		ramAddress := int(actualRamBank)*RamBankSize + int(address-0xA000)
		if ramAddress < len(m.ram) {
			m.ram[ramAddress] = val
		}
	}
}
