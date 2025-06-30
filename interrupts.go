package main

import "fmt"

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
		fmt.Printf("[DEBUG] Writing to IF. Before: %x After: %x\n", i.IF(), val)
		i.interruptFlag.Write(val)
	case 0xFFFF:
		fmt.Printf("[DEBUG] Writing to IE. Before: %x After: %x\n", i.IE(), val)
		i.interruptEnable.Write(val)
	default:
		panic("Invalid address when writing to interrupts")
	}
}

func (i *Interrupts) Request(interrupt uint8) {
	fmt.Printf("[DEBUG] Requesting interrupt %x\n", interrupt)
	res := i.interruptFlag.Read() | (1 << interrupt)
	i.interruptFlag.Write(res)
}

func (i *Interrupts) Clear(interrupt uint8) {
	fmt.Printf("[DEBUG] Clearing interrupt %x\n", interrupt)
	res := i.interruptFlag.Read() & ^(1 << interrupt)
	i.interruptFlag.Write(res)
	fmt.Printf("[DEBUG] After clearning interrupt %x\n", i.IF())
}
