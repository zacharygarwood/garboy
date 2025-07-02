package main

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
