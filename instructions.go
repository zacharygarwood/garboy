package main

import "log"

type Instruction struct {
	Opcode   byte
	Mnemonic string
	Length   uint8
	Cycles   Cycles
	Execute  func(i *Instruction, c *CPU)
}

type Cycles struct {
	branchedCycles uint8
	cycles         uint8
}

func NewCycles(cycles uint8, branchedCycles uint8) Cycles {
	return Cycles{
		cycles:         cycles,
		branchedCycles: branchedCycles,
	}
}

var INVALID_INSTRUCTION Instruction = Instruction{0xFF, "INVALID", 0, NewCycles(1, 1), (*Instruction).invalid_instruction}

var INSTRUCTIONS []Instruction = []Instruction{
	{0x00, "NOP", 1, NewCycles(4, 4), (*Instruction).nop},
	{0x01, "LD BC, u16", 3, NewCycles(12, 12), (*Instruction).ld_r16_imm16},
	{0x02, "LD (BC), A", 1, NewCycles(8, 8), (*Instruction).ld_r16mem_a},
	{0x03, "INC BC", 1, NewCycles(8, 8), (*Instruction).inc_r16},
	{0x04, "INC B", 1, NewCycles(4, 4), (*Instruction).inc_r8},
	{0x05, "DEC B", 1, NewCycles(4, 4), (*Instruction).dec_r8},
	{0x06, "LD B, u8", 2, NewCycles(8, 8), (*Instruction).ld_r8_imm8},
	{0x07, "RLCA", 1, NewCycles(4, 4), (*Instruction).rlca},
	{0x08, "LD (u16), SP", 3, NewCycles(20, 20), (*Instruction).ld_imm16_sp},
	{0x09, "ADD HL, BC", 1, NewCycles(8, 8), (*Instruction).add_hl_r16},
	{0x0a, "LD A, (BC)", 1, NewCycles(8, 8), (*Instruction).ld_a_r16mem},
	{0x0b, "DEC BC", 1, NewCycles(8, 8), (*Instruction).dec_r16},
	{0x0c, "INC C", 1, NewCycles(4, 4), (*Instruction).inc_r8},
	{0x0d, "DEC C", 1, NewCycles(4, 4), (*Instruction).dec_r8},
	{0x0e, "LD C, u8", 2, NewCycles(8, 8), (*Instruction).ld_r8_imm8},
	{0x0f, "RRCA", 1, NewCycles(4, 4), (*Instruction).rrca},
	{0x10, "STOP", 1, NewCycles(4, 4), (*Instruction).stop},
	{0x11, "LD DE, u16", 3, NewCycles(12, 12), (*Instruction).ld_r16_imm16},
	{0x12, "LD (DE), A", 1, NewCycles(8, 8), (*Instruction).ld_r16mem_a},
	{0x13, "INC DE", 1, NewCycles(8, 8), (*Instruction).inc_r16},
	{0x14, "INC D", 1, NewCycles(4, 4), (*Instruction).inc_r8},
	{0x15, "DEC D", 1, NewCycles(4, 4), (*Instruction).dec_r8},
	{0x16, "LD D, u8", 2, NewCycles(8, 8), (*Instruction).ld_r8_imm8},
	{0x17, "RLA", 1, NewCycles(4, 4), (*Instruction).rla},
	{0x18, "JR i8", 2, NewCycles(12, 12), (*Instruction).jr_imm8},
	{0x19, "ADD HL, DE", 1, NewCycles(8, 8), (*Instruction).add_hl_r16},
	{0x1a, "LD A, (DE)", 1, NewCycles(8, 8), (*Instruction).ld_a_r16mem},
	{0x1b, "DEC DE", 1, NewCycles(8, 8), (*Instruction).dec_r16},
	{0x1c, "INC E", 1, NewCycles(4, 4), (*Instruction).inc_r8},
	{0x1d, "DEC E", 1, NewCycles(4, 4), (*Instruction).dec_r8},
	{0x1e, "LD E, u8", 2, NewCycles(8, 8), (*Instruction).ld_r8_imm8},
	{0x1f, "RRA", 1, NewCycles(4, 4), (*Instruction).rra},
	{0x20, "JR NZ, i8", 2, NewCycles(12, 8), (*Instruction).jr_cond_imm8},
	{0x21, "LD HL, u16", 3, NewCycles(12, 12), (*Instruction).ld_r16_imm16},
	{0x22, "LD (HL+), A", 1, NewCycles(8, 8), (*Instruction).ld_r16mem_a},
	{0x23, "INC HL", 1, NewCycles(8, 8), (*Instruction).inc_r16},
	{0x24, "INC H", 1, NewCycles(4, 4), (*Instruction).inc_r8},
	{0x25, "DEC H", 1, NewCycles(4, 4), (*Instruction).dec_r8},
	{0x26, "LD H, u8", 2, NewCycles(8, 8), (*Instruction).ld_r8_imm8},
	{0x27, "DAA", 1, NewCycles(4, 4), (*Instruction).daa},
	{0x28, "JR Z, i8", 2, NewCycles(12, 8), (*Instruction).jr_cond_imm8},
	{0x29, "ADD HL, HL", 1, NewCycles(8, 8), (*Instruction).add_hl_r16},
	{0x2a, "LD A, (HL+)", 1, NewCycles(8, 8), (*Instruction).ld_a_r16mem},
	{0x2b, "DEC HL", 1, NewCycles(8, 8), (*Instruction).dec_r16},
	{0x2c, "INC L", 1, NewCycles(4, 4), (*Instruction).inc_r8},
	{0x2d, "DEC L", 1, NewCycles(4, 4), (*Instruction).dec_r8},
	{0x2e, "LD L, u8", 2, NewCycles(8, 8), (*Instruction).ld_r8_imm8},
	{0x2f, "CPL", 1, NewCycles(4, 4), (*Instruction).cpl},
	{0x30, "JR NC, i8", 2, NewCycles(12, 8), (*Instruction).jr_cond_imm8},
	{0x31, "LD SP, u16", 3, NewCycles(12, 12), (*Instruction).ld_r16_imm16},
	{0x32, "LD (HL-), A", 1, NewCycles(8, 8), (*Instruction).ld_r16mem_a},
	{0x33, "INC SP", 1, NewCycles(8, 8), (*Instruction).inc_r16},
	{0x34, "INC (HL)", 1, NewCycles(12, 12), (*Instruction).inc_r8},
	{0x35, "DEC (HL)", 1, NewCycles(12, 12), (*Instruction).dec_r8},
	{0x36, "LD (HL), u8", 2, NewCycles(12, 12), (*Instruction).ld_r8_imm8},
	{0x37, "SCF", 1, NewCycles(4, 4), (*Instruction).scf},
	{0x38, "JR C, i8", 2, NewCycles(12, 8), (*Instruction).jr_cond_imm8},
	{0x39, "ADD HL, SP", 1, NewCycles(8, 8), (*Instruction).add_hl_r16},
	{0x3a, "LD A, (HL-)", 1, NewCycles(8, 8), (*Instruction).ld_a_r16mem},
	{0x3b, "DEC SP", 1, NewCycles(8, 8), (*Instruction).dec_r16},
	{0x3c, "INC A", 1, NewCycles(4, 4), (*Instruction).inc_r8},
	{0x3d, "DEC A", 1, NewCycles(4, 4), (*Instruction).dec_r8},
	{0x3e, "LD A, u8", 2, NewCycles(8, 8), (*Instruction).ld_r8_imm8},
	{0x3f, "CCF", 1, NewCycles(4, 4), (*Instruction).ccf},
	{0x40, "LD B, B", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x41, "LD B, C", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x42, "LD B, D", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x43, "LD B, E", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x44, "LD B, H", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x45, "LD B, L", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x46, "LD B, (HL)", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x47, "LD B, A", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x48, "LD C, B", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x49, "LD C, C", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x4a, "LD C, D", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x4b, "LD C, E", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x4c, "LD C, H", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x4d, "LD C, L", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x4e, "LD C, (HL)", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x4f, "LD C, A", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x50, "LD D, B", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x51, "LD D, C", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x52, "LD D, D", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x53, "LD D, E", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x54, "LD D, H", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x55, "LD D, L", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x56, "LD D, (HL)", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x57, "LD D, A", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x58, "LD E, B", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x59, "LD E, C", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x5a, "LD E, D", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x5b, "LD E, E", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x5c, "LD E, H", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x5d, "LD E, L", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x5e, "LD E, (HL)", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x5f, "LD E, A", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x60, "LD H, B", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x61, "LD H, C", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x62, "LD H, D", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x63, "LD H, E", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x64, "LD H, H", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x65, "LD H, L", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x66, "LD H, (HL)", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x67, "LD H, A", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x68, "LD L, B", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x69, "LD L, C", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x6a, "LD L, D", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x6b, "LD L, E", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x6c, "LD L, H", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x6d, "LD L, L", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x6e, "LD L, (HL)", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x6f, "LD L, A", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x70, "LD (HL), B", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x71, "LD (HL), C", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x72, "LD (HL), D", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x73, "LD (HL), E", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x74, "LD (HL), H", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x75, "LD (HL), L", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x76, "HALT", 1, NewCycles(4, 4), (*Instruction).halt},
	{0x77, "LD (HL), A", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x78, "LD A, B", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x79, "LD A, C", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x7a, "LD A, D", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x7b, "LD A, E", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x7c, "LD A, H", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x7d, "LD A, L", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x7e, "LD A, (HL)", 1, NewCycles(8, 8), (*Instruction).ld_r8_r8},
	{0x7f, "LD A, A", 1, NewCycles(4, 4), (*Instruction).ld_r8_r8},
	{0x80, "ADD A, B", 1, NewCycles(4, 4), (*Instruction).add_a_r8},
	{0x81, "ADD A, C", 1, NewCycles(4, 4), (*Instruction).add_a_r8},
	{0x82, "ADD A, D", 1, NewCycles(4, 4), (*Instruction).add_a_r8},
	{0x83, "ADD A, E", 1, NewCycles(4, 4), (*Instruction).add_a_r8},
	{0x84, "ADD A, H", 1, NewCycles(4, 4), (*Instruction).add_a_r8},
	{0x85, "ADD A, L", 1, NewCycles(4, 4), (*Instruction).add_a_r8},
	{0x86, "ADD A, (HL)", 1, NewCycles(8, 8), (*Instruction).add_a_r8},
	{0x87, "ADD A, A", 1, NewCycles(4, 4), (*Instruction).add_a_r8},
	{0x88, "ADC A, B", 1, NewCycles(4, 4), (*Instruction).adc_a_r8},
	{0x89, "ADC A, C", 1, NewCycles(4, 4), (*Instruction).adc_a_r8},
	{0x8a, "ADC A, D", 1, NewCycles(4, 4), (*Instruction).adc_a_r8},
	{0x8b, "ADC A, E", 1, NewCycles(4, 4), (*Instruction).adc_a_r8},
	{0x8c, "ADC A, H", 1, NewCycles(4, 4), (*Instruction).adc_a_r8},
	{0x8d, "ADC A, L", 1, NewCycles(4, 4), (*Instruction).adc_a_r8},
	{0x8e, "ADC A, (HL)", 1, NewCycles(8, 8), (*Instruction).adc_a_r8},
	{0x8f, "ADC A, A", 1, NewCycles(4, 4), (*Instruction).adc_a_r8},
	{0x90, "SUB A, B", 1, NewCycles(4, 4), (*Instruction).sub_a_r8},
	{0x91, "SUB A, C", 1, NewCycles(4, 4), (*Instruction).sub_a_r8},
	{0x92, "SUB A, D", 1, NewCycles(4, 4), (*Instruction).sub_a_r8},
	{0x93, "SUB A, E", 1, NewCycles(4, 4), (*Instruction).sub_a_r8},
	{0x94, "SUB A, H", 1, NewCycles(4, 4), (*Instruction).sub_a_r8},
	{0x95, "SUB A, L", 1, NewCycles(4, 4), (*Instruction).sub_a_r8},
	{0x96, "SUB A, (HL)", 1, NewCycles(8, 8), (*Instruction).sub_a_r8},
	{0x97, "SUB A, A", 1, NewCycles(4, 4), (*Instruction).sub_a_r8},
	{0x98, "SBC A, B", 1, NewCycles(4, 4), (*Instruction).sbc_a_r8},
	{0x99, "SBC A, C", 1, NewCycles(4, 4), (*Instruction).sbc_a_r8},
	{0x9a, "SBC A, D", 1, NewCycles(4, 4), (*Instruction).sbc_a_r8},
	{0x9b, "SBC A, E", 1, NewCycles(4, 4), (*Instruction).sbc_a_r8},
	{0x9c, "SBC A, H", 1, NewCycles(4, 4), (*Instruction).sbc_a_r8},
	{0x9d, "SBC A, L", 1, NewCycles(4, 4), (*Instruction).sbc_a_r8},
	{0x9e, "SBC A, (HL)", 1, NewCycles(8, 8), (*Instruction).sbc_a_r8},
	{0x9f, "SBC A, A", 1, NewCycles(4, 4), (*Instruction).sbc_a_r8},
	{0xa0, "AND A, B", 1, NewCycles(4, 4), (*Instruction).and_a_r8},
	{0xa1, "AND A, C", 1, NewCycles(4, 4), (*Instruction).and_a_r8},
	{0xa2, "AND A, D", 1, NewCycles(4, 4), (*Instruction).and_a_r8},
	{0xa3, "AND A, E", 1, NewCycles(4, 4), (*Instruction).and_a_r8},
	{0xa4, "AND A, H", 1, NewCycles(4, 4), (*Instruction).and_a_r8},
	{0xa5, "AND A, L", 1, NewCycles(4, 4), (*Instruction).and_a_r8},
	{0xa6, "AND A, (HL)", 1, NewCycles(8, 8), (*Instruction).and_a_r8},
	{0xa7, "AND A, A", 1, NewCycles(4, 4), (*Instruction).and_a_r8},
	{0xa8, "XOR A, B", 1, NewCycles(4, 4), (*Instruction).xor_a_r8},
	{0xa9, "XOR A, C", 1, NewCycles(4, 4), (*Instruction).xor_a_r8},
	{0xaa, "XOR A, D", 1, NewCycles(4, 4), (*Instruction).xor_a_r8},
	{0xab, "XOR A, E", 1, NewCycles(4, 4), (*Instruction).xor_a_r8},
	{0xac, "XOR A, H", 1, NewCycles(4, 4), (*Instruction).xor_a_r8},
	{0xad, "XOR A, L", 1, NewCycles(4, 4), (*Instruction).xor_a_r8},
	{0xae, "XOR A, (HL)", 1, NewCycles(8, 8), (*Instruction).xor_a_r8},
	{0xaf, "XOR A, A", 1, NewCycles(4, 4), (*Instruction).xor_a_r8},
	{0xb0, "OR A, B", 1, NewCycles(4, 4), (*Instruction).or_a_r8},
	{0xb1, "OR A, C", 1, NewCycles(4, 4), (*Instruction).or_a_r8},
	{0xb2, "OR A, D", 1, NewCycles(4, 4), (*Instruction).or_a_r8},
	{0xb3, "OR A, E", 1, NewCycles(4, 4), (*Instruction).or_a_r8},
	{0xb4, "OR A, H", 1, NewCycles(4, 4), (*Instruction).or_a_r8},
	{0xb5, "OR A, L", 1, NewCycles(4, 4), (*Instruction).or_a_r8},
	{0xb6, "OR (HL)", 1, NewCycles(8, 8), (*Instruction).or_a_r8},
	{0xb7, "OR A, A", 1, NewCycles(4, 4), (*Instruction).or_a_r8},
	{0xb8, "CP A, B", 1, NewCycles(4, 4), (*Instruction).cp_a_r8},
	{0xb9, "CP A, C", 1, NewCycles(4, 4), (*Instruction).cp_a_r8},
	{0xba, "CP A, D", 1, NewCycles(4, 4), (*Instruction).cp_a_r8},
	{0xbb, "CP A, E", 1, NewCycles(4, 4), (*Instruction).cp_a_r8},
	{0xbc, "CP A, H", 1, NewCycles(4, 4), (*Instruction).cp_a_r8},
	{0xbd, "CP A, L", 1, NewCycles(4, 4), (*Instruction).cp_a_r8},
	{0xbe, "CP (HL)", 1, NewCycles(8, 8), (*Instruction).cp_a_r8},
	{0xbf, "CP A", 1, NewCycles(4, 4), (*Instruction).cp_a_r8},
	{0xc0, "RET NZ", 1, NewCycles(20, 8), (*Instruction).ret_cond},
	{0xc1, "POP BC", 1, NewCycles(12, 12), (*Instruction).pop_r16stk},
	{0xc2, "JP NZ, u16", 3, NewCycles(16, 12), (*Instruction).jp_cond_imm16},
	{0xc3, "JP u16", 3, NewCycles(16, 16), (*Instruction).jp_imm16},
	{0xc4, "CALL NZ, u16", 3, NewCycles(24, 12), (*Instruction).call_cond_imm16},
	{0xc5, "PUSH BC", 1, NewCycles(16, 16), (*Instruction).push_r16stk},
	{0xc6, "ADD A, u8", 2, NewCycles(8, 8), (*Instruction).add_a_imm8},
	{0xc7, "RST 00H", 1, NewCycles(16, 16), (*Instruction).rst_tgt3},
	{0xc8, "RET Z", 1, NewCycles(20, 8), (*Instruction).ret_cond},
	{0xc9, "RET", 1, NewCycles(16, 16), (*Instruction).ret},
	{0xca, "JP Z, u16", 3, NewCycles(16, 12), (*Instruction).jp_cond_imm16},
	INVALID_INSTRUCTION,
	{0xcc, "CALL Z, u16", 3, NewCycles(24, 12), (*Instruction).call_cond_imm16},
	{0xcd, "CALL u16", 3, NewCycles(24, 24), (*Instruction).call_imm16},
	{0xce, "ADC A, u8", 2, NewCycles(8, 8), (*Instruction).adc_a_imm8},
	{0xcf, "RST 08H", 1, NewCycles(16, 16), (*Instruction).rst_tgt3},
	{0xd0, "RET NC", 1, NewCycles(20, 8), (*Instruction).ret_cond},
	{0xd1, "POP DE", 1, NewCycles(12, 12), (*Instruction).pop_r16stk},
	{0xd2, "JP NC, u16", 3, NewCycles(16, 12), (*Instruction).jp_cond_imm16},
	INVALID_INSTRUCTION,
	{0xd4, "CALL NC, u16", 3, NewCycles(24, 12), (*Instruction).call_cond_imm16},
	{0xd5, "PUSH DE", 1, NewCycles(16, 16), (*Instruction).push_r16stk},
	{0xd6, "SUB A, u8", 2, NewCycles(8, 8), (*Instruction).sub_a_imm8},
	{0xd7, "RST 10H", 1, NewCycles(16, 16), (*Instruction).rst_tgt3},
	{0xd8, "RET C", 1, NewCycles(20, 8), (*Instruction).ret_cond},
	{0xd9, "RETI", 1, NewCycles(16, 16), (*Instruction).reti},
	{0xda, "JP C, u16", 3, NewCycles(16, 12), (*Instruction).jp_cond_imm16},
	INVALID_INSTRUCTION,
	{0xdc, "CALL C, u16", 3, NewCycles(24, 12), (*Instruction).call_cond_imm16},
	INVALID_INSTRUCTION,
	{0xde, "SBC A, u8", 2, NewCycles(8, 8), (*Instruction).sbc_a_imm8},
	{0xdf, "RST 18H", 1, NewCycles(16, 16), (*Instruction).rst_tgt3},
	{0xe0, "LDH (u8), A", 2, NewCycles(12, 12), (*Instruction).ldh_imm8_a},
	{0xe1, "POP HL", 1, NewCycles(12, 12), (*Instruction).pop_r16stk},
	{0xe2, "LD (C), A", 1, NewCycles(8, 8), (*Instruction).ldh_c_a},
	INVALID_INSTRUCTION,
	INVALID_INSTRUCTION,
	{0xe5, "PUSH HL", 1, NewCycles(16, 16), (*Instruction).push_r16stk},
	{0xe6, "AND A, u8", 2, NewCycles(8, 8), (*Instruction).and_a_imm8},
	{0xe7, "RST 20H", 1, NewCycles(16, 16), (*Instruction).rst_tgt3},
	{0xe8, "ADD SP, i8", 2, NewCycles(16, 16), (*Instruction).add_sp_imm8},
	{0xe9, "JP HL", 1, NewCycles(4, 4), (*Instruction).jp_hl},
	{0xea, "LD (u16), A", 3, NewCycles(16, 16), (*Instruction).ld_imm16_a},
	INVALID_INSTRUCTION,
	INVALID_INSTRUCTION,
	INVALID_INSTRUCTION,
	{0xee, "XOR A, u8", 2, NewCycles(8, 8), (*Instruction).xor_a_imm8},
	{0xef, "RST 28H", 1, NewCycles(16, 16), (*Instruction).rst_tgt3},
	{0xf0, "LDH A, (u8)", 2, NewCycles(12, 12), (*Instruction).ldh_a_imm8},
	{0xf1, "POP AF", 1, NewCycles(12, 12), (*Instruction).pop_r16stk},
	{0xf2, "LD A, (C)", 1, NewCycles(8, 8), (*Instruction).ldh_a_c},
	{0xf3, "DI", 1, NewCycles(4, 4), (*Instruction).di},
	INVALID_INSTRUCTION,
	{0xf5, "PUSH AF", 1, NewCycles(16, 16), (*Instruction).push_r16stk},
	{0xf6, "OR A, u8", 2, NewCycles(8, 8), (*Instruction).or_a_imm8},
	{0xf7, "RST 30H", 1, NewCycles(16, 16), (*Instruction).rst_tgt3},
	{0xf8, "LD HL, SP+r8", 2, NewCycles(12, 12), (*Instruction).ld_hl_sp_plus_imm8},
	{0xf9, "LD SP, HL", 1, NewCycles(8, 8), (*Instruction).ld_sp_hl},
	{0xfa, "LD A, (u16)", 3, NewCycles(16, 16), (*Instruction).ld_a_imm16},
	{0xfb, "EI", 1, NewCycles(4, 4), (*Instruction).ei},
	{0xfe, "CP A, u8", 2, NewCycles(8, 8), (*Instruction).cp_a_imm8},
	{0xff, "RST 38H", 1, NewCycles(16, 16), (*Instruction).rst_tgt3},
}

