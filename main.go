package main

func main() {
	cartridge := NewCartridge("./test_roms/07-jr,jp,call,ret,rst.gb", 0x2000)
	ppu := NewPPU()
	cpu := NewCPU(cartridge, ppu)

	cpu.SkipBootROM()

	maxCycles := 20_000_000
	totalCycles := 0

	for totalCycles < maxCycles {
		cpu.PrintState()
		cycles := cpu.Step()
		totalCycles += int(cycles)
	}
}
