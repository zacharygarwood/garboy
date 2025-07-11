package main

import "fmt"

var BOOT_ROM = [256]byte{
	0x31, 0xFE, 0xFF, 0xAF, 0x21, 0xFF, 0x9F, 0x32, 0xCB, 0x7C, 0x20, 0xFB, 0x21, 0x26, 0xFF, 0x0E,
	0x11, 0x3E, 0x80, 0x32, 0xE2, 0x0C, 0x3E, 0xF3, 0xE2, 0x32, 0x3E, 0x77, 0x77, 0x3E, 0xFC, 0xE0,
	0x47, 0x11, 0x04, 0x01, 0x21, 0x10, 0x80, 0x1A, 0xCD, 0x95, 0x00, 0xCD, 0x96, 0x00, 0x13, 0x7B,
	0xFE, 0x34, 0x20, 0xF3, 0x11, 0xD8, 0x00, 0x06, 0x08, 0x1A, 0x13, 0x22, 0x23, 0x05, 0x20, 0xF9,
	0x3E, 0x19, 0xEA, 0x10, 0x99, 0x21, 0x2F, 0x99, 0x0E, 0x0C, 0x3D, 0x28, 0x08, 0x32, 0x0D, 0x20,
	0xF9, 0x2E, 0x0F, 0x18, 0xF3, 0x67, 0x3E, 0x64, 0x57, 0xE0, 0x42, 0x3E, 0x91, 0xE0, 0x40, 0x04,
	0x1E, 0x02, 0x0E, 0x0C, 0xF0, 0x44, 0xFE, 0x90, 0x20, 0xFA, 0x0D, 0x20, 0xF7, 0x1D, 0x20, 0xF2,
	0x0E, 0x13, 0x24, 0x7C, 0x1E, 0x83, 0xFE, 0x62, 0x28, 0x06, 0x1E, 0xC1, 0xFE, 0x64, 0x20, 0x06,
	0x7B, 0xE2, 0x0C, 0x3E, 0x87, 0xE2, 0xF0, 0x42, 0x90, 0xE0, 0x42, 0x15, 0x20, 0xD2, 0x05, 0x20,
	0x4F, 0x16, 0x20, 0x18, 0xCB, 0x4F, 0x06, 0x04, 0xC5, 0xCB, 0x11, 0x17, 0xC1, 0xCB, 0x11, 0x17,
	0x05, 0x20, 0xF5, 0x22, 0x23, 0x22, 0x23, 0xC9, 0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
	0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
	0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
	0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E, 0x3C, 0x42, 0xB9, 0xA5, 0xB9, 0xA5, 0x42, 0x3C,
	0x21, 0x04, 0x01, 0x11, 0xA8, 0x00, 0x1A, 0x13, 0xBE, 0x00, 0x00, 0x23, 0x7D, 0xFE, 0x34, 0x20,
	0xF5, 0x06, 0x19, 0x78, 0x86, 0x23, 0x05, 0x20, 0xFB, 0x86, 0x00, 0x00, 0x3E, 0x01, 0xE0, 0x50,
}

const (
	// Start of memory addresses
	RomBank1Address     = 0x4000
	VramAddress         = 0x8000
	ExternalRamAddress  = 0xA000
	WramAddress         = 0xC000
	EchoRamAddress      = 0xE000
	OamAddress          = 0xFE00
	NotUsableAddress    = 0xFEA0
	IoRegistersAddress  = 0xFF00
	SerialBufferAddress = 0xFF01
	HramAddress         = 0xFF80

	// End of memory addresses
	VramEndAddress = 0x9FFF
	OamEndAddress  = 0xFE9F

	// Timer register addresses
	DivAddress  = 0xFF04
	TimaAddress = 0xFF05
	TmaAddress  = 0xFF06
	TacAddress  = 0xFF07

	// Interrupt addresses
	InterruptFlagAddress   = 0xFF0F
	InterruptEnableAddress = 0xFFFF

	// PPU addresses
	LcdControlAddress  = 0xFF40
	LcdStatusAddress   = 0xFF41
	ScrollYAddress     = 0xFF42
	ScrollXAddress     = 0xFF43
	LyAddress          = 0xFF44
	LycAddress         = 0xFF45
	BgPaletteAddress   = 0xFF47
	ObP0PaletteAddress = 0xFF48
	ObP1PaletteAddress = 0xFF49
	WindowYAddress     = 0xFF4A
	WindowXAddress     = 0xFF4B
)

type MMU struct {
	cartridge  *Cartridge
	ppu        *PPU
	timer      *Timer
	interrupts *Interrupts

	wram Memory
	hram Memory
	io   Memory

	bootROM        Memory
	bootROMEnabled bool
}