var CB_INSTRUCTIONS []Instruction = []Instruction{
	{0x00, "RLC B", 2, NewCycles(8, 8), (*Instruction).rlc_r8},
	{0x01, "RLC C", 2, NewCycles(8, 8), (*Instruction).rlc_r8},
	{0x02, "RLC D", 2, NewCycles(8, 8), (*Instruction).rlc_r8},
	{0x03, "RLC E", 2, NewCycles(8, 8), (*Instruction).rlc_r8},
	{0x04, "RLC H", 2, NewCycles(8, 8), (*Instruction).rlc_r8},
	{0x05, "RLC L", 2, NewCycles(8, 8), (*Instruction).rlc_r8},
	{0x06, "RLC (HL)", 2, NewCycles(16, 16), (*Instruction).rlc_r8},
	{0x07, "RLC A", 2, NewCycles(8, 8), (*Instruction).rlc_r8},
	{0x08, "RRC B", 2, NewCycles(8, 8), (*Instruction).rrc_r8},
	{0x09, "RRC C", 2, NewCycles(8, 8), (*Instruction).rrc_r8},
	{0x0a, "RRC D", 2, NewCycles(8, 8), (*Instruction).rrc_r8},
	{0x0b, "RRC E", 2, NewCycles(8, 8), (*Instruction).rrc_r8},
	{0x0c, "RRC H", 2, NewCycles(8, 8), (*Instruction).rrc_r8},
	{0x0d, "RRC L", 2, NewCycles(8, 8), (*Instruction).rrc_r8},
	{0x0e, "RRC (HL)", 2, NewCycles(16, 16), (*Instruction).rrc_r8},
	{0x0f, "RRC A", 2, NewCycles(8, 8), (*Instruction).rrc_r8},
	{0x10, "RL B", 2, NewCycles(8, 8), (*Instruction).rl_r8},
	{0x11, "RL C", 2, NewCycles(8, 8), (*Instruction).rl_r8},
	{0x12, "RL D", 2, NewCycles(8, 8), (*Instruction).rl_r8},
	{0x13, "RL E", 2, NewCycles(8, 8), (*Instruction).rl_r8},
	{0x14, "RL H", 2, NewCycles(8, 8), (*Instruction).rl_r8},
	{0x15, "RL L", 2, NewCycles(8, 8), (*Instruction).rl_r8},
	{0x16, "RL (HL)", 2, NewCycles(16, 16), (*Instruction).rl_r8},
	{0x17, "RL A", 2, NewCycles(8, 8), (*Instruction).rl_r8},
	{0x18, "RR B", 2, NewCycles(8, 8), (*Instruction).rr_r8},
	{0x19, "RR C", 2, NewCycles(8, 8), (*Instruction).rr_r8},
	{0x1a, "RR D", 2, NewCycles(8, 8), (*Instruction).rr_r8},
	{0x1b, "RR E", 2, NewCycles(8, 8), (*Instruction).rr_r8},
	{0x1c, "RR H", 2, NewCycles(8, 8), (*Instruction).rr_r8},
	{0x1d, "RR L", 2, NewCycles(8, 8), (*Instruction).rr_r8},
	{0x1e, "RR (HL)", 2, NewCycles(16, 16), (*Instruction).rr_r8},
	{0x1f, "RR A", 2, NewCycles(8, 8), (*Instruction).rr_r8},
	{0x20, "SLA B", 2, NewCycles(8, 8), (*Instruction).sla_r8},
	{0x21, "SLA C", 2, NewCycles(8, 8), (*Instruction).sla_r8},
	{0x22, "SLA D", 2, NewCycles(8, 8), (*Instruction).sla_r8},
	{0x23, "SLA E", 2, NewCycles(8, 8), (*Instruction).sla_r8},
	{0x24, "SLA H", 2, NewCycles(8, 8), (*Instruction).sla_r8},
	{0x25, "SLA L", 2, NewCycles(8, 8), (*Instruction).sla_r8},
	{0x26, "SLA (HL)", 2, NewCycles(16, 16), (*Instruction).sla_r8},
	{0x27, "SLA A", 2, NewCycles(8, 8), (*Instruction).sla_r8},
	{0x28, "SRA B", 2, NewCycles(8, 8), (*Instruction).sra_r8},
	{0x29, "SRA C", 2, NewCycles(8, 8), (*Instruction).sra_r8},
	{0x2a, "SRA D", 2, NewCycles(8, 8), (*Instruction).sra_r8},
	{0x2b, "SRA E", 2, NewCycles(8, 8), (*Instruction).sra_r8},
	{0x2c, "SRA H", 2, NewCycles(8, 8), (*Instruction).sra_r8},
	{0x2d, "SRA L", 2, NewCycles(8, 8), (*Instruction).sra_r8},
	{0x2e, "SRA (HL)", 2, NewCycles(16, 16), (*Instruction).sra_r8},
	{0x2f, "SRA A", 2, NewCycles(8, 8), (*Instruction).sra_r8},
	{0x30, "SWAP B", 2, NewCycles(8, 8), (*Instruction).swap_r8},
	{0x31, "SWAP C", 2, NewCycles(8, 8), (*Instruction).swap_r8},
	{0x32, "SWAP D", 2, NewCycles(8, 8), (*Instruction).swap_r8},
	{0x33, "SWAP E", 2, NewCycles(8, 8), (*Instruction).swap_r8},
	{0x34, "SWAP H", 2, NewCycles(8, 8), (*Instruction).swap_r8},
	{0x35, "SWAP L", 2, NewCycles(8, 8), (*Instruction).swap_r8},
	{0x36, "SWAP (HL)", 2, NewCycles(16, 16), (*Instruction).swap_r8},
	{0x37, "SWAP A", 2, NewCycles(8, 8), (*Instruction).swap_r8},
	{0x38, "SRL B", 2, NewCycles(8, 8), (*Instruction).srl_r8},
	{0x39, "SRL C", 2, NewCycles(8, 8), (*Instruction).srl_r8},
	{0x3a, "SRL D", 2, NewCycles(8, 8), (*Instruction).srl_r8},
	{0x3b, "SRL E", 2, NewCycles(8, 8), (*Instruction).srl_r8},
	{0x3c, "SRL H", 2, NewCycles(8, 8), (*Instruction).srl_r8},
	{0x3d, "SRL L", 2, NewCycles(8, 8), (*Instruction).srl_r8},
	{0x3e, "SRL (HL)", 2, NewCycles(16, 16), (*Instruction).srl_r8},
	{0x3f, "SRL A", 2, NewCycles(8, 8), (*Instruction).srl_r8},
	{0x40, "BIT 0, B", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x41, "BIT 0, C", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x42, "BIT 0, D", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x43, "BIT 0, E", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x44, "BIT 0, H", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x45, "BIT 0, L", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x46, "BIT 0, (HL)", 2, NewCycles(16, 16), (*Instruction).bit_b3_r8},
	{0x47, "BIT 0, A", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x48, "BIT 1, B", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x49, "BIT 1, C", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x4a, "BIT 1, D", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x4b, "BIT 1, E", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x4c, "BIT 1, H", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x4d, "BIT 1, L", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x4e, "BIT 1, (HL)", 2, NewCycles(16, 16), (*Instruction).bit_b3_r8},
	{0x4f, "BIT 1, A", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x50, "BIT 2, B", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x51, "BIT 2, C", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x52, "BIT 2, D", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x53, "BIT 2, E", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x54, "BIT 2, H", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x55, "BIT 2, L", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x56, "BIT 2, (HL)", 2, NewCycles(16, 16), (*Instruction).bit_b3_r8},
	{0x57, "BIT 2, A", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x58, "BIT 3, B", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x59, "BIT 3, C", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x5a, "BIT 3, D", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x5b, "BIT 3, E", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x5c, "BIT 3, H", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x5d, "BIT 3, L", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x5e, "BIT 3, (HL)", 2, NewCycles(16, 16), (*Instruction).bit_b3_r8},
	{0x5f, "BIT 3, A", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x60, "BIT 4, B", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x61, "BIT 4, C", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x62, "BIT 4, D", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x63, "BIT 4, E", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x64, "BIT 4, H", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x65, "BIT 4, L", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x66, "BIT 4, (HL)", 2, NewCycles(16, 16), (*Instruction).bit_b3_r8},
	{0x67, "BIT 4, A", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x68, "BIT 5, B", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x69, "BIT 5, C", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x6a, "BIT 5, D", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x6b, "BIT 5, E", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x6c, "BIT 5, H", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x6d, "BIT 5, L", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x6e, "BIT 5, (HL)", 2, NewCycles(16, 16), (*Instruction).bit_b3_r8},
	{0x6f, "BIT 5, A", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x70, "BIT 6, B", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x71, "BIT 6, C", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x72, "BIT 6, D", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x73, "BIT 6, E", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x74, "BIT 6, H", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x75, "BIT 6, L", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x76, "BIT 6, (HL)", 2, NewCycles(16, 16), (*Instruction).bit_b3_r8},
	{0x77, "BIT 6, A", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x78, "BIT 7, B", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x79, "BIT 7, C", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x7a, "BIT 7, D", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x7b, "BIT 7, E", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x7c, "BIT 7, H", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x7d, "BIT 7, L", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x7e, "BIT 7, (HL)", 2, NewCycles(16, 16), (*Instruction).bit_b3_r8},
	{0x7f, "BIT 7, A", 2, NewCycles(8, 8), (*Instruction).bit_b3_r8},
	{0x80, "RES 0, B", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x81, "RES 0, C", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x82, "RES 0, D", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x83, "RES 0, E", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x84, "RES 0, H", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x85, "RES 0, L", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x86, "RES 0, (HL)", 2, NewCycles(16, 16), (*Instruction).res_b3_r8},
	{0x87, "RES 0, A", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x88, "RES 1, B", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x89, "RES 1, C", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x8a, "RES 1, D", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x8b, "RES 1, E", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x8c, "RES 1, H", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x8d, "RES 1, L", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x8e, "RES 1, (HL)", 2, NewCycles(16, 16), (*Instruction).res_b3_r8},
	{0x8f, "RES 1, A", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x90, "RES 2, B", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x91, "RES 2, C", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x92, "RES 2, D", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x93, "RES 2, E", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x94, "RES 2, H", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x95, "RES 2, L", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x96, "RES 2, (HL)", 2, NewCycles(16, 16), (*Instruction).res_b3_r8},
	{0x97, "RES 2, A", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x98, "RES 3, B", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x99, "RES 3, C", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x9a, "RES 3, D", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x9b, "RES 3, E", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x9c, "RES 3, H", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x9d, "RES 3, L", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0x9e, "RES 3, (HL)", 2, NewCycles(16, 16), (*Instruction).res_b3_r8},
	{0x9f, "RES 3, A", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xa0, "RES 4, B", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xa1, "RES 4, C", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xa2, "RES 4, D", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xa3, "RES 4, E", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xa4, "RES 4, H", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xa5, "RES 4, L", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xa6, "RES 4, (HL)", 2, NewCycles(16, 16), (*Instruction).res_b3_r8},
	{0xa7, "RES 4, A", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xa8, "RES 5, B", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xa9, "RES 5, C", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xaa, "RES 5, D", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xab, "RES 5, E", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xac, "RES 5, H", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xad, "RES 5, L", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xae, "RES 5, (HL)", 2, NewCycles(16, 16), (*Instruction).res_b3_r8},
	{0xaf, "RES 5, A", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xb0, "RES 6, B", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xb1, "RES 6, C", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xb2, "RES 6, D", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xb3, "RES 6, E", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xb4, "RES 6, H", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xb5, "RES 6, L", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xb6, "RES 6, (HL)", 2, NewCycles(16, 16), (*Instruction).res_b3_r8},
	{0xb7, "RES 6, A", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xb8, "RES 7, B", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xb9, "RES 7, C", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xba, "RES 7, D", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xbb, "RES 7, E", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xbc, "RES 7, H", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xbd, "RES 7, L", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xbe, "RES 7, (HL)", 2, NewCycles(16, 16), (*Instruction).res_b3_r8},
	{0xbf, "RES 7, A", 2, NewCycles(8, 8), (*Instruction).res_b3_r8},
	{0xc0, "SET 0, B", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xc1, "SET 0, C", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xc2, "SET 0, D", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xc3, "SET 0, E", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xc4, "SET 0, H", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xc5, "SET 0, L", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xc6, "SET 0, (HL)", 2, NewCycles(16, 16), (*Instruction).set_b3_r8},
	{0xc7, "SET 0, A", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xc8, "SET 1, B", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xc9, "SET 1, C", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xca, "SET 1, D", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xcb, "SET 1, E", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xcc, "SET 1, H", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xcd, "SET 1, L", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xce, "SET 1, (HL)", 2, NewCycles(16, 16), (*Instruction).set_b3_r8},
	{0xcf, "SET 1, A", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xd0, "SET 2, B", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xd1, "SET 2, C", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xd2, "SET 2, D", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xd3, "SET 2, E", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xd4, "SET 2, H", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xd5, "SET 2, L", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xd6, "SET 2, (HL)", 2, NewCycles(16, 16), (*Instruction).set_b3_r8},
	{0xd7, "SET 2, A", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xd8, "SET 3, B", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xd9, "SET 3, C", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xda, "SET 3, D", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xdb, "SET 3, E", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xdc, "SET 3, H", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xdd, "SET 3, L", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xde, "SET 3, (HL)", 2, NewCycles(16, 16), (*Instruction).set_b3_r8},
	{0xdf, "SET 3, A", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xe0, "SET 4, B", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xe1, "SET 4, C", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xe2, "SET 4, D", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xe3, "SET 4, E", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xe4, "SET 4, H", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xe5, "SET 4, L", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xe6, "SET 4, (HL)", 2, NewCycles(16, 16), (*Instruction).set_b3_r8},
	{0xe7, "SET 4, A", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xe8, "SET 5, B", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xe9, "SET 5, C", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xea, "SET 5, D", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xeb, "SET 5, E", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xec, "SET 5, H", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xed, "SET 5, L", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xee, "SET 5, (HL)", 2, NewCycles(16, 16), (*Instruction).set_b3_r8},
	{0xef, "SET 5, A", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xf0, "SET 6, B", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xf1, "SET 6, C", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xf2, "SET 6, D", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xf3, "SET 6, E", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xf4, "SET 6, H", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xf5, "SET 6, L", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xf6, "SET 6, (HL)", 2, NewCycles(16, 16), (*Instruction).set_b3_r8},
	{0xf7, "SET 6, A", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xf8, "SET 7, B", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xf9, "SET 7, C", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xfa, "SET 7, D", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xfb, "SET 7, E", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xfc, "SET 7, H", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xfd, "SET 7, L", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
	{0xfe, "SET 7, (HL)", 2, NewCycles(16, 16), (*Instruction).set_b3_r8},
	{0xff, "SET 7, A", 2, NewCycles(8, 8), (*Instruction).set_b3_r8},
}

func (c *CPU) getImm8() uint8 {
	imm := c.mmu.Read(c.reg.pc.Read())
	c.reg.pc.Increment()
	return imm
}

func (c *CPU) getImm16() uint16 {
	lsb := uint16(c.getImm8())
	msb := uint16(c.getImm8())
	return (msb << 8) | lsb
}

// MemoryReference8 conforms to the Register8 interface
func (c *CPU) byteAt(addr uint16) *MemoryReference8 {
	return &MemoryReference8{
		cpu:  c,
		addr: addr,
	}
}

func (c *CPU) getRegister8(opcode byte, bits []int) Register8 {
	index := ExtractBits(opcode, bits)
	switch index {
	case 0:
		return c.reg.b
	case 1:
		return c.reg.c
	case 2:
		return c.reg.d
	case 3:
		return c.reg.e
	case 4:
		return c.reg.h
	case 5:
		return c.reg.l
	case 6:
		return c.byteAt(c.reg.hl.Read())
	case 7:
		return c.reg.a
	default:
		panic("Invalid index passed to getRegister8()")
	}
}

func (c *CPU) getRegister16(opcode byte, bits []int) Register16 {
	index := ExtractBits(opcode, bits)
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
		panic("Invalid index passed to getRegister16()")
	}
}

