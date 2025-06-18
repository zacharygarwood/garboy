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
	// TODO: Add 0xCB prefix check here to use CB_INSTRUCTIONS
	c.instr = INSTRUCTIONS[opcode]
}

// Executes the current instruction
func (c *CPU) execute() {
	c.instr.Execute(c)
}

// Gets the r16 to use given an index
func (c *CPU) getRegister16(index uint8) Register16 {
	switch index {
	case 0:
		return c.reg.bc
	case 1:
		return c.reg.de
	case 2:
		return c.reg.hl
	case 3:
		return c.reg.sp
	default:
		panic("Invalid index passed to GetRegister16")
	}
}
