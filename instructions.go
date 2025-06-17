package main

import "log"

type Instruction struct {
	Opcode   byte
	Mnemonic string
	Length   uint8
	Cycles   Cycles
	Execute  func(c *CPU)
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

var INVALID_INSTRUCTION Instruction = Instruction{0xFF, "INVALID", 0, NewCycles(1, 1), (*CPU).invalid_instruction}

var INSTRUCTIONS []Instruction = []Instruction{
	{0x00, "NOP", 1, NewCycles(4, 4), (*CPU).nop},
	{0x01, "LD BC, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x02, "LD (BC), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x03, "INC BC", 1, NewCycles(8, 8), (*CPU).inc_r16},
	{0x04, "INC B", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x05, "DEC B", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x06, "LD B, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x07, "RLCA", 1, NewCycles(4, 4), (*CPU).rlca},
	{0x08, "LD (u16), SP", 3, NewCycles(20, 20), (*CPU).ld_imm16_sp},
	{0x09, "ADD HL, BC", 1, NewCycles(8, 8), (*CPU).add_hl_r16},
	{0x0a, "LD A, (BC)", 1, NewCycles(8, 8), (*CPU).ld_a_r16mem},
	{0x0b, "DEC BC", 1, NewCycles(8, 8), (*CPU).dec_r16},
	{0x0c, "INC C", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x0d, "DEC C", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x0e, "LD C, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x0f, "RRCA", 1, NewCycles(4, 4), (*CPU).rrca},
	{0x10, "STOP 0", 1, NewCycles(4, 4), (*CPU).stop},
	{0x11, "LD DE, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x12, "LD (DE), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x13, "INC DE", 1, NewCycles(8, 8), (*CPU).inc_r16},
	{0x14, "INC D", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x15, "DEC D", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x16, "LD D, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x17, "RLA", 1, NewCycles(4, 4), (*CPU).rla},
	{0x18, "JR i8", 2, NewCycles(12, 12), (*CPU).jr_imm8},
	{0x19, "ADD HL, DE", 1, NewCycles(8, 8), (*CPU).add_hl_r16},
	{0x1a, "LD A, (DE)", 1, NewCycles(8, 8), (*CPU).ld_a_r16mem},
	{0x1b, "DEC DE", 1, NewCycles(8, 8), (*CPU).dec_r16},
	{0x1c, "INC E", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x1d, "DEC E", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x1e, "LD E, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x1f, "RRA", 1, NewCycles(4, 4), (*CPU).rra},
	{0x20, "JR NZ, i8", 2, NewCycles(12, 8), (*CPU).jr_cond_imm8},
	{0x21, "LD HL, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x22, "LD (HL+), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x23, "INC HL", 1, NewCycles(8, 8), (*CPU).inc_r16},
	{0x24, "INC H", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x25, "DEC H", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x26, "LD H, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x27, "DAA", 1, NewCycles(4, 4), (*CPU).daa},
	{0x28, "JR Z, i8", 2, NewCycles(12, 8), (*CPU).jr_cond_imm8},
	{0x29, "ADD HL, HL", 1, NewCycles(8, 8), (*CPU).add_hl_r16},
	{0x2a, "LD A, (HL+)", 1, NewCycles(8, 8), (*CPU).ld_a_r16mem},
	{0x2b, "DEC HL", 1, NewCycles(8, 8), (*CPU).dec_r16},
	{0x2c, "INC L", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x2d, "DEC L", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x2e, "LD L, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x2f, "CPL", 1, NewCycles(4, 4), (*CPU).cpl},
	{0x30, "JR NC, i8", 2, NewCycles(12, 8), (*CPU).jr_cond_imm8},
	{0x31, "LD SP, u16", 3, NewCycles(12, 12), (*CPU).ld_r16_imm16},
	{0x32, "LD (HL-), A", 1, NewCycles(8, 8), (*CPU).ld_r16mem_a},
	{0x33, "INC SP", 1, NewCycles(8, 8), (*CPU).inc_r16},
	{0x34, "INC (HL)", 1, NewCycles(12, 12), (*CPU).inc_r8},
	{0x35, "DEC (HL)", 1, NewCycles(12, 12), (*CPU).dec_r8},
	{0x36, "LD (HL), u8", 2, NewCycles(12, 12), (*CPU).ld_r8_imm8},
	{0x37, "SCF", 1, NewCycles(4, 4), (*CPU).scf},
	{0x38, "JR C, i8", 2, NewCycles(12, 8), (*CPU).jr_cond_imm8},
	{0x39, "ADD HL, SP", 1, NewCycles(8, 8), (*CPU).add_hl_r16},
	{0x3a, "LD A, (HL-)", 1, NewCycles(8, 8), (*CPU).ld_a_r16mem},
	{0x3b, "DEC SP", 1, NewCycles(8, 8), (*CPU).dec_r16},
	{0x3c, "INC A", 1, NewCycles(4, 4), (*CPU).inc_r8},
	{0x3d, "DEC A", 1, NewCycles(4, 4), (*CPU).dec_r8},
	{0x3e, "LD A, u8", 2, NewCycles(8, 8), (*CPU).ld_r8_imm8},
	{0x3f, "CCF", 1, NewCycles(4, 4), (*CPU).ccf},
	{0x40, "LD B, B", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x41, "LD B, C", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x42, "LD B, D", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x43, "LD B, E", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x44, "LD B, H", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x45, "LD B, L", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x46, "LD B, (HL)", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x47, "LD B, A", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x48, "LD C, B", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x49, "LD C, C", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x4a, "LD C, D", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x4b, "LD C, E", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x4c, "LD C, H", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x4d, "LD C, L", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x4e, "LD C, (HL)", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x4f, "LD C, A", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x50, "LD D, B", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x51, "LD D, C", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x52, "LD D, D", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x53, "LD D, E", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x54, "LD D, H", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x55, "LD D, L", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x56, "LD D, (HL)", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x57, "LD D, A", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x58, "LD E, B", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x59, "LD E, C", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x5a, "LD E, D", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x5b, "LD E, E", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x5c, "LD E, H", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x5d, "LD E, L", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x5e, "LD E, (HL)", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x5f, "LD E, A", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x60, "LD H, B", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x61, "LD H, C", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x62, "LD H, D", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x63, "LD H, E", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x64, "LD H, H", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x65, "LD H, L", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x66, "LD H, (HL)", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x67, "LD H, A", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x68, "LD L, B", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x69, "LD L, C", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x6a, "LD L, D", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x6b, "LD L, E", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x6c, "LD L, H", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x6d, "LD L, L", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x6e, "LD L, (HL)", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x6f, "LD L, A", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x70, "LD (HL), B", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x71, "LD (HL), C", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x72, "LD (HL), D", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x73, "LD (HL), E", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x74, "LD (HL), H", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x75, "LD (HL), L", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x76, "HALT", 1, NewCycles(4, 4), (*CPU).halt},
	{0x77, "LD (HL), A", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x78, "LD A, B", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x79, "LD A, C", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x7a, "LD A, D", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x7b, "LD A, E", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x7c, "LD A, H", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x7d, "LD A, L", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x7e, "LD A, (HL)", 1, NewCycles(8, 8), (*CPU).ld_r8_r8},
	{0x7f, "LD A, A", 1, NewCycles(4, 4), (*CPU).ld_r8_r8},
	{0x80, "ADD A, B", 1, NewCycles(4, 4), (*CPU).add_a_r8},
	{0x81, "ADD A, C", 1, NewCycles(4, 4), (*CPU).add_a_r8},
	{0x82, "ADD A, D", 1, NewCycles(4, 4), (*CPU).add_a_r8},
	{0x83, "ADD A, E", 1, NewCycles(4, 4), (*CPU).add_a_r8},
	{0x84, "ADD A, H", 1, NewCycles(4, 4), (*CPU).add_a_r8},
	{0x85, "ADD A, L", 1, NewCycles(4, 4), (*CPU).add_a_r8},
	{0x86, "ADD A, (HL)", 1, NewCycles(8, 8), (*CPU).add_a_r8},
	{0x87, "ADD A, A", 1, NewCycles(4, 4), (*CPU).add_a_r8},
	{0x88, "ADC A, B", 1, NewCycles(4, 4), (*CPU).adc_a_r8},
	{0x89, "ADC A, C", 1, NewCycles(4, 4), (*CPU).adc_a_r8},
	{0x8a, "ADC A, D", 1, NewCycles(4, 4), (*CPU).adc_a_r8},
	{0x8b, "ADC A, E", 1, NewCycles(4, 4), (*CPU).adc_a_r8},
	{0x8c, "ADC A, H", 1, NewCycles(4, 4), (*CPU).adc_a_r8},
	{0x8d, "ADC A, L", 1, NewCycles(4, 4), (*CPU).adc_a_r8},
	{0x8e, "ADC A, (HL)", 1, NewCycles(8, 8), (*CPU).adc_a_r8},
	{0x8f, "ADC A, A", 1, NewCycles(4, 4), (*CPU).adc_a_r8},
	{0x90, "SUB B", 1, NewCycles(4, 4), (*CPU).sub_a_r8},
	{0x91, "SUB C", 1, NewCycles(4, 4), (*CPU).sub_a_r8},
	{0x92, "SUB D", 1, NewCycles(4, 4), (*CPU).sub_a_r8},
	{0x93, "SUB E", 1, NewCycles(4, 4), (*CPU).sub_a_r8},
	{0x94, "SUB H", 1, NewCycles(4, 4), (*CPU).sub_a_r8},
	{0x95, "SUB L", 1, NewCycles(4, 4), (*CPU).sub_a_r8},
	{0x96, "SUB (HL)", 1, NewCycles(8, 8), (*CPU).sub_a_r8},
	{0x97, "SUB A", 1, NewCycles(4, 4), (*CPU).sub_a_r8},
	{0x98, "SBC A, B", 1, NewCycles(4, 4), (*CPU).sbc_a_r8},
	{0x99, "SBC A, C", 1, NewCycles(4, 4), (*CPU).sbc_a_r8},
	{0x9a, "SBC A, D", 1, NewCycles(4, 4), (*CPU).sbc_a_r8},
	{0x9b, "SBC A, E", 1, NewCycles(4, 4), (*CPU).sbc_a_r8},
	{0x9c, "SBC A, H", 1, NewCycles(4, 4), (*CPU).sbc_a_r8},
	{0x9d, "SBC A, L", 1, NewCycles(4, 4), (*CPU).sbc_a_r8},
	{0x9e, "SBC A, (HL)", 1, NewCycles(8, 8), (*CPU).sbc_a_r8},
	{0x9f, "SBC A, A", 1, NewCycles(4, 4), (*CPU).sbc_a_r8},
	{0xa0, "AND B", 1, NewCycles(4, 4), (*CPU).and_a_r8},
	{0xa1, "AND C", 1, NewCycles(4, 4), (*CPU).and_a_r8},
	{0xa2, "AND D", 1, NewCycles(4, 4), (*CPU).and_a_r8},
	{0xa3, "AND E", 1, NewCycles(4, 4), (*CPU).and_a_r8},
	{0xa4, "AND H", 1, NewCycles(4, 4), (*CPU).and_a_r8},
	{0xa5, "AND L", 1, NewCycles(4, 4), (*CPU).and_a_r8},
	{0xa6, "AND (HL)", 1, NewCycles(8, 8), (*CPU).and_a_r8},
	{0xa7, "AND A", 1, NewCycles(4, 4), (*CPU).and_a_r8},
	{0xa8, "XOR B", 1, NewCycles(4, 4), (*CPU).xor_a_r8},
	{0xa9, "XOR C", 1, NewCycles(4, 4), (*CPU).xor_a_r8},
	{0xaa, "XOR D", 1, NewCycles(4, 4), (*CPU).xor_a_r8},
	{0xab, "XOR E", 1, NewCycles(4, 4), (*CPU).xor_a_r8},
	{0xac, "XOR H", 1, NewCycles(4, 4), (*CPU).xor_a_r8},
	{0xad, "XOR L", 1, NewCycles(4, 4), (*CPU).xor_a_r8},
	{0xae, "XOR (HL)", 1, NewCycles(8, 8), (*CPU).xor_a_r8},
	{0xaf, "XOR A", 1, NewCycles(4, 4), (*CPU).xor_a_r8},
	{0xb0, "OR B", 1, NewCycles(4, 4), (*CPU).or_a_r8},
	{0xb1, "OR C", 1, NewCycles(4, 4), (*CPU).or_a_r8},
	{0xb2, "OR D", 1, NewCycles(4, 4), (*CPU).or_a_r8},
	{0xb3, "OR E", 1, NewCycles(4, 4), (*CPU).or_a_r8},
	{0xb4, "OR H", 1, NewCycles(4, 4), (*CPU).or_a_r8},
	{0xb5, "OR L", 1, NewCycles(4, 4), (*CPU).or_a_r8},
	{0xb6, "OR (HL)", 1, NewCycles(8, 8), (*CPU).or_a_r8},
	{0xb7, "OR A", 1, NewCycles(4, 4), (*CPU).or_a_r8},
	{0xb8, "CP B", 1, NewCycles(4, 4), (*CPU).cp_a_r8},
	{0xb9, "CP C", 1, NewCycles(4, 4), (*CPU).cp_a_r8},
	{0xba, "CP D", 1, NewCycles(4, 4), (*CPU).cp_a_r8},
	{0xbb, "CP E", 1, NewCycles(4, 4), (*CPU).cp_a_r8},
	{0xbc, "CP H", 1, NewCycles(4, 4), (*CPU).cp_a_r8},
	{0xbd, "CP L", 1, NewCycles(4, 4), (*CPU).cp_a_r8},
	{0xbe, "CP (HL)", 1, NewCycles(8, 8), (*CPU).cp_a_r8},
	{0xbf, "CP A", 1, NewCycles(4, 4), (*CPU).cp_a_r8},
	{0xc0, "RET NZ", 1, NewCycles(20, 8), (*CPU).ret_cond},
	{0xc1, "POP BC", 1, NewCycles(12, 12), (*CPU).pop_r16stk},
	{0xc2, "JP NZ, u16", 3, NewCycles(16, 12), (*CPU).jp_cond_imm16},
	{0xc3, "JP u16", 3, NewCycles(16, 16), (*CPU).jp_imm16},
	{0xc4, "CALL NZ, u16", 3, NewCycles(24, 12), (*CPU).call_cond_imm16},
	{0xc5, "PUSH BC", 1, NewCycles(16, 16), (*CPU).push_r16stk},
	{0xc6, "ADD A, u8", 2, NewCycles(8, 8), (*CPU).add_a_imm8},
	{0xc7, "RST 00H", 1, NewCycles(16, 16), (*CPU).rst_tgt3},
	{0xc8, "RET Z", 1, NewCycles(20, 8), (*CPU).ret_cond},
	{0xc9, "RET", 1, NewCycles(16, 16), (*CPU).ret},
	{0xca, "JP Z, u16", 3, NewCycles(16, 12), (*CPU).jp_cond_imm16},
	INVALID_INSTRUCTION,
	{0xcc, "CALL Z, u16", 3, NewCycles(24, 12), (*CPU).call_cond_imm16},
	{0xcd, "CALL u16", 3, NewCycles(24, 24), (*CPU).call_imm16},
	{0xce, "ADC A, u8", 2, NewCycles(8, 8), (*CPU).adc_a_imm8},
	{0xcf, "RST 08H", 1, NewCycles(16, 16), (*CPU).rst_tgt3},
	{0xd0, "RET NC", 1, NewCycles(20, 8), (*CPU).ret_cond},
	{0xd1, "POP DE", 1, NewCycles(12, 12), (*CPU).pop_r16stk},
	{0xd2, "JP NC, u16", 3, NewCycles(16, 12), (*CPU).jp_cond_imm16},
	INVALID_INSTRUCTION,
	{0xd4, "CALL NC, u16", 3, NewCycles(24, 12), (*CPU).call_cond_imm16},
	{0xd5, "PUSH DE", 1, NewCycles(16, 16), (*CPU).push_r16stk},
	{0xd6, "SUB u8", 2, NewCycles(8, 8), (*CPU).sub_a_r8},
	{0xd7, "RST 10H", 1, NewCycles(16, 16), (*CPU).rst_tgt3},
	{0xd8, "RET C", 1, NewCycles(20, 8), (*CPU).ret_cond},
	{0xd9, "RETI", 1, NewCycles(16, 16), (*CPU).reti},
	{0xda, "JP C, u16", 3, NewCycles(16, 12), (*CPU).jp_cond_imm16},
	INVALID_INSTRUCTION,
	{0xdc, "CALL C, u16", 3, NewCycles(24, 12), (*CPU).call_cond_imm16},
	INVALID_INSTRUCTION,
	{0xde, "SBC A, u8", 2, NewCycles(8, 8), (*CPU).sbc_a_imm8},
	{0xdf, "RST 18H", 1, NewCycles(16, 16), (*CPU).rst_tgt3},
	{0xe0, "LDH (u8), A", 2, NewCycles(12, 12), (*CPU).ldh_imm8_a},
	{0xe1, "POP HL", 1, NewCycles(12, 12), (*CPU).pop_r16stk},
	{0xe2, "LD (C), A", 1, NewCycles(8, 8), (*CPU).ldh_c_a},
	INVALID_INSTRUCTION,
	INVALID_INSTRUCTION,
	{0xe5, "PUSH HL", 1, NewCycles(16, 16), (*CPU).push_r16stk},
	{0xe6, "AND u8", 2, NewCycles(8, 8), (*CPU).and_a_r8},
	{0xe7, "RST 20H", 1, NewCycles(16, 16), (*CPU).rst_tgt3},
	{0xe8, "ADD SP, i8", 2, NewCycles(16, 16), (*CPU).add_sp_imm8},
	{0xe9, "JP (HL)", 1, NewCycles(4, 4), (*CPU).jp_imm16},
	{0xea, "LD (u16), A", 3, NewCycles(16, 16), (*CPU).ld_imm16_a},
	INVALID_INSTRUCTION,
	INVALID_INSTRUCTION,
	INVALID_INSTRUCTION,
	{0xee, "XOR u8", 2, NewCycles(8, 8), (*CPU).xor_a_r8},
	{0xef, "RST 28H", 1, NewCycles(16, 16), (*CPU).rst_tgt3},
	{0xf0, "LDH A, (u8)", 2, NewCycles(12, 12), (*CPU).ldh_a_imm8},
	{0xf1, "POP AF", 1, NewCycles(12, 12), (*CPU).pop_r16stk},
	{0xf2, "LD A, (C)", 1, NewCycles(8, 8), (*CPU).ldh_a_c},
	{0xf3, "DI", 1, NewCycles(4, 4), (*CPU).di},
	INVALID_INSTRUCTION,
	{0xf5, "PUSH AF", 1, NewCycles(16, 16), (*CPU).push_r16stk},
	{0xf6, "OR u8", 2, NewCycles(8, 8), (*CPU).or_a_r8},
	{0xf7, "RST 30H", 1, NewCycles(16, 16), (*CPU).rst_tgt3},
	{0xf8, "LD HL, SP+r8", 2, NewCycles(12, 12), (*CPU).ld_hl_sp_plus_imm8},
	{0xf9, "LD SP, HL", 1, NewCycles(8, 8), (*CPU).ld_sp_hl},
	{0xfa, "LD A, (u16)", 3, NewCycles(16, 16), (*CPU).ld_a_imm16},
	{0xfb, "EI", 1, NewCycles(4, 4), (*CPU).ei},
	{0xfe, "CP u8", 2, NewCycles(8, 8), (*CPU).cp_a_r8},
	{0xff, "RST 38H", 1, NewCycles(16, 16), (*CPU).rst_tgt3},
}

var CB_INSTRUCTIONS []Instruction = []Instruction{
	{0x00, "RLC B", 2, NewCycles(8, 8), (*CPU).rlc_r8},
	{0x01, "RLC C", 2, NewCycles(8, 8), (*CPU).rlc_r8},
	{0x02, "RLC D", 2, NewCycles(8, 8), (*CPU).rlc_r8},
	{0x03, "RLC E", 2, NewCycles(8, 8), (*CPU).rlc_r8},
	{0x04, "RLC H", 2, NewCycles(8, 8), (*CPU).rlc_r8},
	{0x05, "RLC L", 2, NewCycles(8, 8), (*CPU).rlc_r8},
	{0x06, "RLC (HL)", 2, NewCycles(16, 16), (*CPU).rlc_r8},
	{0x07, "RLC A", 2, NewCycles(8, 8), (*CPU).rlc_r8},
	{0x08, "RRC B", 2, NewCycles(8, 8), (*CPU).rrc_r8},
	{0x09, "RRC C", 2, NewCycles(8, 8), (*CPU).rrc_r8},
	{0x0a, "RRC D", 2, NewCycles(8, 8), (*CPU).rrc_r8},
	{0x0b, "RRC E", 2, NewCycles(8, 8), (*CPU).rrc_r8},
	{0x0c, "RRC H", 2, NewCycles(8, 8), (*CPU).rrc_r8},
	{0x0d, "RRC L", 2, NewCycles(8, 8), (*CPU).rrc_r8},
	{0x0e, "RRC (HL)", 2, NewCycles(16, 16), (*CPU).rrc_r8},
	{0x0f, "RRC A", 2, NewCycles(8, 8), (*CPU).rrc_r8},
	{0x10, "RL B", 2, NewCycles(8, 8), (*CPU).rl_r8},
	{0x11, "RL C", 2, NewCycles(8, 8), (*CPU).rl_r8},
	{0x12, "RL D", 2, NewCycles(8, 8), (*CPU).rl_r8},
	{0x13, "RL E", 2, NewCycles(8, 8), (*CPU).rl_r8},
	{0x14, "RL H", 2, NewCycles(8, 8), (*CPU).rl_r8},
	{0x15, "RL L", 2, NewCycles(8, 8), (*CPU).rl_r8},
	{0x16, "RL (HL)", 2, NewCycles(16, 16), (*CPU).rl_r8},
	{0x17, "RL A", 2, NewCycles(8, 8), (*CPU).rl_r8},
	{0x18, "RR B", 2, NewCycles(8, 8), (*CPU).rr_r8},
	{0x19, "RR C", 2, NewCycles(8, 8), (*CPU).rr_r8},
	{0x1a, "RR D", 2, NewCycles(8, 8), (*CPU).rr_r8},
	{0x1b, "RR E", 2, NewCycles(8, 8), (*CPU).rr_r8},
	{0x1c, "RR H", 2, NewCycles(8, 8), (*CPU).rr_r8},
	{0x1d, "RR L", 2, NewCycles(8, 8), (*CPU).rr_r8},
	{0x1e, "RR (HL)", 2, NewCycles(16, 16), (*CPU).rr_r8},
	{0x1f, "RR A", 2, NewCycles(8, 8), (*CPU).rr_r8},
	{0x20, "SLA B", 2, NewCycles(8, 8), (*CPU).sla_r8},
	{0x21, "SLA C", 2, NewCycles(8, 8), (*CPU).sla_r8},
	{0x22, "SLA D", 2, NewCycles(8, 8), (*CPU).sla_r8},
	{0x23, "SLA E", 2, NewCycles(8, 8), (*CPU).sla_r8},
	{0x24, "SLA H", 2, NewCycles(8, 8), (*CPU).sla_r8},
	{0x25, "SLA L", 2, NewCycles(8, 8), (*CPU).sla_r8},
	{0x26, "SLA (HL)", 2, NewCycles(16, 16), (*CPU).sla_r8},
	{0x27, "SLA A", 2, NewCycles(8, 8), (*CPU).sla_r8},
	{0x28, "SRA B", 2, NewCycles(8, 8), (*CPU).sra_r8},
	{0x29, "SRA C", 2, NewCycles(8, 8), (*CPU).sra_r8},
	{0x2a, "SRA D", 2, NewCycles(8, 8), (*CPU).sra_r8},
	{0x2b, "SRA E", 2, NewCycles(8, 8), (*CPU).sra_r8},
	{0x2c, "SRA H", 2, NewCycles(8, 8), (*CPU).sra_r8},
	{0x2d, "SRA L", 2, NewCycles(8, 8), (*CPU).sra_r8},
	{0x2e, "SRA (HL)", 2, NewCycles(16, 16), (*CPU).sra_r8},
	{0x2f, "SRA A", 2, NewCycles(8, 8), (*CPU).sra_r8},
	{0x30, "SWAP B", 2, NewCycles(8, 8), (*CPU).swap_r8},
	{0x31, "SWAP C", 2, NewCycles(8, 8), (*CPU).swap_r8},
	{0x32, "SWAP D", 2, NewCycles(8, 8), (*CPU).swap_r8},
	{0x33, "SWAP E", 2, NewCycles(8, 8), (*CPU).swap_r8},
	{0x34, "SWAP H", 2, NewCycles(8, 8), (*CPU).swap_r8},
	{0x35, "SWAP L", 2, NewCycles(8, 8), (*CPU).swap_r8},
	{0x36, "SWAP (HL)", 2, NewCycles(16, 16), (*CPU).swap_r8},
	{0x37, "SWAP A", 2, NewCycles(8, 8), (*CPU).swap_r8},
	{0x38, "SRL B", 2, NewCycles(8, 8), (*CPU).srl_r8},
	{0x39, "SRL C", 2, NewCycles(8, 8), (*CPU).srl_r8},
	{0x3a, "SRL D", 2, NewCycles(8, 8), (*CPU).srl_r8},
	{0x3b, "SRL E", 2, NewCycles(8, 8), (*CPU).srl_r8},
	{0x3c, "SRL H", 2, NewCycles(8, 8), (*CPU).srl_r8},
	{0x3d, "SRL L", 2, NewCycles(8, 8), (*CPU).srl_r8},
	{0x3e, "SRL (HL)", 2, NewCycles(16, 16), (*CPU).srl_r8},
	{0x3f, "SRL A", 2, NewCycles(8, 8), (*CPU).srl_r8},
	{0x40, "BIT 0, B", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x41, "BIT 0, C", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x42, "BIT 0, D", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x43, "BIT 0, E", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x44, "BIT 0, H", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x45, "BIT 0, L", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x46, "BIT 0, (HL)", 2, NewCycles(16, 16), (*CPU).bit_b3_r8},
	{0x47, "BIT 0, A", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x48, "BIT 1, B", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x49, "BIT 1, C", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x4a, "BIT 1, D", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x4b, "BIT 1, E", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x4c, "BIT 1, H", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x4d, "BIT 1, L", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x4e, "BIT 1, (HL)", 2, NewCycles(16, 16), (*CPU).bit_b3_r8},
	{0x4f, "BIT 1, A", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x50, "BIT 2, B", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x51, "BIT 2, C", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x52, "BIT 2, D", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x53, "BIT 2, E", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x54, "BIT 2, H", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x55, "BIT 2, L", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x56, "BIT 2, (HL)", 2, NewCycles(16, 16), (*CPU).bit_b3_r8},
	{0x57, "BIT 2, A", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x58, "BIT 3, B", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x59, "BIT 3, C", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x5a, "BIT 3, D", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x5b, "BIT 3, E", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x5c, "BIT 3, H", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x5d, "BIT 3, L", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x5e, "BIT 3, (HL)", 2, NewCycles(16, 16), (*CPU).bit_b3_r8},
	{0x5f, "BIT 3, A", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x60, "BIT 4, B", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x61, "BIT 4, C", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x62, "BIT 4, D", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x63, "BIT 4, E", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x64, "BIT 4, H", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x65, "BIT 4, L", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x66, "BIT 4, (HL)", 2, NewCycles(16, 16), (*CPU).bit_b3_r8},
	{0x67, "BIT 4, A", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x68, "BIT 5, B", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x69, "BIT 5, C", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x6a, "BIT 5, D", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x6b, "BIT 5, E", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x6c, "BIT 5, H", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x6d, "BIT 5, L", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x6e, "BIT 5, (HL)", 2, NewCycles(16, 16), (*CPU).bit_b3_r8},
	{0x6f, "BIT 5, A", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x70, "BIT 6, B", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x71, "BIT 6, C", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x72, "BIT 6, D", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x73, "BIT 6, E", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x74, "BIT 6, H", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x75, "BIT 6, L", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x76, "BIT 6, (HL)", 2, NewCycles(16, 16), (*CPU).bit_b3_r8},
	{0x77, "BIT 6, A", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x78, "BIT 7, B", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x79, "BIT 7, C", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x7a, "BIT 7, D", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x7b, "BIT 7, E", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x7c, "BIT 7, H", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x7d, "BIT 7, L", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x7e, "BIT 7, (HL)", 2, NewCycles(16, 16), (*CPU).bit_b3_r8},
	{0x7f, "BIT 7, A", 2, NewCycles(8, 8), (*CPU).bit_b3_r8},
	{0x80, "RES 0, B", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x81, "RES 0, C", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x82, "RES 0, D", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x83, "RES 0, E", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x84, "RES 0, H", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x85, "RES 0, L", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x86, "RES 0, (HL)", 2, NewCycles(16, 16), (*CPU).res_b3_r8},
	{0x87, "RES 0, A", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x88, "RES 1, B", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x89, "RES 1, C", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x8a, "RES 1, D", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x8b, "RES 1, E", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x8c, "RES 1, H", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x8d, "RES 1, L", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x8e, "RES 1, (HL)", 2, NewCycles(16, 16), (*CPU).res_b3_r8},
	{0x8f, "RES 1, A", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x90, "RES 2, B", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x91, "RES 2, C", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x92, "RES 2, D", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x93, "RES 2, E", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x94, "RES 2, H", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x95, "RES 2, L", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x96, "RES 2, (HL)", 2, NewCycles(16, 16), (*CPU).res_b3_r8},
	{0x97, "RES 2, A", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x98, "RES 3, B", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x99, "RES 3, C", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x9a, "RES 3, D", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x9b, "RES 3, E", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x9c, "RES 3, H", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x9d, "RES 3, L", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0x9e, "RES 3, (HL)", 2, NewCycles(16, 16), (*CPU).res_b3_r8},
	{0x9f, "RES 3, A", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xa0, "RES 4, B", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xa1, "RES 4, C", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xa2, "RES 4, D", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xa3, "RES 4, E", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xa4, "RES 4, H", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xa5, "RES 4, L", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xa6, "RES 4, (HL)", 2, NewCycles(16, 16), (*CPU).res_b3_r8},
	{0xa7, "RES 4, A", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xa8, "RES 5, B", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xa9, "RES 5, C", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xaa, "RES 5, D", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xab, "RES 5, E", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xac, "RES 5, H", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xad, "RES 5, L", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xae, "RES 5, (HL)", 2, NewCycles(16, 16), (*CPU).res_b3_r8},
	{0xaf, "RES 5, A", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xb0, "RES 6, B", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xb1, "RES 6, C", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xb2, "RES 6, D", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xb3, "RES 6, E", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xb4, "RES 6, H", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xb5, "RES 6, L", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xb6, "RES 6, (HL)", 2, NewCycles(16, 16), (*CPU).res_b3_r8},
	{0xb7, "RES 6, A", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xb8, "RES 7, B", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xb9, "RES 7, C", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xba, "RES 7, D", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xbb, "RES 7, E", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xbc, "RES 7, H", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xbd, "RES 7, L", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xbe, "RES 7, (HL)", 2, NewCycles(16, 16), (*CPU).res_b3_r8},
	{0xbf, "RES 7, A", 2, NewCycles(8, 8), (*CPU).res_b3_r8},
	{0xc0, "SET 0, B", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xc1, "SET 0, C", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xc2, "SET 0, D", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xc3, "SET 0, E", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xc4, "SET 0, H", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xc5, "SET 0, L", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xc6, "SET 0, (HL)", 2, NewCycles(16, 16), (*CPU).set_b3_r8},
	{0xc7, "SET 0, A", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xc8, "SET 1, B", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xc9, "SET 1, C", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xca, "SET 1, D", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xcb, "SET 1, E", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xcc, "SET 1, H", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xcd, "SET 1, L", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xce, "SET 1, (HL)", 2, NewCycles(16, 16), (*CPU).set_b3_r8},
	{0xcf, "SET 1, A", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xd0, "SET 2, B", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xd1, "SET 2, C", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xd2, "SET 2, D", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xd3, "SET 2, E", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xd4, "SET 2, H", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xd5, "SET 2, L", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xd6, "SET 2, (HL)", 2, NewCycles(16, 16), (*CPU).set_b3_r8},
	{0xd7, "SET 2, A", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xd8, "SET 3, B", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xd9, "SET 3, C", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xda, "SET 3, D", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xdb, "SET 3, E", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xdc, "SET 3, H", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xdd, "SET 3, L", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xde, "SET 3, (HL)", 2, NewCycles(16, 16), (*CPU).set_b3_r8},
	{0xdf, "SET 3, A", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xe0, "SET 4, B", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xe1, "SET 4, C", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xe2, "SET 4, D", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xe3, "SET 4, E", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xe4, "SET 4, H", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xe5, "SET 4, L", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xe6, "SET 4, (HL)", 2, NewCycles(16, 16), (*CPU).set_b3_r8},
	{0xe7, "SET 4, A", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xe8, "SET 5, B", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xe9, "SET 5, C", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xea, "SET 5, D", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xeb, "SET 5, E", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xec, "SET 5, H", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xed, "SET 5, L", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xee, "SET 5, (HL)", 2, NewCycles(16, 16), (*CPU).set_b3_r8},
	{0xef, "SET 5, A", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xf0, "SET 6, B", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xf1, "SET 6, C", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xf2, "SET 6, D", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xf3, "SET 6, E", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xf4, "SET 6, H", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xf5, "SET 6, L", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xf6, "SET 6, (HL)", 2, NewCycles(16, 16), (*CPU).set_b3_r8},
	{0xf7, "SET 6, A", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xf8, "SET 7, B", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xf9, "SET 7, C", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xfa, "SET 7, D", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xfb, "SET 7, E", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xfc, "SET 7, H", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xfd, "SET 7, L", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
	{0xfe, "SET 7, (HL)", 2, NewCycles(16, 16), (*CPU).set_b3_r8},
	{0xff, "SET 7, A", 2, NewCycles(8, 8), (*CPU).set_b3_r8},
}

func (c *CPU) invalid_instruction() {
	log.Fatal("Calling invalid instruction")
}

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

// 0x Prefixed instructions
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
