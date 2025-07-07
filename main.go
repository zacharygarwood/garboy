package main

import (
	"time"
)

var (
	MCyclesPerSecond = 1048576.0
	Fps              = 59.7
	CyclesPerFrame   = int(MCyclesPerSecond / Fps)
	TimePerFrame     = time.Second / time.Duration(Fps)
)

func main() {
	cartridge := NewCartridge("./test_roms/02-interrupts.gb", 0x2000)

	interrupts := NewInterrupts()
	ppu := NewPPU(interrupts)
	display := NewDisplay(ppu)
	timer := NewTimer(interrupts)
	mmu := NewMMU(cartridge, ppu, timer, interrupts)
	cpu := NewCPU(mmu, interrupts)

	scheduler := NewScheduler(cpu, ppu, timer)

	//cpu.SkipBootROM()

	go RunDisplay(display)

	for {
		frameStartTime := time.Now()
		for cyclesThisFrame := 0; cyclesThisFrame < CyclesPerFrame; cyclesThisFrame++ {
			scheduler.Step()
		}
		elapsedTime := time.Since(frameStartTime)

		if elapsedTime < TimePerFrame {
			time.Sleep(TimePerFrame - elapsedTime)
		}
	}
}
