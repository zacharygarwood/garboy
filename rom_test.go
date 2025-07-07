package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

	cartridge := NewCartridge(romPath, 0)
	interrupts := NewInterrupts()
	ppu := NewPPU(interrupts)
	timer := NewTimer(interrupts)
	mmu := NewMMU(cartridge, ppu, timer, interrupts)
	cpu := NewCPU(mmu, interrupts)
	scheduler := NewScheduler(cpu, ppu, timer)
	cpu.SkipBootROM()

	const maxCycles = 80_000_000
	for cycles := 0; cycles < maxCycles; cycles++ {
		scheduler.Step()

		if isMooneyeTest {
			pc := cpu.reg.pc.Read()
			if mmu.Read(pc) == 0x18 && mmu.Read(pc+1) == 0xFE {
				b, c, d, e, h, l := cpu.reg.b.Read(), cpu.reg.c.Read(), cpu.reg.d.Read(), cpu.reg.e.Read(), cpu.reg.h.Read(), cpu.reg.l.Read()
				if b == 3 && c == 5 && d == 8 && e == 13 && h == 21 && l == 34 {
					return
				}
				t.Errorf("Mooneye test failed. Register state: B:%d, C:%d, D:%d, E:%d, H:%d, L:%d", b, c, d, e, h, l)
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