func NewMMU(cart *Cartridge, ppu *PPU, timer *Timer, interrupts *Interrupts) *MMU {
	return &MMU{
		cartridge:      cart,
		ppu:            ppu,
		timer:          timer,
		interrupts:     interrupts,
		wram:           NewRAM(0x2000),
		hram:           NewRAM(0x7F),
		io:             NewIORegisters(),
		bootROM:        NewROM(BOOT_ROM[:]),
		bootROMEnabled: true,
	}
}

func (m *MMU) Read(address uint16) byte {
	switch {
	case address <= 0xFF && m.bootROMEnabled:
		return m.bootROM.Read(address)
	case address < VramAddress:
		return m.cartridge.Read(address)
	case address < ExternalRamAddress:
		return m.ppu.Read(address)
	case address < WramAddress:
		return m.cartridge.Read(address)
	case address >= WramAddress && address < OamAddress:
		return m.wram.Read(address & 0x1FFF)
	case address < NotUsableAddress:
		return m.ppu.Read(address)
	case address < IoRegistersAddress:
		return 0xFF // Not usable
	case address >= DivAddress && address <= TacAddress:
		return m.timer.Read(address)
	case address == InterruptFlagAddress:
		return m.interrupts.IF()
	case address >= LcdControlAddress && address <= WindowXAddress:
		return m.ppu.Read(address)
	case address < HramAddress:
		return m.io.Read(address - IoRegistersAddress)
	case address < InterruptEnableAddress:
		return m.hram.Read(address - HramAddress)
	case address == InterruptEnableAddress:
		return m.interrupts.IE()
	default:
		panic("Should not be reading past 0xFFFF")
	}
}

func (m *MMU) Write(address uint16, val byte) {
	switch {
	case address < VramAddress:
		m.cartridge.Write(address, val)
	case address < ExternalRamAddress:
		m.ppu.Write(address, val)
	case address < WramAddress:
		m.cartridge.Write(address, val)
	case address >= WramAddress && address < OamAddress:
		m.wram.Write(address&0x1FFF, val)
	case address < NotUsableAddress:
		m.ppu.Write(address, val)
	case address < IoRegistersAddress:
		return // Not usable
	case address == 0xFF02 && val == 0x81: // FIXME: Used for Blargg's CPU tests
		out := m.Read(0xFF01)
		fmt.Printf("%c", out)
	case address >= DivAddress && address <= TacAddress:
		m.timer.Write(address, val)
	case address == InterruptFlagAddress:
		m.interrupts.Write(address, val)
	case address >= LcdControlAddress && address <= WindowXAddress:
		m.ppu.Write(address, val)
	case address == 0xFF50 && m.bootROMEnabled && val != 0:
		m.bootROMEnabled = false
	case address < HramAddress:
		m.io.Write(address-IoRegistersAddress, val)
	case address < InterruptEnableAddress:
		m.hram.Write(address-HramAddress, val)
	case address == InterruptEnableAddress:
		m.interrupts.Write(address, val)
	default:
		panic("Should not be reading past 0xFFFF")
	}
}

func (m *MMU) ReadWord(address uint16) uint16 {
	lo := m.Read(address)
	hi := m.Read(address + 1)

	return (uint16(hi) << 8) + uint16(lo)
}

func (m *MMU) WriteWord(address uint16, val uint16) {
	hi := uint8((val >> 8) & 0xFF)
	lo := uint8(val & 0xFF)

	m.Write(address, lo)
	m.Write(address+1, hi)
}

func (m *MMU) SetBootRomEnabled(val bool) {
	m.bootROMEnabled = val
}

// TODO: Move to separate module. These are used for the cpu single step instr tests
type MmuInterface interface {
	Read(address uint16) uint8
	Write(address uint16, val uint8)
	ReadWord(address uint16) uint16
	WriteWord(address uint16, val uint16)

	SetBootRomEnabled(val bool)
}

// SingleStepTests require a flat 64K RAM
type MockMmu struct {
	ram            Memory
	bootROMEnabled bool
}

func (m *MockMmu) Read(address uint16) uint8 {
	return m.ram.Read(address)
}

func (m *MockMmu) Write(address uint16, val uint8) {
	m.ram.Write(address, val)
}

func (m *MockMmu) ReadWord(address uint16) uint16 {
	lo := m.Read(address)
	hi := m.Read(address + 1)

	return (uint16(hi) << 8) + uint16(lo)
}

func (m *MockMmu) WriteWord(address uint16, val uint16) {
	hi := uint8((val >> 8) & 0xFF)
	lo := uint8(val & 0xFF)

	m.Write(address, lo)
	m.Write(address+1, hi)
}

func (m *MockMmu) SetBootRomEnabled(val bool) {
	m.bootROMEnabled = val
}
