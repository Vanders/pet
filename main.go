package main

import (
	"fmt"
	"os"

	"github.com/vanders/pet/mos6502"
	"github.com/veandco/go-sdl2/sdl"
)

type (
	Byte = mos6502.Byte
	Word = mos6502.Word
)

func dumpAndExit(cpu *mos6502.CPU, ram *RAM, err error) {
	fmt.Println(err)
	cpu.Dump()

	// Dump the Zero Page
	for n := 0; n < 256; n = n + 4 {
		fmt.Fprintf(cpu.Writer,
			"0x%04x: 0x%02x,\t0x%04x: 0x%02x,\t0x%04x: 0x%02x,\t0x%04x: 0x%02x\n",
			Word(n),
			ram.Read(Word(n)),
			Word(n+1),
			ram.Read(Word(n+1)),
			Word(n+2),
			ram.Read(Word(n+2)),
			Word(n+3),
			ram.Read(Word(n+3)))
	}

	os.Exit(1)
}

func main() {
	// Create a new memory bus
	bus := Bus{}

	// Initialise memory

	// Main memory
	ram := &RAM{
		Base: 0x0000,
		Size: 0x7fff, // 32k
	}
	ram.Reset()
	bus.Map(ram)

	// Screen memory
	sram := &RAM{
		Base: 0x8000,
		Size: 0x1000, // 4k
	}
	sram.Reset()
	bus.Map(sram)

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
	kbd := &Keyboard{
		Buffer: buf,
	}
	kbd.Reset()

	// Create PIAs & VIA
	var pia *PIA

	// PIA1
	pia1 := &PIA1{
		Keyboard:  kbd,
		KbdBuffer: buf,
	}
	pia = &PIA{
		Base: 0xe810,
	}
	pia.PortRead = pia1.PortRead
	pia.PortWrite = pia1.PortWrite
	pia.IRQ = pia1.IRQ
	bus.Map(pia)

	// PIA2
	pia2 := &PIA2{}
	pia = &PIA{
		Base: 0xe820,
	}
	pia.PortRead = pia2.PortRead
	pia.PortWrite = pia2.PortWrite
	pia.IRQ = pia2.IRQ
	bus.Map(pia)

	// VIA
	via := &VIA{
		Base: 0xe840,
	}
	bus.Map(via)

	// Start video
	video := Video{
		Read: bus.Read,
	}
	err := video.Reset()
	if err != nil {
		panic(err)
	}
	defer video.Stop()

	// Initialise the CPU & connect it to the bus
	cpu := mos6502.CPU{
		Read:   bus.Read,
		Write:  bus.Write,
		Writer: nil, // os.Stderr
	}
	cpu.Reset()

	// Execute instructions
	lastTicks := sdl.GetTicks()
	currentTicks := lastTicks
	running := true
	for running {
		err := cpu.Step()
		if err != nil {
			dumpAndExit(&cpu, ram, fmt.Errorf("\nexecution stopped: %s", err))
		}

		// Update GUI
		switch event := video.Event().(type) {
		case EventQuit:
			running = false
			break
		case EventKeypress:
			kbd.Scan(event.Key)
		}

		// Check devices for interrupts
		if bus.CheckInterrupts() {
			cpu.Interrupt()
		}

		// Wait 50ms and then redraw the screen
		currentTicks = sdl.GetTicks()
		if currentTicks-lastTicks > 50 {
			err = video.Redraw()
			if err != nil {
				break
			}
			lastTicks = currentTicks
		}
	}
}

// Pheripheral Interface Adaptor #1
type PIA1 struct {
	ports [4]Byte // 4 8bit ports

	Keyboard  *Keyboard  // Keyboard
	KbdBuffer chan (Key) // Keyboard "buffer"
	key       Key        // Last keypress
}

func (p *PIA1) PortRead(port int) Byte {
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
		return p.Keyboard.Get(row)
	}
	return p.ports[port]
}

func (p *PIA1) PortWrite(port int, data Byte) {
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
	ports [4]Byte // 4 8bit ports
}

func (p *PIA2) PortRead(port int) Byte {
	return p.ports[port]
}

func (p *PIA2) PortWrite(port int, data Byte) {
	p.ports[port] = data
}

func (p *PIA2) IRQ() bool {
	return false
}
