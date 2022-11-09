package main

// VIA models a Versatile Interface Adaptor
type VIA struct {
	Base Word // Base address

	portAOut    Byte
	portBOut    Byte
	portADir    Byte
	portBDir    Byte
	timer1      Word
	timer1latch Word
	timer2      Word
	timer2latch Word
	peripheral  Byte
	ifr         Byte
	ie          Byte
}

func (v *VIA) GetBase() Word {
	return v.Base
}

func (v *VIA) GetSize() Word {
	return Word(16)
}

func (v *VIA) CheckInterrupt() bool {
	return false
}

func (v *VIA) Read(address Word) Byte {
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
		return Byte(0)
	}

	return Byte(0)
}

func (v *VIA) Write(address Word, data Byte) {
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

// Additional helper functions for exposing specific lines

// Return the current state of the CB2 line (Peripheral Control Register bit 2)
func (v *VIA) CB2() Byte {
	return Byte((v.peripheral & 0x02) >> 1)
}