func (c *CPU) getRegister16Stk(opcode byte, bits []int) Register16 {
	index := ExtractBits(opcode, bits)
	switch index {
	case 0:
		return c.reg.bc
	case 1:
		return c.reg.de
	case 2:
		return c.reg.hl
	case 3:
		return c.reg.af
	default:
		panic("Invalid index passed to getRegister16Stk()")
	}
}

func (c *CPU) getRegister16Mem(opcode byte, bits []int) *MemoryReference8 {
	index := ExtractBits(opcode, bits)
	switch index {
	case 0:
		return c.byteAt(c.reg.bc.Read())
	case 1:
		return c.byteAt(c.reg.de.Read())
	case 2:
		hl := c.byteAt(c.reg.hl.Read())
		c.reg.hl.Increment()
		return hl
	case 3:
		hl := c.byteAt(c.reg.hl.Read())
		c.reg.hl.Decrement()
		return hl
	default:
		panic("Invalid index passed to getRegister16Mem()")
	}
}

func (c *CPU) getCond(opcode byte, bits []int) bool {
	index := ExtractBits(opcode, bits)
	switch index {
	case 0:
		return !c.reg.f.Z()
	case 1:
		return c.reg.f.Z()
	case 2:
		return !c.reg.f.C()
	case 3:
		return c.reg.f.C()
	default:
		panic("Invalid index passed to getCond()")
	}
}

