package scheduler

import (
	"garboy/cpu"
	"garboy/display"
	"garboy/timer"
)

type Scheduler struct {
	cpu   *cpu.CPU
	ppu   *display.PPU
	timer *timer.Timer
}

func NewScheduler(cpu *cpu.CPU, ppu *display.PPU, timer *timer.Timer) *Scheduler {
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
