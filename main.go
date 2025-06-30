package main

func main() {
	cartridge := NewCartridge("./test_roms/02-interrupts.gb", 0x2000)
	interrupts := NewInterrupts()
	ppu := NewPPU()
	timer := NewTimer(interrupts)
	mmu := NewMMU(cartridge, ppu, timer, interrupts)
	cpu := NewCPU(mmu, interrupts)

	cpu.SkipBootROM()

	maxCycles := 80_000_000
	totalCycles := 0

	for totalCycles < maxCycles {
		cpu.PrintState()
		cycles := cpu.Tick()
		timer.Tick(uint16(cycles))
		totalCycles += int(cycles)
	}
}
