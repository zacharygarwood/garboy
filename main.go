package main

import (
	"time"

	"garboy/cartridge"
	"garboy/cpu"
	"garboy/display"
	"garboy/interrupts"
	"garboy/mmu"
	"garboy/scheduler"
	"garboy/timer"
)

var (
	CyclesPerSecond = 4194304.0
	Fps             = 59.7
	CyclesPerFrame  = int(CyclesPerSecond / Fps)
	TimePerFrame    = time.Second / time.Duration(Fps)
)

func main() {
	cartridge := cartridge.NewCartridge("./roms/pokemon-red.gb")

	interrupts := interrupts.NewInterrupts()
	ppu := display.NewPPU(interrupts)
	joypad := display.NewJoypad()
	lcd := display.NewDisplay(ppu, joypad)
	timer := timer.NewTimer(interrupts)
	mmu := mmu.NewMMU(cartridge, ppu, timer, joypad, interrupts)
	cpu := cpu.NewCPU(mmu, interrupts)

	scheduler := scheduler.NewScheduler(cpu, ppu, timer)

	// Uncomment this if you aren't patient ;)
	// cpu.SkipBootROM()

	go display.RunDisplay(lcd)

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
