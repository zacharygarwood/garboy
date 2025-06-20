package main

type CPU struct {
	reg *Registers
	mmu *MMU
}

func NewCPU(cartridge *Cartridge, ppu *PPU) *CPU {
	return &CPU{
		reg: NewRegisters(),
		mmu: NewMMU(cartridge, ppu),
	}
}

func (c *CPU) Step() {
	opcode := c.fetch()
	instruction := c.decode(opcode)
	c.execute(instruction)
}

// Fetches the opcode at PC and increments the PC forward one
func (c *CPU) fetch() byte {
	opcode := c.mmu.Read(c.reg.pc.Read())
	c.reg.pc.Increment()
	return opcode
}

// Decodes the opcode and stores the current instruction
func (c *CPU) decode(opcode byte) Instruction {
	// TODO: Add 0xCB prefix check here to use CB_INSTRUCTIONS
	return INSTRUCTIONS[opcode]
}

// Executes the current instruction
func (c *CPU) execute(instr Instruction) {
	instr.Execute(&instr, c)
}
