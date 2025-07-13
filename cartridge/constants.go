package cartridge

const (
	// Main memory
	RomBank1Address    = 0x4000
	VramAddress        = 0x8000
	ExternalRamAddress = 0xA000
	WramAddress        = 0xC000

	// MBC0
	RomBankXEnd = 0x7FFF
	RamStart    = 0xA000
	RamEnd      = 0xBFFF

	// MBC1
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

	// MBC3
	MBC3RamEnableEnd = 0x1FFF
	MBC3RomBankStart = 0x2000
	MBC3RomBankEnd   = 0x3FFF
	MBC3RamBankStart = 0x4000
	MBC3RamBankEnd   = 0x5FFF
	MBC3LatchStart   = 0x6000
	MBC3LatchEnd     = 0x7FFF
)
