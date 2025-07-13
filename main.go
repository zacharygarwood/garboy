package main

import (
	"time"
)

var (
	CyclesPerSecond = 4194304.0
	Fps             = 59.7
	CyclesPerFrame  = int(CyclesPerSecond / Fps)
	TimePerFrame    = time.Second / time.Duration(Fps)
)

func main() {
	cartridge := NewCartridge("./roms/pokemon-red.gb")

	interrupts := NewInterrupts()
	ppu := NewPPU(interrupts)
	joypad := NewJoypad()
	display := NewDisplay(ppu, joypad)
	timer := NewTimer(interrupts)
	mmu := NewMMU(cartridge, ppu, timer, joypad, interrupts)
	cpu := NewCPU(mmu, interrupts)

	scheduler := NewScheduler(cpu, ppu, timer)

	// Uncomment this if you aren't patient ;)
	// cpu.SkipBootROM()

	go RunDisplay(display)

	for {
		frameStartTime := time.Now()
		for cyclesThisFrame := 0; cyclesThisFrame < CyclesPerFrame; {
			cyclesThisFrame += int(scheduler.Step())
		}
		elapsedTime := time.Since(frameStartTime)

		if elapsedTime < TimePerFrame {
			time.Sleep(TimePerFrame - elapsedTime)
		}
	}
}
