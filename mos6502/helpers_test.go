package mos6502

import "testing"

// Start of executable code, or data, after the zero page
const (
	exeStart Word = 0x200
	memMax   Word = 0xffff
)

// "fake" memory that provides a bunch of helper methods
type fakeMem struct {
	mem     [memMax]Byte
	curAddr Word
}

func (m *fakeMem) Reset() {
	for n := Word(0); n < memMax; n++ {
		m.mem[n] = 0x00
	}
	// Start of executable code is above the stack
	m.curAddr = exeStart
	m.SetWord(VEC_RESET, exeStart)
}

func (m *fakeMem) Read(addr Word) Byte {
	return m.mem[addr]
}

func (m *fakeMem) Write(addr Word, data Byte) {
	m.mem[addr] = data
}

func (m *fakeMem) SetByte(addr Word, data Byte) {
	m.Write(addr, data)
}

func (m *fakeMem) SetWord(addr Word, data Word) {
	hi := Byte((data >> 8) & 0xFF)
	lo := Byte(data & 0xFF)

	m.Write(addr, lo)
	m.Write(addr+1, hi)
}

// WriteByte writes a byte to the next address at the counter and increments
// it by 1
func (m *fakeMem) WriteByte(data Byte) {
	m.Write(m.curAddr, data)
	m.curAddr += 1
}

// WriteWord writes a word to the next address at the counter and increments
// it by 2
func (m *fakeMem) WriteWord(data Word) {
	hi := Byte((data >> 8) & 0xFF)
	lo := Byte(data & 0xFF)

	m.Write(m.curAddr, lo)
	m.curAddr += 1
	m.Write(m.curAddr, hi)
	m.curAddr += 1
}

func (m *fakeMem) GetByte(addr Word) Byte {
	return m.Read(addr)
}

func (m *fakeMem) GetWord(addr Word) Word {
	lo := m.Read(addr)
	hi := m.Read(addr + 1)
	return Word(hi)<<8 | Word(lo)
}

func (m *fakeMem) SetCurrentAddress(addr Word) {
	m.curAddr = addr
}

func newMem() *fakeMem {
	m := &fakeMem{}
	m.Reset()

	return m
}

func newCPU(m *fakeMem) *CPU {
	c := &CPU{
		Read:  m.Read,
		Write: m.Write,
	}
	c.Reset()
	c.PC.Set(exeStart)

	return c
}

/*
	Each test must define the following:

- A name
- An instruction to execute (just the address mode is required)
- A setup function (initial state)
- A check function (final state)
*/
type testCase struct {
	name  string
	ins   Instruction
	setup func(*testing.T, *CPU, *fakeMem)
	check func(*testing.T, *CPU, *fakeMem)
}

// Helpers for checking various CPU flags
func NSet(t *testing.T, c *CPU) {
	N := c.Registers.P.N
	if N != true {
		t.Error("negative flag is not set on negative result")
	}
}

func NClear(t *testing.T, c *CPU) {
	N := c.Registers.P.N
	if N != false {
		t.Error("negative flag is set on non-negative result")
	}
}

func ZSet(t *testing.T, c *CPU) {
	Z := c.Registers.P.Z
	if Z != true {
		t.Error("zero flag is not set on zero result")
	}
}

func ZClear(t *testing.T, c *CPU) {
	Z := c.Registers.P.Z
	if Z != false {
		t.Error("zero flag is set on non-zero result")
	}
}

// Helpers for checking CPU registers
func CompareA(t *testing.T, c *CPU, expected Byte) {
	a := c.Registers.A.Get()
	if a != expected {
		t.Errorf("A incorrect: expected 0x%02x, got 0x%02x", expected, a)
	}
}
