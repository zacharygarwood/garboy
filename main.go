package main

func main() {
	cartridge := NewCartridge("./test_roms/04-op r,imm.gb", 0x2000)
	ppu := NewPPU()
	cpu := NewCPU(cartridge, ppu)

	cpu.SkipBootROM()

	maxCycles := 15_000_000
	totalCycles := 0

	for totalCycles < maxCycles {
		cpu.PrintState()
		cycles := cpu.Step()
		totalCycles += int(cycles)
	}
}
