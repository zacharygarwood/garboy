package main

const (
	// Interrupt bits
	VBlankInterrupt = 0
	LcdInterrupt    = 1
	TimerInterrupt  = 2
	SerialInterrupt = 3
	JoypadInterrupt = 4

	// Interrupt Sources
	VBlankInterruptSource = 0x40
	StatInterruptSource   = 0x48
	TimerInterruptSource  = 0x50
	SerialInterruptSource = 0x58
	JoypadInterruptSource = 0x60

	InterruptMCycles = 5
)

type Interrupts struct {
	interruptFlag   *InterruptRegister
	interruptEnable *InterruptRegister
}

func NewInterrupts() *Interrupts {
	return &Interrupts{
		interruptFlag:   &InterruptRegister{},
		interruptEnable: &InterruptRegister{},
	}
}

func (i *Interrupts) IF() uint8 {
	return i.interruptFlag.Read()
}

func (i *Interrupts) IE() uint8 {
	return i.interruptEnable.Read()
}

func (i *Interrupts) Write(address uint16, val uint8) {
	switch address {
	case 0xFF0F:
		i.interruptFlag.Write(val)
	case 0xFFFF:
		i.interruptEnable.Write(val)
	default:
		panic("Invalid address when writing to interrupts")
	}
}

func (i *Interrupts) Request(interrupt uint8) {
	res := i.interruptFlag.Read() | (1 << interrupt)
	i.interruptFlag.Write(res)
}

func (i *Interrupts) Clear(interrupt uint8) {
	res := i.interruptFlag.Read() & ^(1 << interrupt)
	i.interruptFlag.Write(res)
}
