package main

const (
	TacEnable          = 2
	TacClockSelectMask = 0x03
)

var timerBitPositions = []uint16{9, 3, 5, 7}

type Timer struct {
	internalCounter uint16

	tima uint8
	tma  uint8
	tac  uint8

	overflowDelay uint8

	interrupts *Interrupts
}

func NewTimer(interrupts *Interrupts) *Timer {
	return &Timer{
		internalCounter: 0,
		tima:            0,
		tma:             0,
		tac:             0,
		overflowDelay:   0,
		interrupts:      interrupts,
	}
}

func (t *Timer) Step() {
	prevCounter := t.internalCounter
	t.internalCounter++

	if t.overflowDelay > 0 {
		t.overflowDelay--
		if t.overflowDelay == 0 {
			t.tima = t.tma
			t.interrupts.Request(TimerInterrupt)
		}
	}

	if !t.isTimerEnabled() {
		return
	}

	bitPos := t.getTimerBitPosition()

	prevBit := (prevCounter >> bitPos) & 1
	currentBit := (t.internalCounter >> bitPos) & 1

	if prevBit == 1 && currentBit == 0 {
		t.incrementTimer()
	}
}

func (t *Timer) incrementTimer() {
	if t.tima == 0xFF {
		t.tima = 0x00
		t.overflowDelay = 4 // 1 M-Cycle
	} else {
		t.tima++
	}
}

func (t *Timer) isTimerEnabled() bool {
	return IsBitSet(t.tac, TacEnable)
}

func (t *Timer) getTimerBitPosition() uint16 {
	clockSelect := t.tac & TacClockSelectMask
	return timerBitPositions[clockSelect]
}

func (t *Timer) Read(addr uint16) uint8 {
	switch addr {
	case DivAddress:
		return uint8(t.internalCounter >> 8)
	case TimaAddress:
		return t.tima
	case TmaAddress:
		return t.tma
	case TacAddress:
		return t.tac | 0xF8
	default:
		panic("Reading from Timer using an invalid address")
	}
}

func (t *Timer) Write(addr uint16, value uint8) {
	switch addr {
	case DivAddress:
		t.handleDivReset()
	case TimaAddress:
		if t.overflowDelay > 0 {
			t.overflowDelay = 0
		}
		t.tima = value

	case TmaAddress:
		if t.overflowDelay > 0 {
			t.tima = value
		}
		t.tma = value

	case TacAddress:
		t.tac = value & 0x07
	}
}

func (t *Timer) handleDivReset() {
	if t.isTimerEnabled() {
		bitPos := t.getTimerBitPosition()
		if (t.internalCounter>>bitPos)&1 == 1 {
			t.incrementTimer()
		}
	}

	t.internalCounter = 0
}
