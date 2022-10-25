package main

import (
	"github.com/vanders/pet/mos6502"
)

// VIA models a Versatile Interface Adaptor
type VIA struct {
	Base mos6502.Word // Base address

	portAOut    mos6502.Byte
	portBOut    mos6502.Byte
	portADir    mos6502.Byte
	portBDir    mos6502.Byte
	timer1      mos6502.Word
	timer1latch mos6502.Word
	timer2      mos6502.Word
	timer2latch mos6502.Word
	peripheral  mos6502.Byte
	ifr         mos6502.Byte
	ie          mos6502.Byte
}

func (v *VIA) GetBase() mos6502.Word {
	return v.Base
}

func (v *VIA) GetSize() mos6502.Word {
	return mos6502.Word(16)
}

func (v *VIA) Read(address mos6502.Word) mos6502.Byte {
	port := address - v.Base
	switch port {
	case 0x0: // Port B output
		return v.portBOut
	case 0x1: // Port A output
		return v.portAOut
	case 0x2: // Port B direction
		return v.portBDir
	case 0x3: // Port A direction
		return v.portADir
	case 0x4: // Timer 1 low
	case 0x5: // Timer 1 high
	case 0x6: // Timer 1 latch low
	case 0x7: // Timer 1 latch high
	case 0x8: // Timer 2 low
	case 0x9: // Timer 2 high
	case 0xa: // Timer 2 latch low
	case 0xb: // Timer 2 latch high
	case 0xc: // Peripheral control
		return v.peripheral
	case 0xd: // Interrupt flag register (IFR)
		return v.ifr
	case 0xe: // Interrupt enable register
		return v.ie
	case 0xf: // IO Port A output, without handshaking
		return v.portAOut
	default:
		return mos6502.Byte(0)
	}

	return mos6502.Byte(0)
}

func (v *VIA) Write(address mos6502.Word, data mos6502.Byte) {
	port := address - v.Base
	switch port {
	case 0x0: // Port B output
		v.portBOut = data
	case 0x1: // Port A output
		v.portAOut = data
	case 0x2: // Port B direction
		v.portBDir = data
	case 0x3: // Port A direction
		v.portADir = data
	case 0x4: // Timer 1 low
	case 0x5: // Timer 1 high
	case 0x6: // Timer 1 latch low
	case 0x7: // Timer 1 latch high
	case 0x8: // Timer 2 low
	case 0x9: // Timer 2 high
	case 0xa: // Timer 2 latch low
	case 0xb: // Timer 2 latch high
	case 0xc: // Peripheral control
		v.peripheral = data
	case 0xd: // Interrupt flag register (IFR)
		v.ifr = data
	case 0xe: // Interrupt enable register
		v.ie = data
	case 0xf: // IO Port A output, without handshaking
		v.portAOut = data
	}
}
