package main

import (
	"github.com/vanders/pet/mos6502"
)

type RAM struct {
	Base mos6502.Word // Base address
	Size mos6502.Word // Size

	mem []mos6502.Byte
}

func (r *RAM) GetBase() mos6502.Word {
	return r.Base
}

func (r *RAM) GetSize() mos6502.Word {
	return r.Size
}

func (r *RAM) CheckInterrupt() bool {
	return false
}

func (r *RAM) Reset() {
	r.mem = make([]mos6502.Byte, r.Size)

	for n := mos6502.Word(0); n < r.Size; n++ {
		r.mem[n] = 0x00
	}
}

func (r *RAM) Read(address mos6502.Word) mos6502.Byte {
	return r.mem[address-r.Base]
}

func (r *RAM) Write(address mos6502.Word, data mos6502.Byte) {
	r.mem[address-r.Base] = data
}
