package main

// PIA models a Pheripheral Interface Adaptor
type PIA struct {
	Base Word // Base address

	PortRead  func(p int) Byte
	PortWrite func(p int, data Byte)
	IRQ       func() bool
}

func (p *PIA) GetBase() Word {
	return p.Base
}

func (p *PIA) GetSize() Word {
	return Word(4)
}

func (p *PIA) CheckInterrupt() bool {
	return p.IRQ()
}

func (p *PIA) Read(address Word) Byte {
	port := int(address - p.Base)
	if port >= 0 && port <= 3 {
		return p.PortRead(port)
	}
	return Byte(0)
}

func (p *PIA) Write(address Word, data Byte) {
	port := int(address - p.Base)
	if port >= 0 && port <= 3 {
		p.PortWrite(port, data)
	}
}
