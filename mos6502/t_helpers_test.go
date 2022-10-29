package mos6502

import "testing"

// Start of executable code, or data, after the zero page
const (
	exeStart  Word = 0x0200
	dataStart Word = 0x0300
	memMax    Word = 0xffff
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

- An opcode to execute
- A name
- A setup function (initial state)
- A check function (final state)
*/
type testCase struct {
	op    Opcode
	name  string
	setup func(*testing.T, *CPU, *fakeMem)
	check func(*testing.T, *CPU, *fakeMem)
}

type testCases []testCase

func (tests testCases) Run(t *testing.T) {
	m := newMem()
	c := newCPU(m)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setup(t, c, m)
			Call(t, c, testCase.op)
			testCase.check(t, c, m)
		})

		m.Reset()
		c.Reset()
	}
}

// Call the instruction from the given opcode
func Call(t *testing.T, c *CPU, opcode Opcode) {
	ins, ok := c.instructionSet[Opcode(opcode)]
	if !ok {
		t.Errorf("invalid or unknown instruction 0x%2x", opcode)
	}
	err := ins.F(ins)
	if err != nil {
		t.Error(err)
	}
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

func CSet(t *testing.T, c *CPU) {
	C := c.Registers.P.C
	if C != true {
		t.Error("carry flag is not set on carry")
	}
}

func CClear(t *testing.T, c *CPU) {
	C := c.Registers.P.C
	if C != false {
		t.Error("carry flag is set on non-carry")
	}
}

// Helpers for checking CPU registers
func CompareA(t *testing.T, c *CPU, expected Byte) {
	a := c.Registers.A.Get()
	if a != expected {
		t.Errorf("A incorrect: expected 0x%02x, got 0x%02x", expected, a)
	}
}

func CompareX(t *testing.T, c *CPU, expected Byte) {
	x := c.Registers.X.Get()
	if x != expected {
		t.Errorf("X incorrect: expected 0x%02x, got 0x%02x", expected, x)
	}
}

func CompareY(t *testing.T, c *CPU, expected Byte) {
	y := c.Registers.Y.Get()
	if y != expected {
		t.Errorf("Y incorrect: expected 0x%02x, got 0x%02x", expected, y)
	}
}

func CompareSP(t *testing.T, c *CPU, expected Byte) {
	sp := c.Registers.S.Get()
	if sp != expected {
		t.Errorf("SP incorrect: expected 0x%02x, got 0x%02x", expected, sp)
	}
}

// Helpers for checking memory contents
func CompareMem(t *testing.T, m *fakeMem, addr Word, expected Byte) {
	data := m.GetByte(addr)
	if data != expected {
		t.Errorf("addr 0x%04x incorrect: expected 0x%02x, got 0x%02x", addr, expected, data)
	}
}
