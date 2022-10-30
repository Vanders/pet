package mos6502

import (
	"fmt"
)

// ByteRegister defines an 8 bit register
type ByteRegister struct {
	v Byte
}

func (r *ByteRegister) Set(v Byte) {
	r.v = v
}

func (r *ByteRegister) Get() Byte {
	return r.v
}

func (r *ByteRegister) Inc() {
	r.v += 1
}

func (r *ByteRegister) Dec() {
	r.v -= 1
}

func (r ByteRegister) String() string {
	return fmt.Sprintf("$%02x", r.v)
}

// WordRegister defines a 16 bit register
type WordRegister struct {
	v Word
}

func (r *WordRegister) Set(v Word) {
	r.v = v
}

func (r *WordRegister) Get() Word {
	return r.v
}

func (r *WordRegister) Inc() {
	r.v += 1
}

func (r *WordRegister) Dec() {
	r.v -= 1
}

func (r WordRegister) String() string {
	return fmt.Sprintf("$%04x", r.v)
}

// Flags is an 8 bit mask of CPU states
type Flags struct {
	C, Z, I, D, B, V, N bool
}

// Reset sets the initial state of the flags at CPU reset
func (f *Flags) Reset() {
	f.C = false
	f.Z = false
	f.I = true
	f.D = false
	f.B = true
	f.V = false
	f.N = false
}

const (
	BIT_0 = 1 << 0
	BIT_1 = 1 << 1
	BIT_2 = 1 << 2
	BIT_3 = 1 << 3
	BIT_4 = 1 << 4
	BIT_5 = 1 << 5
	BIT_6 = 1 << 6
	BIT_7 = 1 << 7
)

// Update sets the appropriate flags
func (f *Flags) Update(data Byte) {
	f.Z = (data == 0)
	f.N = (data&BIT_7 != 0)
}

// SetCarry sets or clears the carry flag
func (f *Flags) SetCarry(b bool) {
	f.C = b
}

// SetZero sets or clears the zero flag
func (f *Flags) SetZero(b bool) {
	f.Z = b
}

// SetInterrupt sets or clears the interrupt disable flag
func (f *Flags) SetInterrupt(b bool) {
	f.I = b
}

// SetOverflow sets or clears the overflow flag
func (f *Flags) SetOverflow(b bool) {
	f.V = b
}

// SetNegative sets or clears the negative flag
func (f *Flags) SetNegative(b bool) {
	f.N = b
}

// GetByte returns the flags register as a single 8bit byte
func (f *Flags) GetByte() Byte {
	var status Byte
	if f.C == true {
		status |= BIT_0
	}
	if f.Z == true {
		status |= BIT_1
	}
	if f.I == true {
		status |= BIT_2
	}
	if f.D == true {
		status |= BIT_3
	}
	if f.V == true {
		status |= BIT_6
	}
	if f.N == true {
		status |= BIT_7
	}
	return status
}

// SetByte sets the flags register from a single 8bit byte
func (f *Flags) SetByte(b Byte) {
	f.C = (b&BIT_0 != 0)
	f.Z = (b&BIT_1 != 0)
	f.I = (b&BIT_2 != 0)
	f.D = (b&BIT_3 != 0)
	f.V = (b&BIT_6 != 0)
	f.N = (b&BIT_7 != 0)
}

func (f Flags) String() string {
	return fmt.Sprintf("\tC: %t\n\tZ: %t\n\tI: %t\n\tD: %t\n\tB: %t\n\tV: %t\n\tN: %t",
		f.C,
		f.Z,
		f.I,
		f.D,
		f.B,
		f.V,
		f.N)
}
