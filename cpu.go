package main

import "fmt"

type CPU struct {
	reg *Registers
	mmu *MMU

	branched bool
	halted   bool

	interruptMasterEnable bool // IME
}

func NewCPU(cartridge *Cartridge, ppu *PPU) *CPU {
	return &CPU{
		reg: NewRegisters(),
		mmu: NewMMU(cartridge, ppu),
	}
}

func (c *CPU) Step() uint8 {
	opcode := c.fetch()
	instruction := c.decode(opcode)
	return c.execute(instruction)
}

// Fetches the opcode at PC
func (c *CPU) fetch() byte {
	opcode := c.mmu.Read(c.reg.pc.Read())
	c.reg.pc.Increment()
	return opcode
}

// Decodes the opcode returning the instruction. Increments the PC based on the instruction
func (c *CPU) decode(opcode byte) Instruction {
	var instruction Instruction
	if opcode == 0xCB {
		addr := c.getImm8()
		instruction = CB_INSTRUCTIONS[addr]
	} else {
		instruction = INSTRUCTIONS[opcode]
	}
	return instruction
}

// Executes the passed instruction
func (c *CPU) execute(instr Instruction) uint8 {
	c.branched = false // Always set to false before executing
	instr.Execute(&instr, c)

	if c.branched {
		return instr.Cycles.branched
	} else {
		return instr.Cycles.normal
	}
}

func (c *CPU) incrementPC(val uint8) {
	pc := c.reg.pc.Read()
	length := uint16(val)

	c.reg.pc.Write(pc + length)
}

func (c *CPU) SkipBootROM() {
	c.reg.a.Write(0x01)
	c.reg.f.Write(0xB0)
	c.reg.b.Write(0x00)
	c.reg.c.Write(0x13)
	c.reg.d.Write(0x00)
	c.reg.e.Write(0xD8)
	c.reg.h.Write(0x01)
	c.reg.l.Write(0x4D)
	c.reg.sp.Write(0xFFFE)
	c.reg.pc.Write(0x0100)
	c.mmu.bootROMEnabled = false
}

func (c *CPU) PrintState() {
	pc := c.reg.pc.Read()
	fmt.Printf("A:%.2X F:%.2X B:%.2X C:%.2X D:%.2X E:%.2X H:%.2X L:%.2X SP:%.4X PC:%.4X PCMEM:%02X,%02X,%02X,%02X\n",
		c.reg.a.Read(), c.reg.f.Read(), c.reg.b.Read(), c.reg.c.Read(), c.reg.d.Read(), c.reg.e.Read(), c.reg.h.Read(),
		c.reg.l.Read(), c.reg.sp.Read(), pc, c.byteAt(pc).Read(), c.byteAt(pc+1).Read(), c.byteAt(pc+2).Read(), c.byteAt(pc+3).Read())
}