func (c *Instruction) invalid_instruction() {
	panic("Calling invalid instruction")
}

// Block 0
func (i *Instruction) nop(c *CPU) {}

func (i *Instruction) ld_r16_imm16(c *CPU) {
	r16 := c.getRegister16(i.Opcode, []int{5, 4})
	imm16 := c.getImm16()

	r16.Write(imm16)
}

func (i *Instruction) ld_r16mem_a(c *CPU) {
	r16mem := c.getRegister16Mem(i.Opcode, []int{5, 4})
	a := c.reg.a.Read()

	r16mem.Write(a)
}

func (i *Instruction) ld_a_r16mem(c *CPU) {
	r16mem := c.getRegister16Mem(i.Opcode, []int{5, 4}).Read()

	c.reg.a.Write(r16mem)
}

func (i *Instruction) ld_imm16_sp(c *CPU) {
	imm16 := c.getImm16()
	sp := c.reg.sp.Read()

	c.mmu.WriteWord(imm16, sp)
}

func (i *Instruction) inc_r16(c *CPU) {
	r16 := c.getRegister16(i.Opcode, []int{5, 4})

	r16.Increment()
}

func (i *Instruction) dec_r16(c *CPU) {
	r16 := c.getRegister16(i.Opcode, []int{5, 4})

	r16.Decrement()
}

