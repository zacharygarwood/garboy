package main

type CPU struct {
	reg Registers
	mmu MMU

	instr Instruction
}

func NewCPU(cartridge *Cartridge, ppu *PPU) *CPU {
	return &CPU{
		reg: *NewRegisters(),
		mmu: *NewMMU(cartridge, ppu),
	}
}

func (c *CPU) Step() {
	opcode := c.fetch()
	c.decode(opcode)
	c.execute()
}

// Fetches the opcode at PC and increments the PC forward one
func (c *CPU) fetch() byte {
	opcode := c.mmu.Read(c.reg.pc.Read())
	c.reg.pc.Increment()
	return opcode
}

// Decodes the opcode and stores the current instruction
func (c *CPU) decode(opcode byte) {
	c.instr = INSTRUCTIONS[opcode]
}

// Executes the current instruction
func (c *CPU) execute() {
	c.instr.Execute(c)
}
