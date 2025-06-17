package main

type Instruction struct {
	Opcode   byte
	Mnemonic string
	Length   uint8
	Cycles   Cycles // [cycles, branched cycles] where branced cycles may not be present
	Execute  func(c *CPU)
}

type Cycles struct {
	cycles         uint8
	branchedCycles uint8
}

func NewCycles(cycles uint8, branchedCycles uint8) Cycles {
	return Cycles{
		cycles:         cycles,
		branchedCycles: branchedCycles,
	}
}

var INSTRUCTIONS []Instruction = []Instruction{
	{0x00, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0x01, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x02, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x03, "INC BC", 1, NewCycles(4, 4), (*CPU).inc_r16},
	{0x04, "INC B", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x05, "DEC B", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x06, "LD B, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x07, "RLCA", 1, NewCycles(4, 4), (*CPU).rlca},
	{0x08, "LD (u16), SP", 3, NewCycles(20, 20), (*CPU).ld_imm16_sp},
	{0x09, "ADD HL, BC", 1, NewCycles(8, 8), (*CPU).add_hl_r16},
	{0x0A, "LD A, (BC)", 1, NewCycles(8, 8), (*CPU).ld_a_r16mem},
	{0x0B, "DEC BC", 1, NewCycles(8, 8), (*CPU).dec_r16},
	{0x0C, "INC C", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x0D, "DEC C", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x0E, "LD C, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x0F, "RRCA", 1, NewCycles(4, 4), (*CPU).rrca},
	{0x10, "STOP", 1, NewCycles(4, 4), (*CPU).stop},
	{0x11, "LD DE, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x12, "LD (DE), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x13, "INC DE", 1, NewCycles(8, 8), (*CPU).inc_r16},
	{0x14, "INC D", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x15, "DEC D", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x16, "LD D, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x17, "RLA", 1, NewCycles(4, 4), (*CPU).rla},
	{0x18, "JR i8", 2, NewCycles(12, 12), (*CPU).jr_imm8},
	{0x19, "ADD HL, DE", 1, NewCycles(8, 8), (*CPU).add_hl_r16},
	{0x1A, "LD A, (DE)", 1, NewCycles(8, 8), (*CPU).ld_a_r16mem},
	{0x1B, "DEC DE", 1, NewCycles(8, 8), (*CPU).dec_r16},
	{0x1C, "INC E", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x1D, "DEC E", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x1E, "LD E, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x1F, "RRA", 1, NewCycles(4, 4), (*CPU).rra},
	{0x20, "JR NZ, i8", 2, NewCycles(8, 12), (*CPU).jr_cond_imm8},
	{0x21, "LD HL, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x22, "LD (HL+), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x23, "INC HL", 1, NewCycles(8, 8), (*CPU).inc_r16},
	{0x24, "INC H", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x25, "DEC H", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x26, "LD H, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x27, "DAA", 1, NewCycles(4, 4), (*CPU).daa},
	{0x28, "JR Z, i8", 2, NewCycles(8, 12), (*CPU).jr_cond_imm8},
	{0x29, "ADD HL, HL", 1, NewCycles(8, 8), (*CPU).add_hl_r16},
	{0x2A, "LD A, (HL+)", 1, NewCycles(8, 8), (*CPU).ld_a_r16mem},
	{0x2B, "DEC HL", 1, NewCycles(8, 8), (*CPU).dec_r16},
	{0x2C, "INC L", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x2D, "DEC L", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x2E, "LD L, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x2F, "CPL", 1, NewCycles(4, 4), (*CPU).cpl},
	{0x30, "JR NC, i8", 2, NewCycles(8, 12), (*CPU).jr_cond_imm8},
	{0x31, "LD SP, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x32, "LD (HL-), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x33, "INC SP", 1, NewCycles(8, 8), (*CPU).inc_r16},
	{0x34, "", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x35, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x36, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x37, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x38, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x39, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x3A, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x3B, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x3C, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x3D, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x3E, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x3F, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x40, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x41, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x42, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x43, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x44, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0x45, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x46, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x47, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x48, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x49, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x4A, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x4B, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x4C, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x4D, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x4E, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x4F, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x50, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x51, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x52, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x53, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x54, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x55, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0x56, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x57, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x58, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x59, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x5A, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x5B, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x5C, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x5D, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x5E, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x5F, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x60, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x61, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x62, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x63, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x64, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x65, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x66, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0x67, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x68, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x69, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x6A, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x6B, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x6C, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x6D, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x6E, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x6F, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x70, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x71, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x72, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x73, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x74, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x75, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x76, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x77, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0x78, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x79, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x7A, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x7B, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x7C, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x7D, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x7E, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x7F, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x80, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x81, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x82, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x83, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x84, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x85, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x86, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x87, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x88, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x89, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0x8A, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x8B, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x8C, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x8D, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x8E, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x8F, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x90, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x91, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x92, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x93, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x94, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x95, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x96, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x97, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x98, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x99, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x9A, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0x9B, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x9C, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x9D, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x9E, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x9F, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xA0, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xA1, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xA2, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xA3, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xA4, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xA5, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xA6, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xA7, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xA8, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xA9, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xAA, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xAB, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0xAC, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xAD, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0xAE, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xAF, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xB0, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xB1, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xB2, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xB3, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xB4, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xB5, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xB6, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xB7, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xB8, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xB9, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xBA, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xBB, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xBC, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0xBD, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xBE, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0xBF, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xC0, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xC1, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xC2, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xC3, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xC4, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xC5, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xC6, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xC7, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xC8, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xC9, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xCA, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xCB, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xCC, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xCD, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0xCE, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xCF, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0xD0, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xD1, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xD2, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xD3, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xD4, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xD5, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xD6, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xD7, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xD8, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xD9, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xDA, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xDB, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xDC, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xDD, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xDE, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0xDF, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xE0, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0xE1, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xE2, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xE3, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xE4, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xE5, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xE6, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xE7, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xE8, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xE9, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xEA, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xEB, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xEC, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xED, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xEE, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xEF, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xF0, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0xF1, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xF2, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0xF3, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xF4, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xF5, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xF6, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xF7, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xF8, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xF9, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xFA, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xFB, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xFC, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xFD, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xFE, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0xFF, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
}

// Pan Docs reference: https://gbdev.io/pandocs/CPU_Instruction_Set.html

// Block 0
func (c *CPU) nop() {
	// TODO
}

func (c *CPU) ld_r16_imm16() {
	// TODO
}

func (c *CPU) ld_r16mem_a() {
	// TODO
}

func (c *CPU) ld_a_r16mem() {
	// TODO
}

func (c *CPU) ld_imm16_sp() {
	// TODO
}

func (c *CPU) inc_r16() {
	// TODO
}

func (c *CPU) dec_r16() {
	// TODO
}

func (c *CPU) add_hl_r16() {
	// TODO
}

func (c *CPU) inc_r8() {
	// TODO
}

func (c *CPU) dec_r8() {
	// TODO
}

func (c *CPU) ld_r8_imm8() {
	// TODO
}

func (c *CPU) rlca() {
	// TODO
}

func (c *CPU) rrca() {
	// TODO
}

func (c *CPU) rla() {
	// TODO
}

func (c *CPU) rra() {
	// TODO
}

func (c *CPU) daa() {
	// TODO
}

func (c *CPU) cpl() {
	// TODO
}

func (c *CPU) scf() {
	// TODO
}

func (c *CPU) ccf() {
	// TODO
}

func (c *CPU) jr_imm8() {
	// TODO
}

func (c *CPU) jr_cond_imm8() {
	// TODO
}

func (c *CPU) stop() {
	// TODO
}

// Block 1
func (c *CPU) ld_r8_r8() {
	// TODO
	// Exception: ld [hl] [hl] yields the halt instruction
}

func (c *CPU) halt() {
	// TODO
}

// Block 2
func (c *CPU) add_a_r8() {
	// TODO
}

func (c *CPU) adc_a_r8() {
	// TODO
}

func (c *CPU) sub_a_r8() {
	// TODO
}

func (c *CPU) sbc_a_r8() {
	// TODO
}

func (c *CPU) and_a_r8() {
	// TODO
}

func (c *CPU) xor_a_r8() {
	// TODO
}

func (c *CPU) or_a_r8() {
	// TODO
}

func (c *CPU) cp_a_r8() {
	// TODO
}

// Block 3
func (c *CPU) add_a_imm8() {
	// TODO
}

func (c *CPU) adc_a_imm8() {
	// TODO
}

func (c *CPU) sub_a_imm8() {
	// TODO
}

func (c *CPU) sbc_a_imm8() {
	// TODO
}

func (c *CPU) and_a_imm8() {
	// TODO
}

func (c *CPU) xor_a_imm8() {
	// TODO
}

func (c *CPU) or_a_imm8() {
	// TODO
}

func (c *CPU) cp_a_imm8() {
	// TODO
}

func (c *CPU) ret_cond() {
	// TODO
}

func (c *CPU) ret() {
	// TODO
}

func (c *CPU) reti() {
	// TODO
}

func (c *CPU) jp_cond_imm16() {
	// TODO
}

func (c *CPU) jp_imm16() {
	// TODO
}

func (c *CPU) jp_hl() {
	// TODO
}

func (c *CPU) call_cond_imm16() {
	// TODO
}

func (c *CPU) call_imm16() {
	// TODO
}

func (c *CPU) rst_tgt3() {
	// TODO
}

func (c *CPU) pop_r16stk() {
	// TODO
}

func (c *CPU) push_r16stk() {
	// TODO
}

func (c *CPU) ldh_c_a() {
	// TODO
}

func (c *CPU) ldh_imm8_a() {
	// TODO
}

func (c *CPU) ld_imm16_a() {
	// TODO
}

func (c *CPU) ldh_a_c() {
	// TODO
}

func (c *CPU) ldh_a_imm8() {
	// TODO
}

func (c *CPU) ld_a_imm16() {
	// TODO
}

func (c *CPU) add_sp_imm8() {
	// TODO
}

// ld hl, sp + imm8
func (c *CPU) ld_hl_sp_plus_imm8() {
	// TODO
}

func (c *CPU) ld_sp_hl() {
	// TODO
}

func (c *CPU) di() {
	// TODO
}

func (c *CPU) ei() {
	// TODO
}

// 0xCB Prefixed instructions
func (c *CPU) rlc_r8() {
	// TODO
}

func (c *CPU) rrc_r8() {
	// TODO
}

func (c *CPU) rl_r8() {
	// TODO
}

func (c *CPU) rr_r8() {
	// TODO
}

func (c *CPU) sla_r8() {
	// TODO
}

func (c *CPU) sra_r8() {
	// TODO
}

func (c *CPU) swap_r8() {
	// TODO
}

func (c *CPU) srl_r8() {
	// TODO
}

func (c *CPU) bit_b3_r8() {
	// TODO
}

func (c *CPU) res_b3_r8() {
	// TODO
}

func (c *CPU) set_b3_r8() {
	// TODO
}
