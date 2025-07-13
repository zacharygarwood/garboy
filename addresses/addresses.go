package addresses

const (
	// Main memory
	RomBank1     = 0x4000
	Vram         = 0x8000
	ExternalRam  = 0xA000
	Wram         = 0xC000
	EchoRam      = 0xE000
	Oam          = 0xFE00
	NotUsable    = 0xFEA0
	IoRegisters  = 0xFF00
	SerialBuffer = 0xFF01
	Hram         = 0xFF80

	// End of memory addresses
	VramEnd = 0x9FFF
	OamEnd  = 0xFE9F

	// Timer register addresses
	Div  = 0xFF04
	Tima = 0xFF05
	Tma  = 0xFF06
	Tac  = 0xFF07

	// Interrupt addresses
	InterruptFlag   = 0xFF0F
	InterruptEnable = 0xFFFF

	// PPU addresses
	LcdControl  = 0xFF40
	LcdStatus   = 0xFF41
	ScrollY     = 0xFF42
	ScrollX     = 0xFF43
	Ly          = 0xFF44
	Lyc         = 0xFF45
	Dma         = 0xFF46
	BgPalette   = 0xFF47
	ObP0Palette = 0xFF48
	ObP1Palette = 0xFF49
	WindowY     = 0xFF4A
	WindowX     = 0xFF4B

	// MBC0
	RomBankXEnd = 0x7FFF
	RamStart    = 0xA000
	RamEnd      = 0xBFFF

	// MBC1
	RomBank0End   = 0x3FFF
	RomBankXStart = 0x4000

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
