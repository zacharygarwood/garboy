package main

import "fmt"

const (
	InterruptCycles = 20
	HaltedCycles    = 4 // Treated as NOP
)

var InterruptSources = []uint16{
	VBlankInterruptSource,
	StatInterruptSource,
	TimerInterruptSource,
	SerialInterruptSource,
	JoypadInterruptSource,
}

type CPU struct {
	reg        *Registers
	mmu        MmuInterface
	interrupts *Interrupts

	branched              bool
	halted                bool
	haltBug               bool
	interruptMasterEnable bool // IME

	imeDelay uint8
}

func NewCPU(mmu MmuInterface, interrupts *Interrupts) *CPU {
	return &CPU{
		reg:                   NewRegisters(),
		mmu:                   mmu,
		interrupts:            interrupts,
		branched:              false,
		halted:                false,
		interruptMasterEnable: true,
	}
}

func (c *CPU) Step() uint16 {
	if c.imeDelay > 0 {
		c.imeDelay--
		if c.imeDelay == 0 {
			c.interruptMasterEnable = true
		}
	}

	if c.handleInterrupts() {
		return InterruptCycles
	}

	if c.halted {
		return HaltedCycles
	}

	//c.PrintState()

	opcode := c.fetch()
	instruction := c.decode(opcode)
	return c.execute(instruction)
}

// Fetches the opcode at PC
func (c *CPU) fetch() byte {
	opcode := c.mmu.Read(c.reg.pc.Read())

	if c.haltBug {
		c.haltBug = false
	} else {
		c.reg.pc.Increment()
	}
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
func (c *CPU) execute(instr Instruction) uint16 {
	c.branched = false // Always set to false before executing
	instr.Execute(&instr, c)

	if c.branched {
		return instr.Cycles.branched
	} else {
		return instr.Cycles.normal
	}
}

func (c *CPU) handleInterrupts() bool {
	ie := c.interrupts.IE()
	iff := c.interrupts.IF()
	pending := ie & iff

	if pending == 0 {
		return false
	}

	c.halted = false

	if !c.interruptMasterEnable {
		return false
	}

	for i := 0; i < len(InterruptSources); i++ {
		interrupt := uint8(i)
		if IsBitSet(pending, interrupt) {
			c.interrupts.Clear(interrupt)
			c.interruptMasterEnable = false

			c.Push16(c.reg.pc.Read())
			c.reg.pc.Write(InterruptSources[i])
			return true
		}
	}
	return false
}

// Pushes val on to the stack (SP)
func (c *CPU) Push16(val uint16) {
	sp := c.reg.sp.Read() - 2
	c.mmu.WriteWord(sp, val)
	c.reg.sp.Write(sp)
}

// Pops val off the stack (SP)
func (c *CPU) Pop16() uint16 {
	sp := c.reg.sp.Read()
	val := c.mmu.ReadWord(sp)
	c.reg.sp.Write(sp + 2)
	return val
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

	c.mmu.SetBootRomEnabled(false)
}

func (c *CPU) PrintState() {
	pc := c.reg.pc.Read()
	fmt.Printf("[CPU] A:%.2X F:%.2X B:%.2X C:%.2X D:%.2X E:%.2X H:%.2X L:%.2X SP:%.4X PC:%.4X PCMEM:%02X,%02X,%02X,%02X\n",
		c.reg.a.Read(), c.reg.f.Read(), c.reg.b.Read(), c.reg.c.Read(), c.reg.d.Read(), c.reg.e.Read(), c.reg.h.Read(),
		c.reg.l.Read(), c.reg.sp.Read(), pc, c.byteAt(pc).Read(), c.byteAt(pc+1).Read(), c.byteAt(pc+2).Read(), c.byteAt(pc+3).Read())
}

func (c *CPU) PrintStateDecimal() {
	pc := c.reg.pc.Read()
	fmt.Printf("A:%d F:%d B:%d C:%d D:%d E:%d H:%d L:%d SP:%d PC:%d PCMEM:%d,%d,%d,%d\n",
		c.reg.a.Read(), c.reg.f.Read(), c.reg.b.Read(), c.reg.c.Read(), c.reg.d.Read(), c.reg.e.Read(), c.reg.h.Read(),
		c.reg.l.Read(), c.reg.sp.Read(), pc, c.byteAt(pc).Read(), c.byteAt(pc+1).Read(), c.byteAt(pc+2).Read(), c.byteAt(pc+3).Read())
}
