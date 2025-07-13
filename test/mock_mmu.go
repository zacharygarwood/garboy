package main

import (
	"garboy/memory"
)

// SingleStepTests require a flat 64K RAM
type MockMmu struct {
	ram            memory.Memory
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
