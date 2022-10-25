package main

import (
	"fmt"
	"os"

	"github.com/vanders/pet/mos6502"
)

func dumpAndExit(cpu *mos6502.CPU, err error) {
	fmt.Println(err)
	cpu.Dump()
	os.Exit(1)
}

func main() {
	var mem Memory
	mem.Reset()

	// Load ROMs
	basicLo := &ROM{
		Base: 0xc000,
		Size: 0x1000, // 4k
	}
	basicLo.Reset()
	basicLo.Load("roms/basic-2-c000.901465-01.bin")
	mem.Map(basicLo)

	basicHi := &ROM{
		Base: 0xd000,
		Size: 0x1000, // 4k
	}
	basicHi.Reset()
	basicHi.Load("roms/basic-2-d000.901465-02.bin")
	mem.Map(basicHi)

	edit := &ROM{
		Base: 0xe000,
		Size: 0x800, // 2k
	}
	edit.Reset()
	edit.Load("roms/edit-2-n.901447-24.bin")
	mem.Map(edit)

	kernal := &ROM{
		Base: 0xf000,
		Size: 0x1000, // 4k
	}
	kernal.Reset()
	kernal.Load("roms/kernal-2.901465-03.bin")
	mem.Map(kernal)

	// Create PIAs & VIA
	var pia *PIA

	pia1 := &PIA1{}
	pia = &PIA{
		Base: 0xe810,
	}
	pia.PortRead = pia1.PortRead
	pia.PortWrite = pia1.PortWrite
	mem.Map(pia)

	pia2 := &PIA2{}
	pia = &PIA{
		Base: 0xe820,
	}
	pia.PortRead = pia2.PortRead
	pia.PortWrite = pia2.PortWrite
	mem.Map(pia)

	via := &VIA{
		Base: 0xe840,
	}
	mem.Map(via)

	cpu := mos6502.CPU{
		Read:   mem.Read,
		Write:  mem.Write,
		Writer: os.Stderr,
	}
	cpu.Reset()

	// Execute as many instructions as possible
	for {
		err := cpu.Step()
		if err != nil {
			dumpAndExit(&cpu, fmt.Errorf("execution stopped: %s", err))
		}
	}
}

// Pheripheral Interface Adaptor #1
type PIA1 struct {
	ports [4]mos6502.Byte // 4 8bit ports
}

func (p *PIA1) PortRead(port int) mos6502.Byte {
	switch port {
	case 0:
		/* DICcKKKK
		D=Diagnostic sense
		I=IEEE EOI in
		C=Cassette sense #2
		c=Cassette sense #1
		K=Keyboard Row select
		*/
		return p.ports[port] | mos6502.BIT_7 // Diagnostic Sense is always high
	}
	return p.ports[port]
}

func (p *PIA1) PortWrite(port int, data mos6502.Byte) {
	p.ports[port] = data
}

// Pheripheral Interface Adaptor #2
type PIA2 struct {
	ports [4]mos6502.Byte // 4 8bit ports
}

func (p *PIA2) PortRead(port int) mos6502.Byte {
	return p.ports[port]
}

func (p *PIA2) PortWrite(port int, data mos6502.Byte) {
	p.ports[port] = data
}
