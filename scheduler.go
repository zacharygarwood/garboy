package main

type Scheduler struct {
	cpu   *CPU
	ppu   *PPU
	timer *Timer
}

func NewScheduler(cpu *CPU, ppu *PPU, timer *Timer) *Scheduler {
	return &Scheduler{
		cpu:   cpu,
		ppu:   ppu,
		timer: timer,
	}
}

func (s *Scheduler) Step() uint16 {
	cycles := s.cpu.Step()
	s.timer.Step(cycles)
	s.ppu.Step(cycles)
	return cycles
}
