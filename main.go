package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/vanders/pet/mos6502"
)

func dumpAndExit(cpu *mos6502.CPU, err error) {
	fmt.Println(err)
	cpu.Dump()
	os.Exit(1)
}

func main() {
	var wg sync.WaitGroup

	// Create a new memory bus
	bus := Bus{}

	// Initialise memory
	ram := &RAM{}
	ram.Reset()
	bus.Map(ram)

	// Load ROMs
	basicLo := &ROM{
		Base: 0xc000,
		Size: 0x1000, // 4k
	}
	basicLo.Reset()
	basicLo.Load("roms/basic-2-c000.901465-01.bin")
	bus.Map(basicLo)

	basicHi := &ROM{
		Base: 0xd000,
		Size: 0x1000, // 4k
	}
	basicHi.Reset()
	basicHi.Load("roms/basic-2-d000.901465-02.bin")
	bus.Map(basicHi)

	edit := &ROM{
		Base: 0xe000,
		Size: 0x800, // 2k
	}
	edit.Reset()
	edit.Load("roms/edit-2-n.901447-24.bin")
	bus.Map(edit)

	kernal := &ROM{
		Base: 0xf000,
		Size: 0x1000, // 4k
	}
	kernal.Reset()
	kernal.Load("roms/kernal-2.901465-03.bin")
	bus.Map(kernal)

	// Configure keyboard
	buf := make(chan Key, 1)
	kbd := Keyboard{
		Buffer: buf,
	}
	kbd.Reset()

	wg.Add(1)
	go func() {
		defer wg.Done()
		kbd.Scan()
	}()

	// Create PIAs & VIA
	var pia *PIA

	pia1 := &PIA1{
		KbdBuffer: buf,
	}
	pia = &PIA{
		Base: 0xe810,
	}
	pia.PortRead = pia1.PortRead
	pia.PortWrite = pia1.PortWrite
	pia.IRQ = pia1.IRQ
	bus.Map(pia)

	pia2 := &PIA2{}
	pia = &PIA{
		Base: 0xe820,
	}
	pia.PortRead = pia2.PortRead
	pia.PortWrite = pia2.PortWrite
	pia.IRQ = pia2.IRQ
	bus.Map(pia)

	via := &VIA{
		Base: 0xe840,
	}
	bus.Map(via)

	cpu := mos6502.CPU{
		Read:   bus.Read,
		Write:  bus.Write,
		Writer: os.Stderr,
	}
	cpu.Reset()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Execute as many instructions as possible
		for {
			err := cpu.Step()
			if err != nil {
				dumpAndExit(&cpu, fmt.Errorf("\nexecution stopped: %s", err))
			}

			// Check devices for interrupts
			if bus.CheckInterrupts() {
				cpu.Interrupt()
			}
		}
	}()

	wg.Wait()
}

// Pheripheral Interface Adaptor #1
type PIA1 struct {
	ports [4]mos6502.Byte // 4 8bit ports

	KbdBuffer chan (Key) // Keyboard "buffer"
	key       Key        // Last keypress
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
	case 2:
		// KKKKKKKK	K=Keyboard Row Input
		// get keyboard scan row (bits 0-3)
		row := p.ports[0] & 0x0f
		// startup can sometimes set the row to 0x0f
		if row > 9 {
			return mos6502.Byte(0xff)
		}
		// does the row being scanned have a keypress?
		if row == mos6502.Byte(p.key.row) {
			// return the key bit
			return mos6502.Byte(0xff - (0x01 << p.key.bit))
		} else {
			// nothing here
			return mos6502.Byte(0xff)
		}
	}
	return p.ports[port]
}

func (p *PIA1) PortWrite(port int, data mos6502.Byte) {
	p.ports[port] = data
}

func (p *PIA1) IRQ() bool {
	select {
	case key := <-p.KbdBuffer:
		// got a key
		p.key = key
		return true
	default:
		return false
	}
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

func (p *PIA2) IRQ() bool {
	return false
}
