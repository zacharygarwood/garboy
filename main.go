package main

import "fmt"

func main() {
	cartridge := NewCartridge("./test_roms/cpu_instrs.gb", 0x2000)
	ppu := NewPPU()
	cpu := NewCPU(cartridge, ppu)

	for range 256 {
		if !cpu.Step() {
			return
		}
	}

	fmt.Printf("Boot ROM implemented!")
}
