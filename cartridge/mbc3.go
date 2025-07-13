package cartridge

import (
	"time"

	"garboy/addresses"
)

type MBC3 struct {
	rom          []byte
	ram          []byte
	romBank      byte
	ramBank      byte
	ramEnabled   bool
	hasRam       bool
	hasTimer     bool
	rtcSelect    byte
	rtcRegs      [5]byte // Seconds, Minutes, Hours, Days Low, Days High
	rtcLatch     [5]byte
	rtcLatchData byte
	rtcBaseTime  time.Time
	rtcHalt      bool
}

func NewMBC3(data []byte, header CartridgeHeader) *MBC3 {
	mbc := &MBC3{
		rom:         data,
		romBank:     1,
		ramBank:     0,
		ramEnabled:  false,
		hasRam:      header.CartType == 0x10 || header.CartType == 0x12 || header.CartType == 0x13,
		hasTimer:    header.CartType == 0x0F || header.CartType == 0x10,
		rtcBaseTime: time.Now(),
	}

	if mbc.hasRam {
		ramSize := getRamSize(header.RamSize)
		mbc.ram = make([]byte, ramSize)
	}

	return mbc
}

func (m *MBC3) Read(address uint16) byte {
	switch {
	case address < addresses.RomBank1:
		return m.rom[address]
	case address < addresses.Vram:
		bankOffset := int(m.romBank) * RomBankSize
		romAddress := bankOffset + int(address-addresses.RomBank1)
		if romAddress < len(m.rom) {
			return m.rom[romAddress]
		}
		return 0xFF
	case address >= addresses.ExternalRam && address < addresses.Wram:
		if !m.ramEnabled {
			return 0xFF
		}

		if m.ramBank <= 0x03 {
			if !m.hasRam {
				return 0xFF
			}
			ramAddress := int(m.ramBank)*RamBankSize + int(address-addresses.ExternalRam)
			if ramAddress < len(m.ram) {
				return m.ram[ramAddress]
			}
			return 0xFF
		} else if m.ramBank >= 0x08 && m.ramBank <= 0x0C && m.hasTimer {
			return m.rtcLatch[m.ramBank-0x08]
		}
		return 0xFF
	default:
		return 0xFF
	}
}

func (m *MBC3) Write(address uint16, val byte) {
	switch {
	case address < addresses.MBC3RomBankStart:
		m.ramEnabled = (val & 0x0F) == 0x0A
	case address < addresses.MBC3RamBankStart:
		bank := val & 0x7F
		if bank == 0 {
			bank = 1
		}
		m.romBank = bank
	case address < addresses.MBC3LatchStart:
		m.ramBank = val
	case address < addresses.Vram:
		if m.hasTimer {
			if m.rtcLatchData == 0x00 && val == 0x01 {
				m.latchRTC()
			}
			m.rtcLatchData = val
		}
	case address >= addresses.ExternalRam && address < addresses.Wram:
		if !m.ramEnabled {
			return
		}

		if m.ramBank <= 0x03 {
			if !m.hasRam {
				return
			}
			ramAddress := int(m.ramBank)*RamBankSize + int(address-addresses.ExternalRam)
			if ramAddress < len(m.ram) {
				m.ram[ramAddress] = val
			}
		} else if m.ramBank >= 0x08 && m.ramBank <= 0x0C && m.hasTimer {
			m.writeRTC(m.ramBank-0x08, val)
		}
	}
}

func (m *MBC3) latchRTC() {
	if m.rtcHalt {
		copy(m.rtcLatch[:], m.rtcRegs[:])
		return
	}

	elapsed := time.Since(m.rtcBaseTime)
	totalSeconds := int(elapsed.Seconds())

	seconds := totalSeconds % 60
	minutes := (totalSeconds / 60) % 60
	hours := (totalSeconds / 3600) % 24
	days := totalSeconds / 86400

	m.rtcLatch[0] = byte(seconds)
	m.rtcLatch[1] = byte(minutes)
	m.rtcLatch[2] = byte(hours)
	m.rtcLatch[3] = byte(days & 0xFF)

	dayHigh := byte((days >> 8) & 0x01)
	carry := byte(0)
	if days > 511 {
		carry = 0x80
	}
	halt := byte(0)
	if m.rtcHalt {
		halt = 0x40
	}

	m.rtcLatch[4] = dayHigh | carry | halt
}

func (m *MBC3) writeRTC(reg byte, val byte) {
	m.rtcRegs[reg] = val

	switch reg {
	case 4: // Days high
		m.rtcHalt = (val & 0x40) != 0
		if !m.rtcHalt {
			m.rtcBaseTime = time.Now()
		}
	}
}
