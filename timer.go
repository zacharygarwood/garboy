package main

import (
	"fmt"
)

const (
	TacEnable          = 2
	TacClockSelectMask = 0x03
)

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

func (t *Timer) Step(cycles uint16) {
	t.systemCounter += cycles

	if t.isTimerEnabled() {
		t.timerCounter += cycles
		frequency := t.timerFrequency()

		if t.timerCounter >= frequency {
			t.timerCounter -= frequency
			if t.tima == 0xFF {
				t.tima = t.tma
				t.interrupts.Request(TimerInterrupt)
			} else {
				t.tima++
			}
		}
	}
}

func (t *Timer) isTimerEnabled() bool {
	return IsBitSet(t.tac, TacEnable)
}

func (t *Timer) timerFrequency() uint16 {
	clockSelect := t.tac & TacClockSelectMask
	return timerFrequencies[clockSelect]
}

func (t *Timer) Read(address uint16) uint8 {
	switch address {
	case DivAddress:
		return uint8(t.systemCounter >> 8)
	case TimaAddress:
		return t.tima
	case TmaAddress:
		return t.tma
	case TacAddress:
		return t.tac
	default:
		panic("Invalid address trying to read from Timer")
	}
}

func (t *Timer) Write(address uint16, val uint8) {
	switch address {
	case DivAddress:
		// Writing to DIV resets it
		t.systemCounter = 0
		t.timerCounter = 0
	case TimaAddress:
		t.tima = val
	case TmaAddress:
		t.tma = val
	case TacAddress:
		t.tac = val & 0x07
	}
}

func (t *Timer) PrintState() {
	fmt.Printf("[TIMER] DIV:%x TIMA:%x TMA:%x TAC:%x\n", t.systemCounter>>8, t.tima, t.tma, t.tac)
}
