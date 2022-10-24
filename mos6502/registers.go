package mos6502

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
