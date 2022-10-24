package main

import (
	"io/ioutil"

	"github.com/vanders/pet/mos6502"
)

type Memory struct {
	mem [mos6502.MAX_ADDR]mos6502.Byte
}

func (m *Memory) Reset() {
	for n := mos6502.STACK_TOP; n < mos6502.MAX_ADDR; n++ {
		m.mem[n] = 0x00
	}
}

func (m *Memory) Read(address mos6502.Word) mos6502.Byte {
	return m.mem[address]
}

func (m *Memory) Write(address mos6502.Word, value mos6502.Byte) {
	m.mem[address] = value
}

// XXX This is a hack until we have proper object mapping and this can be replaced with a ROM object
func (m *Memory) Load(address mos6502.Word, filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if address+mos6502.Word(len(data)) > mos6502.MAX_ADDR-1 {
		panic("can't load that there (too big)")
	}
	for n := 0; n < len(data); n++ {
		m.mem[address+mos6502.Word(n)] = mos6502.Byte(data[n])
	}
}
