package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"garboy/cpu"
	"garboy/interrupts"
	"garboy/memory"
)

type TestState struct {
	PC  uint16      `json:"pc"`
	SP  uint16      `json:"sp"`
	A   uint8       `json:"a"`
	B   uint8       `json:"b"`
	C   uint8       `json:"c"`
	D   uint8       `json:"d"`
	E   uint8       `json:"e"`
	F   uint8       `json:"f"`
	H   uint8       `json:"h"`
	L   uint8       `json:"l"`
	RAM [][2]uint16 `json:"ram"`
}

type TestCycle struct {
	Address uint16
	Value   uint8
	Type    string
}

func (tc *TestCycle) UnmarshalJSON(data []byte) error {
	var raw []interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) != 3 {
		return fmt.Errorf("expected cycle array to have 3 elements, got %d", len(raw))
	}

	addr, ok := raw[0].(float64)
	if !ok {
		return fmt.Errorf("expected cycle address to be a number")
	}
	val, ok := raw[1].(float64)
	if !ok {
		return fmt.Errorf("expected cycle value to be a number")
	}
	typeStr, ok := raw[2].(string)
	if !ok {
		return fmt.Errorf("expected cycle type to be a string")
	}

	tc.Address = uint16(addr)
	tc.Value = uint8(val)
	tc.Type = typeStr

	return nil
}

type TestCase struct {
	Name         string      `json:"name"`
	InitialState TestState   `json:"initial"`
	FinalState   TestState   `json:"final"`
	Cycles       []TestCycle `json:"cycles"`
}

func setupCPUForTest(t *testing.T, state TestState) (*cpu.CPU, memory.Memory) {
	mmu := &MockMmu{
		ram: memory.NewRAM(0x10000),
	}

	interrupts := interrupts.NewInterrupts()
	cpu := cpu.NewCPU(mmu, interrupts)

	cpu.SkipBootROM()

	a, f, b, c, d, e, h, l, sp, pc := cpu.GetState()

	pc.Write(state.PC)
	sp.Write(state.SP)
	a.Write(state.A)
	b.Write(state.B)
	c.Write(state.C)
	d.Write(state.D)
	e.Write(state.E)
	f.Write(state.F)
	h.Write(state.H)
	l.Write(state.L)

	for _, ramState := range state.RAM {
		mmu.Write(ramState[0], uint8(ramState[1]))
	}

	return cpu, mmu
}

func assertState(t *testing.T, testName string, cpu *cpu.CPU, mmu memory.Memory, expected TestState) {
	a, f, b, c, d, e, h, l, sp, pc := cpu.GetState()
	if pc.Read() != expected.PC {
		t.Errorf("%s: PC mismatch. Got %04X, want %04X\n", testName, pc, expected.PC)
	}
	if sp.Read() != expected.SP {
		t.Errorf("%s: SP mismatch. Got %04X, want %04X\n", testName, sp, expected.SP)
	}
	if a.Read() != expected.A {
		t.Errorf("%s: A mismatch. Got %02X, want %02X\n", testName, a.Read(), expected.A)
	}
	if b.Read() != expected.B {
		t.Errorf("%s: B mismatch. Got %02X, want %02X\n", testName, b.Read(), expected.B)
	}
	if c.Read() != expected.C {
		t.Errorf("%s: C mismatch. Got %02X, want %02X\n", testName, c.Read(), expected.C)
	}
	if d.Read() != expected.D {
		t.Errorf("%s: D mismatch. Got %02X, want %02X\n", testName, d.Read(), expected.D)
	}
	if e.Read() != expected.E {
		t.Errorf("%s: E mismatch. Got %02X, want %02X\n", testName, e.Read(), expected.E)
	}
	if (f.Read() & 0xF0) != (expected.F & 0xF0) {
		t.Errorf("%s: F mismatch. Got %02X, want %02X\n", testName, f.Read(), expected.F)
	}
	if h.Read() != expected.H {
		t.Errorf("%s: H mismatch. Got %02X, want %02X\n", testName, h.Read(), expected.H)
	}
	if l.Read() != expected.L {
		t.Errorf("%s: L mismatch. Got %02X, want %02X\n", testName, l.Read(), expected.L)
	}

	for _, ramState := range expected.RAM {
		addr, expectedVal := ramState[0], uint8(ramState[1])
		actualVal := mmu.Read(addr)

		if actualVal != expectedVal {
			t.Errorf("%s: RAM mismatch at %04X. Got %02X, want %02X\n", testName, addr, actualVal, expectedVal)
		}
	}
}

func runCPUTest(t *testing.T, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", path, err)
	}

	var tests []TestCase
	if err := json.Unmarshal(data, &tests); err != nil {
		t.Fatalf("Failed to parse JSON from %s: %v", path, err)
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cpu, mmu := setupCPUForTest(t, tc.InitialState)

			cpu.Step()

			assertState(t, tc.Name, cpu, mmu, tc.FinalState)
		})
	}
}

func TestCpuInstructions(t *testing.T) {
	testDir := "json_tests/"

	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Skipf("Skipping CPU tests: test directory not found at %s", testDir)
		return
	}

	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			t.Run(info.Name(), func(t *testing.T) {
				runCPUTest(t, path)
			})
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk test directory: %v", err)
	}
}
