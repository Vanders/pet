package main

import (
	"github.com/vanders/pet/mos6502"
)

// PIA models a Pheripheral Interface Adaptor
type PIA struct {
	Base mos6502.Word // Base address

	PortRead  func(p int) mos6502.Byte
	PortWrite func(p int, data mos6502.Byte)
	IRQ       func() bool
}

func (p *PIA) GetBase() mos6502.Word {
	return p.Base
}

func (p *PIA) GetSize() mos6502.Word {
	return mos6502.Word(4)
}

func (p *PIA) CheckInterrupt() bool {
	return p.IRQ()
}

func (p *PIA) Read(address mos6502.Word) mos6502.Byte {
	port := int(address - p.Base)
	if port >= 0 && port <= 3 {
		return p.PortRead(port)
	}
	return mos6502.Byte(0)
}

func (p *PIA) Write(address mos6502.Word, data mos6502.Byte) {
	port := int(address - p.Base)
	if port >= 0 && port <= 3 {
		p.PortWrite(port, data)
	}
}
