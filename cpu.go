package main

import "fmt"

type CPU struct {
	reg Registers
	mmu MMU
}

func NewCPU(cartridge *Cartridge, ppu *PPU) *CPU {
	return &CPU{
		reg: *NewRegisters(),
		mmu: *NewMMU(cartridge, ppu),
	}
}

func (c *CPU) Step() bool {
	opcode := c.fetch()
	return c.decodeAndExecute(opcode)
}

// Fetches the opcode at PC and increments the PC forward one
func (c *CPU) fetch() byte {
	opcode := c.mmu.Read(c.reg.pc.Read())
	c.reg.pc.Increment()
	return opcode
}

func (c *CPU) decodeAndExecute(opcode uint8) bool {
	switch opcode {
	default:
		fmt.Printf("Unknown opcode: %#2x\n", opcode)
		fmt.Printf("cpu.PC=0x%04x\n", c.reg.pc)
		fmt.Printf("cpu.SP=0x%04x\n", c.reg.sp)
		return false
	}
	return true
}
