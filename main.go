package main

func main() {
	cartridge := NewCartridge("./test_roms/11-op a,(hl).gb", 0x2000)
	ppu := NewPPU()
	cpu := NewCPU(cartridge, ppu)

	cpu.SkipBootROM()

	maxCycles := 80_000_000
	totalCycles := 0

	for totalCycles < maxCycles {
		cpu.PrintState()
		cycles := cpu.Step()
		totalCycles += int(cycles)
	}
}
