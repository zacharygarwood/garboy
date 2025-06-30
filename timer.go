package main

import "fmt"

var timerFrequencies = []uint16{1024, 16, 64, 256}

type Timer struct {
	tima uint8
	tma  uint8
	tac  uint8

	systemCounter uint16
	timerCounter  uint16

	interrupts *Interrupts
}

func NewTimer(interrupts *Interrupts) *Timer {
	return &Timer{
		interrupts: interrupts,
	}
}

func (t *Timer) Tick(cycles uint16) {
	t.systemCounter += cycles

	if t.isTimerEnabled() {
		t.timerCounter += cycles
		frequency := t.timerFrequency()

		if t.timerCounter >= frequency {
			t.timerCounter -= frequency
			if t.tima == 0xFF {
				t.tima = t.tma
				t.interrupts.Request(2)
			} else {
				t.tima++
			}
		}
	}
}

func (t *Timer) isTimerEnabled() bool {
	return IsBitSet(t.tac, 2)
}

func (t *Timer) timerFrequency() uint16 {
	clockSelect := t.tac & 0x03
	return timerFrequencies[clockSelect]
}

func (t *Timer) Read(address uint16) uint8 {
	switch address {
	case 0xFF04:
		return uint8(t.systemCounter >> 8) // DIV
	case 0xFF05:
		return t.tima
	case 0xFF06:
		return t.tma
	case 0xFF07:
		return t.tac
	default:
		panic("Invalid address trying to read from Timer")
	}
}

func (t *Timer) Write(address uint16, val uint8) {
	switch address {
	case 0xFF04:
		fmt.Printf("[DEBUG] Resetting systemCounter and timerCounter on write to 0xFF04\n")
		// Writing to DIV resets it
		t.systemCounter = 0
		t.timerCounter = 0
	case 0xFF05:
		t.tima = val
	case 0xFF06:
		t.tma = val
	case 0xFF07:
		t.tac = val & 0x07
	}
}