func (i *Instruction) add_hl_r16(c *CPU) {
	hl := c.reg.hl.Read()
	r16 := c.getRegister16(i.Opcode, []int{5, 4}).Read()

	res := hl + r16
	c.reg.hl.Write(res)

	c.reg.f.SetN(false)
	c.reg.f.SetH(IsHalfCarry16(hl, r16))
	c.reg.f.SetC(res < hl)
}

func (i *Instruction) inc_r8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{5, 4, 3})
	oldR8 := r8.Read()

	res := r8.Increment()

	c.reg.f.SetC(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(IsHalfCarry8(oldR8, 1))
}

func (i *Instruction) dec_r8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{5, 4, 3})
	oldR8 := r8.Read()

	res := r8.Decrement()

	c.reg.f.SetC(res == 0)
	c.reg.f.SetN(true)
	c.reg.f.SetH(IsHalfBorrow8(oldR8, 1))
}

func (i *Instruction) ld_r8_imm8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{5, 4, 3})
	imm8 := c.getImm8()

	r8.Write(imm8)
}

func (i *Instruction) rlca(c *CPU) {
	a := c.reg.a.Read()

	res, carry := RotateLeft(a)
	c.reg.a.Write(res)

	c.reg.f.SetZ(false)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) rrca(c *CPU) {
	a := c.reg.a.Read()

	res, carry := RotateRight(a)
	c.reg.a.Write(res)

	c.reg.f.SetZ(false)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) rla(c *CPU) {
	a := c.reg.a.Read()

	res, carry := RotateLeftThroughCarry(a, c.reg.f.C())
	c.reg.a.Write(res)

	c.reg.f.SetZ(false)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) rra(c *CPU) {
	a := c.reg.a.Read()

	res, carry := RotateRightThroughCarry(a, c.reg.f.C())
	c.reg.a.Write(res)

	c.reg.f.SetZ(false)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) daa(c *CPU) {
	a := c.reg.a.Read()
	adjustment := uint8(0)
	subtract := c.reg.f.N()

	if c.reg.f.H() || (!subtract && (a&0xF) > 0x9) {
		adjustment |= 0x06
	}
	if c.reg.f.C() || (!subtract && a > 0x99) {
		adjustment |= 0x60
		c.reg.f.SetC(true)
	}

	res := a
	if subtract {
		res -= adjustment
	} else {
		res += adjustment
	}
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetH(false)
}

