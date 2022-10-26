package main

import (
	"github.com/vanders/pet/mos6502"
)

type RAM struct {
	mem [mos6502.MAX_ADDR]mos6502.Byte
}

func (r *RAM) GetBase() mos6502.Word {
	return mos6502.Word(0)
}

func (r *RAM) GetSize() mos6502.Word {
	return mos6502.Word(mos6502.MAX_ADDR - 1)
}

func (r *RAM) CheckInterrupt() bool {
	return false
}

func (r *RAM) Reset() {
	for n := mos6502.STACK_TOP; n < mos6502.MAX_ADDR; n++ {
		r.mem[n] = 0x00
	}
}

func (r *RAM) Read(address mos6502.Word) mos6502.Byte {
	return r.mem[address]
}

func (r *RAM) Write(address mos6502.Word, data mos6502.Byte) {
	r.mem[address] = data
}
