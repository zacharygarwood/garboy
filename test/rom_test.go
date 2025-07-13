package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"garboy/cartridge"
	"garboy/cpu"
	"garboy/display"
	"garboy/interrupts"
	"garboy/mmu"
	"garboy/scheduler"
	"garboy/timer"
)

func runRomTest(t *testing.T, romPath string) {
	var r, w, originalStdout *os.File
	var outputBuffer bytes.Buffer

	isBlarggTest := strings.Contains(romPath, "blargg")
	isMooneyeTest := strings.Contains(romPath, "mooneye")

	if isBlarggTest {
		originalStdout = os.Stdout
		r, w, _ = os.Pipe()
		os.Stdout = w
	}

	cartridge := cartridge.NewCartridge(romPath)
	interrupts := interrupts.NewInterrupts()
	ppu := display.NewPPU(interrupts)
	timer := timer.NewTimer(interrupts)
	joypad := display.NewJoypad()
	mmu := mmu.NewMMU(cartridge, ppu, timer, joypad, interrupts)
	cpu := cpu.NewCPU(mmu, interrupts)
	scheduler := scheduler.NewScheduler(cpu, ppu, timer)
	cpu.SkipBootROM()

	const maxCycles = 80_000_000
	for cycles := 0; cycles < maxCycles; cycles++ {
		scheduler.Step()

		_, _, b, c, d, e, h, l, _, pc := cpu.GetState()

		if isMooneyeTest {
			if mmu.Read(pc.Read()) == 0x18 && mmu.Read(pc.Read()+1) == 0xFE {
				if b.Read() == 3 && c.Read() == 5 && d.Read() == 8 && e.Read() == 13 && h.Read() == 21 && l.Read() == 34 {
					return
				}
				t.Errorf("Mooneye test failed. Register state: B:%d, C:%d, D:%d, E:%d, H:%d, L:%d", b.Read(), c.Read(), d.Read(), e.Read(), h.Read(), l.Read())
				return
			}
		}
	}

	if isBlarggTest {
		os.Stdout = originalStdout
		w.Close()

		io.Copy(&outputBuffer, r)

		output := outputBuffer.String()
		if strings.Contains(output, "Passed") {
			return
		}
		t.Errorf("Blargg test failed or timed out. Captured output:\n%s", output)
		return
	}

	t.Errorf("Test timed out after executing max cycles")
}

func TestRoms(t *testing.T) {
	testDirs := []string{
		"./test_roms/blargg",
		"./test_roms/mooneye",
	}

	for _, dir := range testDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Logf("Skipping ROM tests: Directory not found at %s", dir)
			continue
		}

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && filepath.Ext(path) == ".gb" {
				romName := info.Name()
				t.Run(romName, func(t *testing.T) {
					runRomTest(t, path)
				})
			}
			return nil
		})
	}
}