func (i *Instruction) cpl(c *CPU) {
	a := c.reg.a.Read()

	c.reg.a.Write(^a)

	c.reg.f.SetN(true)
	c.reg.f.SetH(true)
}

func (i *Instruction) scf(c *CPU) {
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(true)
}

func (i *Instruction) ccf(c *CPU) {
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(!c.reg.f.C())
}

func (i *Instruction) jr_imm8(c *CPU) {
	imm8 := int8(c.getImm8())
	pc := c.reg.pc.Read()

	res := uint16(int(pc) + int(imm8))
	c.reg.pc.Write(res)
}

func (i *Instruction) jr_cond_imm8(c *CPU) {
	cond := c.getCond(i.Opcode, []int{4, 3})
	imm8 := int8(c.getImm8())
	pc := c.reg.pc.Read()

	if cond {
		res := uint16(int(pc) + int(imm8))
		c.reg.pc.Write(res)
	}
}

func (i *Instruction) stop(c *CPU) {}

// Block 1
func (i *Instruction) ld_r8_r8(c *CPU) {
	dst := c.getRegister8(i.Opcode, []int{5, 4, 3})
	src := c.getRegister8(i.Opcode, []int{2, 1, 0}).Read()

	dst.Write(src)
}

func (i *Instruction) halt(c *CPU) {
	c.halted = true
}

