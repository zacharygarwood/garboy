package main

type CPU struct {
	reg Registers
	mmu MMU
}

func (c *CPU) Step() {
	// opcode := c.fetch()
	// cycles := c.decodeAndExecute(opcode)
}

func (c *CPU) fetch() byte {
	return c.mmu.Read(c.reg.pc.Read())
}

func (c *CPU) decodeAndExecute(opcode uint8) uint8 {
	// TODO
	return 0x00
}
