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
	f.B = true
	f.D = false
	f.I = true
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

// SetOverflow sets or clears the overflow flag
func (f *Flags) SetOverflow(b bool) {
	f.V = b
}

// SetNegative sets or clears the negative flag
func (f *Flags) SetNegative(b bool) {
	f.N = b
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
