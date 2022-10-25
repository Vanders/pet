package main

import (
	"github.com/vanders/pet/mos6502"
)

// PIA models a Pheripheral Interface Adaptor
type PIA struct {
	Base  mos6502.Word    // Base address
	Ports [4]mos6502.Byte // 4 8 bit ports
}

func (p *PIA) GetBase() mos6502.Word {
	return p.Base
}

func (p *PIA) GetSize() mos6502.Word {
	return mos6502.Word(4)
}

func (p *PIA) Read(address mos6502.Word) mos6502.Byte {
	port := address - p.Base
	if port >= 0 && port <= 3 {
		return p.Ports[port]
	}
	return mos6502.Byte(0)
}

func (p *PIA) Write(address mos6502.Word, data mos6502.Byte) {
	port := address - p.Base
	if port >= 0 && port <= 3 {
		p.Ports[port] = data
	}
}
