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

// Pan Docs reference: https://gbdev.io/pandocs/CPU_Instruction_Set.html

// Block 0
func (c *CPU) nop() {
	// TODO
}

func (c *CPU) ld_r16_imm16(dst Register16) {
	// TODO
}

func (c *CPU) ld_r16mem_a(dst *Register16) {
	// TODO
}

func (c *CPU) ld_a_r16mem(src *Register16) {
	// TODO
}

func (c *CPU) ld_imm16_sp() {
	// TODO
}

func (c *CPU) inc_r16(operand Register16) {
	// TODO
}

func (c *CPU) dec_r16(operand Register16) {
	// TODO
}

func (c *CPU) add_hl_r16(operand Register16) {
	// TODO
}

func (c *CPU) inc_r8(operand Register8) {
	// TODO
}

func (c *CPU) dec_r8(operand Register8) {
	// TODO
}

func (c *CPU) ld_r8_imm8(dst Register8) {
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

func (c *CPU) jr_cond_imm8(cond bool) {
	// TODO
}

func (c *CPU) stop() {
	// TODO
}

// Block 1
func (c *CPU) ld_r8_r8(dst Register8, src Register8) {
	// TODO
	// Exception: ld [hl] [hl] yields the halt instruction
}

func (c *CPU) halt() {
	// TODO
}

// Block 2
func (c *CPU) add_a_r8(operand Register8) {
	// TODO
}

func (c *CPU) adc_a_r8(operand Register8) {
	// TODO
}

func (c *CPU) sub_a_r8(operand Register8) {
	// TODO
}

func (c *CPU) sbc_a_r8(operand Register8) {
	// TODO
}

func (c *CPU) and_a_r8(operand Register8) {
	// TODO
}

func (c *CPU) xor_a_r8(operand Register8) {
	// TODO
}

func (c *CPU) or_a_r8(operand Register8) {
	// TODO
}

func (c *CPU) cp_a_r8(operand Register8) {
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

func (c *CPU) ret_cond(cond bool) {
	// TODO
}

func (c *CPU) ret() {
	// TODO
}

func (c *CPU) reti() {
	// TODO
}

func (c *CPU) jp_cond_imm16(cond bool) {
	// TODO
}

func (c *CPU) jp_imm16() {
	// TODO
}

func (c *CPU) jp_hl() {
	// TODO
}

func (c *CPU) call_cond_imm16(cond bool) {
	// TODO
}

func (c *CPU) call_imm16() {
	// TODO
}

func (c *CPU) rst_tgt3(target uint16) {
	// TODO
}

func (c *CPU) pop_r16stk(register Register16) {
	// TODO
}

func (c *CPU) push_r16stk(register Register16) {
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
func (c *CPU) rlc_r8(operand Register8) {
	// TODO
}

func (c *CPU) rrc_r8(operand Register8) {
	// TODO
}

func (c *CPU) rl_r8(operand Register8) {
	// TODO
}

func (c *CPU) rr_r8(operand Register8) {
	// TODO
}

func (c *CPU) sla_r8(operand Register8) {
	// TODO
}

func (c *CPU) sra_r8(operand Register8) {
	// TODO
}

func (c *CPU) swap_r8(operand Register8) {
	// TODO
}

func (c *CPU) srl_r8(operand Register8) {
	// TODO
}

func (c *CPU) bit_b3_r8(bitIndex uint16, operand Register8) {
	// TODO
}

func (c *CPU) res_b3_r8(bitIndex uint16, operand Register8) {
	// TODO
}

func (c *CPU) set_b3_r8(bitIndex uint16, operand Register8) {
	// TODO
}