// Block 2
func (i *Instruction) add_a_r8(c *CPU) {
	a := c.reg.a.Read()
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0}).Read()

	res := a + r8
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(IsHalfCarry8(a, r8))
	c.reg.f.SetC(res < a)
}

func (i *Instruction) adc_a_r8(c *CPU) {
	a := c.reg.a.Read()
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0}).Read()
	carry := AsUint8(c.reg.f.C())

	res := a + r8 + carry
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(IsHalfCarry8(a, r8+carry))
	c.reg.f.SetC(res < a)
}

func (i *Instruction) sub_a_r8(c *CPU) {
	a := c.reg.a.Read()
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0}).Read()

	res := a - r8
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(true)
	c.reg.f.SetH(IsHalfBorrow8(a, r8))
	c.reg.f.SetC(r8 > a)
}

func (i *Instruction) sbc_a_r8(c *CPU) {
	a := c.reg.a.Read()
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0}).Read()
	carry := AsUint8(c.reg.f.C())

	res := a - r8 - carry
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(true)
	c.reg.f.SetH(IsHalfBorrow8(a, r8))
	c.reg.f.SetC(r8+carry > a)
}

func (i *Instruction) and_a_r8(c *CPU) {
	a := c.reg.a.Read()
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0}).Read()

	res := a & r8
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(true)
	c.reg.f.SetC(false)
}

func (i *Instruction) xor_a_r8(c *CPU) {
	a := c.reg.a.Read()
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0}).Read()

	res := a ^ r8
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(false)
}

func (i *Instruction) or_a_r8(c *CPU) {
	a := c.reg.a.Read()
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0}).Read()

	res := a | r8
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(false)
}

func (i *Instruction) cp_a_r8(c *CPU) {
	a := c.reg.a.Read()
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0}).Read()

	res := a - r8

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(true)
	c.reg.f.SetH(IsHalfBorrow8(a, r8))
	c.reg.f.SetC(r8 > a)
}

// Block 3
func (i *Instruction) add_a_imm8(c *CPU) {
	a := c.reg.a.Read()
	imm8 := c.getImm8()

	res := a + imm8
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(IsHalfCarry8(a, imm8))
	c.reg.f.SetC(res < a)
}

func (i *Instruction) adc_a_imm8(c *CPU) {
	a := c.reg.a.Read()
	imm8 := c.getImm8()
	carry := AsUint8(c.reg.f.C())

	res := a + imm8 + carry
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(IsHalfCarry8(a, imm8+carry))
	c.reg.f.SetC(res < a)
}

func (i *Instruction) sub_a_imm8(c *CPU) {
	a := c.reg.a.Read()
	imm8 := c.getImm8()

	res := a - imm8
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(true)
	c.reg.f.SetH(IsHalfBorrow8(a, imm8))
	c.reg.f.SetC(imm8 > a)
}

func (i *Instruction) sbc_a_imm8(c *CPU) {
	a := c.reg.a.Read()
	imm8 := c.getImm8()
	carry := AsUint8(c.reg.f.C())

	res := a - imm8 - carry
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(true)
	c.reg.f.SetH(IsHalfBorrow8(a, imm8+carry))
	c.reg.f.SetC((imm8 + carry) > a)
}

func (i *Instruction) and_a_imm8(c *CPU) {
	a := c.reg.a.Read()
	imm8 := c.getImm8()

	res := a & imm8
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(true)
	c.reg.f.SetC(false)
}

func (i *Instruction) xor_a_imm8(c *CPU) {
	a := c.reg.a.Read()
	imm8 := c.getImm8()

	res := a ^ imm8
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(false)
}

func (i *Instruction) or_a_imm8(c *CPU) {
	a := c.reg.a.Read()
	imm8 := c.getImm8()

	res := a | imm8
	c.reg.a.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(false)
}

func (i *Instruction) cp_a_imm8(c *CPU) {
	a := c.reg.a.Read()
	imm8 := c.getImm8()

	res := a - imm8

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(true)
	c.reg.f.SetH(IsHalfBorrow8(a, imm8))
	c.reg.f.SetC(imm8 > a)
}

func (i *Instruction) ret_cond(c *CPU) {
	cond := c.getCond(i.Opcode, []int{4, 3})
	if cond {
		spMem := c.mmu.ReadWord(c.reg.sp.Read())

		c.reg.pc.Write(spMem)
		c.reg.sp.Write(c.reg.sp.Read() + 2)
	}
}

