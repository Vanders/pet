package main

import (
	"io/ioutil"

	"github.com/vanders/pet/mos6502"
)

type ROM struct {
	Base mos6502.Word
	Size mos6502.Word

	mem []mos6502.Byte
}

func (r *ROM) GetBase() mos6502.Word {
	return r.Base
}

func (r *ROM) GetSize() mos6502.Word {
	return r.Size
}

func (r *ROM) CheckInterrupt() bool {
	return false
}

func (r *ROM) Reset() {
	r.mem = make([]mos6502.Byte, r.Size)

	for n := mos6502.Word(0); n < r.Size; n++ {
		r.mem[n] = 0x00
	}
}

func (r *ROM) Read(address mos6502.Word) mos6502.Byte {
	return r.mem[address-r.Base]
}

func (r *ROM) Write(mos6502.Word, mos6502.Byte) {
	// ROM
}

func (r *ROM) Load(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if mos6502.Word(len(data)) > r.Size {
		panic("can't load that there (too big)")
	}
	for n := 0; n < len(data); n++ {
		r.mem[mos6502.Word(n)] = mos6502.Byte(data[n])
	}
}
