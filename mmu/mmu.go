package mmu

import (
	"fmt"

	"garboy/addresses"
	"garboy/cartridge"
	"garboy/display"
	"garboy/interrupts"
	"garboy/memory"
	"garboy/timer"
)

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

type MmuInterface interface {
	Read(address uint16) uint8
	Write(address uint16, val uint8)
	ReadWord(address uint16) uint16
	WriteWord(address uint16, val uint16)
	SetBootRomEnabled(val bool)
}

type MMU struct {
	cartridge  *cartridge.Cartridge
	ppu        *display.PPU
	timer      *timer.Timer
	joypad     *display.Joypad
	interrupts *interrupts.Interrupts

	wram memory.Memory
	hram memory.Memory
	io   memory.Memory

	bootROM        memory.Memory
	bootROMEnabled bool
}

func NewMMU(cart *cartridge.Cartridge, ppu *display.PPU, timer *timer.Timer, joypad *display.Joypad, interrupts *interrupts.Interrupts) *MMU {
	return &MMU{
		cartridge:      cart,
		ppu:            ppu,
		timer:          timer,
		joypad:         joypad,
		interrupts:     interrupts,
		wram:           memory.NewRAM(0x2000),
		hram:           memory.NewRAM(0x7F),
		io:             memory.NewIORegisters(),
		bootROM:        memory.NewROM(BOOT_ROM[:]),
		bootROMEnabled: true,
	}
}

func (m *MMU) Read(address uint16) byte {
	switch {
	case address <= 0xFF && m.bootROMEnabled:
		return m.bootROM.Read(address)
	case address < addresses.Vram:
		return m.cartridge.Read(address)
	case address < addresses.ExternalRam:
		return m.ppu.Read(address)
	case address < addresses.Wram:
		return m.cartridge.Read(address)
	case address >= addresses.Wram && address < addresses.Oam:
		return m.wram.Read(address & 0x1FFF)
	case address < addresses.NotUsable:
		return m.ppu.Read(address)
	case address < addresses.IoRegisters:
		return 0xFF // Not usable
	case address == addresses.IoRegisters:
		return m.joypad.Read()
	case address >= addresses.Div && address <= addresses.Tac:
		return m.timer.Read(address)
	case address == addresses.InterruptFlag:
		return m.interrupts.IF()
	case address >= addresses.LcdControl && address <= addresses.WindowX:
		return m.ppu.Read(address)
	case address < addresses.Hram:
		return m.io.Read(address - addresses.IoRegisters)
	case address < addresses.InterruptEnable:
		return m.hram.Read(address - addresses.Hram)
	case address == addresses.InterruptEnable:
		return m.interrupts.IE()
	default:
		panic("Should not be reading past 0xFFFF")
	}
}

func (m *MMU) Write(address uint16, val byte) {
	switch {
	case address < addresses.Vram:
		m.cartridge.Write(address, val)
	case address < addresses.ExternalRam:
		m.ppu.Write(address, val)
	case address < addresses.Wram:
		m.cartridge.Write(address, val)
	case address >= addresses.Wram && address < addresses.Oam:
		m.wram.Write(address&0x1FFF, val)
	case address < addresses.NotUsable:
		m.ppu.Write(address, val)
	case address < addresses.IoRegisters:
		return // Not usable
	case address == addresses.IoRegisters:
		m.joypad.Write(val)
	case address == addresses.SerialTransfer && val == 0x81: // NOTE: Used for Blargg's CPU tests. Serial not implemented
		out := m.Read(addresses.SerialBuffer)
		fmt.Printf("%c", out)
	case address >= addresses.Div && address <= addresses.Tac:
		m.timer.Write(address, val)
	case address == addresses.InterruptFlag:
		m.interrupts.Write(address, val)
	case address == addresses.Dma:
		m.DmaTransfer(val)
	case address >= addresses.LcdControl && address <= addresses.WindowX:
		m.ppu.Write(address, val)
	case address == addresses.BootRomControl && m.bootROMEnabled && val != 0:
		m.bootROMEnabled = false
	case address < addresses.Hram:
		m.io.Write(address-addresses.IoRegisters, val)
	case address < addresses.InterruptEnable:
		m.hram.Write(address-addresses.Hram, val)
	case address == addresses.InterruptEnable:
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

func (m *MMU) DmaTransfer(val uint8) {
	src := uint16(val) * 0x100
	dst := uint16(0xFE00)

	segmentSize := uint16(0xA0)
	for i := uint16(0); i < segmentSize; i++ {
		transfer := m.Read(src + i)
		m.Write(dst+i, transfer)
	}
}