func (i *Instruction) ret(c *CPU) {
	spMem := c.mmu.ReadWord(c.reg.sp.Read())

	c.reg.pc.Write(spMem)
	c.reg.sp.Write(c.reg.sp.Read() + 2)
}

func (i *Instruction) reti(c *CPU) {
	spMem := c.mmu.ReadWord(c.reg.sp.Read())

	c.reg.pc.Write(spMem)
	c.reg.sp.Write(c.reg.sp.Read() + 2)
	c.IME = true
}

func (i *Instruction) jp_cond_imm16(c *CPU) {
	cond := c.getCond(i.Opcode, []int{4, 3})
	imm16 := c.getImm16()

	if cond {
		c.reg.pc.Write(imm16)
	}
}

func (i *Instruction) jp_imm16(c *CPU) {
	imm16 := c.getImm16()

	c.reg.pc.Write(imm16)
}

func (i *Instruction) jp_hl(c *CPU) {
	c.reg.pc.Write(c.reg.hl.Read())
}

func (i *Instruction) call_cond_imm16(c *CPU) {
	cond := c.getCond(i.Opcode, []int{4, 3})
	imm16 := c.getImm16()
	sp := c.reg.sp.Read()
	pc := c.reg.pc.Read()

	if cond {
		c.mmu.WriteWord(sp-2, pc)
		c.reg.sp.Write(sp - 2)
		c.reg.pc.Write(imm16)
	}
}

func (i *Instruction) call_imm16(c *CPU) {
	imm16 := c.getImm16()
	sp := c.reg.sp.Read()
	pc := c.reg.pc.Read()

	c.mmu.WriteWord(sp-2, pc)
	c.reg.sp.Write(sp - 2)
	c.reg.pc.Write(imm16)
}

func (i *Instruction) rst_tgt3(c *CPU) {
	tgt := ExtractBits(i.Opcode, []int{5, 4, 3})
	sp := c.reg.sp.Read()
	pc := c.reg.pc.Read()

	c.mmu.WriteWord(sp-2, pc)
	c.reg.sp.Write(sp - 2)
	c.reg.pc.Write(uint16(tgt))
}

func (i *Instruction) pop_r16stk(c *CPU) {
	r16Stk := c.getRegister16Stk(i.Opcode, []int{5, 4})
	sp := c.reg.sp.Read()

	spMem := c.mmu.ReadWord(sp)
	c.reg.sp.Write(sp + 2)
	r16Stk.Write(spMem)
}

func (i *Instruction) push_r16stk(c *CPU) {
	r16Stk := c.getRegister16Stk(i.Opcode, []int{5, 4})
	sp := c.reg.sp.Read()

	c.mmu.WriteWord(sp-2, r16Stk.Read())
	c.reg.sp.Write(sp - 2)
}

func (i *Instruction) ldh_c_a(c *CPU) {
	cMem := c.byteAt(0xFF00 + uint16(c.reg.c.Read()))
	a := c.reg.a.Read()

	cMem.Write(a)
}

func (i *Instruction) ldh_imm8_a(c *CPU) {
	imm8 := c.getImm8()
	imm8Mem := c.byteAt(0xFF00 + uint16(imm8))
	a := c.reg.a.Read()

	imm8Mem.Write(a)
}

func (i *Instruction) ld_imm16_a(c *CPU) {
	imm16 := c.getImm16()
	imm16Mem := c.byteAt(imm16)
	a := c.reg.a.Read()

	imm16Mem.Write(a)
}

func (i *Instruction) ldh_a_c(c *CPU) {
	cMem := c.byteAt(0xFF00 + uint16(c.reg.c.Read()))

	c.reg.a.Write(cMem.Read())
}

func (i *Instruction) ldh_a_imm8(c *CPU) {
	imm8 := c.getImm8()
	imm8Mem := c.byteAt(0xFF00 + uint16(imm8))

	c.reg.a.Write(imm8Mem.Read())
}

func (i *Instruction) ld_a_imm16(c *CPU) {
	imm16 := c.getImm16()
	imm16Mem := c.byteAt(imm16)

	c.reg.a.Write(imm16Mem.Read())
}

func (i *Instruction) add_sp_imm8(c *CPU) {
	imm8 := c.getImm8()
	sp := c.reg.sp.Read()

	res := uint16(int(sp) + int(imm8))
	c.reg.sp.Write(res)

	c.reg.f.SetZ(false)
	c.reg.f.SetN(false)
	c.reg.f.SetH(IsHalfCarry8(uint8(sp&0xFF), imm8))
	c.reg.f.SetC(uint8(res&0xFF) < uint8(sp&0xFF))
}

// ld hl, sp + imm8
func (i *Instruction) ld_hl_sp_plus_imm8(c *CPU) {
	imm8 := c.getImm8()
	sp := c.reg.sp.Read()

	res := uint16(int(sp) + int(imm8))
	c.reg.hl.Write(res)

	c.reg.f.SetZ(false)
	c.reg.f.SetN(false)
	c.reg.f.SetH(IsHalfCarry8(uint8(sp&0xFF), imm8))
	c.reg.f.SetC(uint8(res&0xFF) < uint8(sp&0xFF))
}

func (i *Instruction) ld_sp_hl(c *CPU) {
	hl := c.reg.hl.Read()

	c.reg.sp.Write(hl)
}

func (i *Instruction) di(c *CPU) {
	c.IME = false
}

func (i *Instruction) ei(c *CPU) {
	c.IME = true
}

// 0xCB Prefixed instructions
func (i *Instruction) rlc_r8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	res, carry := RotateLeft(r8.Read())
	r8.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) rrc_r8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	res, carry := RotateRight(r8.Read())
	r8.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) rl_r8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	res, carry := RotateLeftThroughCarry(r8.Read(), c.reg.f.C())
	r8.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) rr_r8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	res, carry := RotateRightThroughCarry(r8.Read(), c.reg.f.C())
	r8.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) sla_r8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	res, carry := ShiftLeftArithmetic(r8.Read())
	r8.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) sra_r8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	res, carry := ShiftRightArithmetic(r8.Read())
	r8.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) swap_r8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	res := Swap(r8.Read())
	r8.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(false)
}

func (i *Instruction) srl_r8(c *CPU) {
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	res, carry := ShiftRightLogic(r8.Read())
	r8.Write(res)

	c.reg.f.SetZ(res == 0)
	c.reg.f.SetN(false)
	c.reg.f.SetH(false)
	c.reg.f.SetC(carry)
}

func (i *Instruction) bit_b3_r8(c *CPU) {
	b3 := ExtractBits(i.Opcode, []int{5, 4, 3})
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	c.reg.f.SetZ(!IsBitSet(r8.Read(), b3))
	c.reg.f.SetN(false)
	c.reg.f.SetH(true)
}

func (i *Instruction) res_b3_r8(c *CPU) {
	b3 := ExtractBits(i.Opcode, []int{5, 4, 3})
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	res := ResetBit(r8.Read(), b3)
	r8.Write(res)
}

func (i *Instruction) set_b3_r8(c *CPU) {
	b3 := ExtractBits(i.Opcode, []int{5, 4, 3})
	r8 := c.getRegister8(i.Opcode, []int{2, 1, 0})

	res := SetBit(r8.Read(), b3)
	r8.Write(res)
}
