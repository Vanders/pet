package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/sqweek/dialog"
	"github.com/vanders/pet/mos6502"
)

type (
	Byte = mos6502.Byte
	Word = mos6502.Word
)

func dumpAndExit(cpu *mos6502.CPU, ram *RAM, err error) {
	fmt.Println(err)
	dump(cpu, ram)
	os.Exit(1)
}

func dump(cpu *mos6502.CPU, ram *RAM) {
	cpu.Dump()

	if cpu.Writer != nil {
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
	}
}

type PET struct {
	cpu  *mos6502.CPU
	bus  *Bus
	ram  *RAM
	pia1 *PIA1
	pia2 *PIA2
	via  *VIA

	cassette *Cassette

	gui *GUI

	mos6502.WordReadWrite
}

const (
	// Kernal vectors
	VEC_LOAD = 0xffd5
	VEC_SAVE = 0xffd8

	// Trap function selectors
	TRAP_LOAD = 0x01
	TRAP_SAVE = 0x02

	// Zero page
	TXTTAB = 0x28
	VARTAB = 0x2a
)

func main() {
	var (
		writer io.Writer //= os.Stderr
		ctx    context.Context
		wg     sync.WaitGroup
	)
	ctx, cancel := context.WithCancel(context.Background())

	debug := flag.Bool("d", false, "enable CPU dissasembly")
	romVersion := flag.Int("r", 2, "ROM version (2 or 4)")
	ramSize := flag.Int("m", 32, "RAM size (Megabytes)")
	flag.Parse()

	if *debug {
		writer = os.Stderr
	}

	// Create a new memory bus
	bus := Bus{}

	// Initialise memory

	// Main memory
	ram := &RAM{
		Base: 0x0000,
		Size: Word(*ramSize * 1024), // 0x7fff, // 32k
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
	switch *romVersion {
	case 0: // Special case for diagnostic ROMs
		u2 := &ROM{
			Base: 0xf000,
			Size: 0x800, // 2k
		}
		u2.Reset()
		u2.Load("roms/U-2 DIA")
		bus.Map(u2)

		u3 := &ROM{
			Base: 0xf800,
			Size: 0x800, // 2k
		}
		u3.Reset()
		u3.Load("roms/U-3 DIA")
		bus.Map(u3)
	case 2:
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
		kernal.PatchVector(VEC_LOAD, LOADPATCH_v2)
		kernal.PatchVector(VEC_SAVE, SAVEPATCH_v2)
		bus.Map(kernal)
	case 4:
		basic1 := &ROM{
			Base: 0xb000,
			Size: 0x1000, // 4k
		}
		basic1.Reset()
		basic1.Load("roms/basic-4-b000.901465-23.bin")
		bus.Map(basic1)

		basic2 := &ROM{
			Base: 0xc000,
			Size: 0x1000, // 4k
		}
		basic2.Reset()
		basic2.Load("roms/basic-4-c000.901465-20.bin")
		bus.Map(basic2)

		basic3 := &ROM{
			Base: 0xd000,
			Size: 0x1000, // 4k
		}
		basic3.Reset()
		basic3.Load("roms/basic-4-d000.901465-21.bin")
		bus.Map(basic3)

		edit := &ROM{
			Base: 0xe000,
			Size: 0x800, // 2k
		}
		edit.Reset()
		edit.Load("roms/edit-4-40-n-50Hz.901498-01.bin")
		bus.Map(edit)

		kernal := &ROM{
			Base: 0xf000,
			Size: 0x1000, // 4k
		}
		kernal.Reset()
		kernal.Load("roms/kernal-4.901465-22.bin")
		kernal.PatchVector(VEC_LOAD, LOADPATCH_v4)
		kernal.PatchVector(VEC_SAVE, SAVEPATCH_v4)
		bus.Map(kernal)
	default:
		fmt.Fprintf(os.Stderr, "Invalid ROM version %d\n", *romVersion)
		os.Exit(1)
	}

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
		Base:      0xe810,
		PortRead:  pia1.PortRead,
		PortWrite: pia1.PortWrite,
		IRQ:       pia1.IRQ,
	}
	bus.Map(pia)

	// PIA2
	pia2 := &PIA2{}
	pia = &PIA{
		Base:      0xe820,
		PortRead:  pia2.PortRead,
		PortWrite: pia2.PortWrite,
		IRQ:       pia2.IRQ,
	}
	bus.Map(pia)

	// VIA
	via := &VIA{
		Base: 0xe840,
	}
	bus.Map(via)

	// Initialise video

	/* The character ROM is special as it is not mapped to the main memory bus
	like the other ROMS. Hence, it has a size but it's base "address" is 0x0000
	and the video circuitry/routine generates an address directly into the ROM.
	*/
	charROM := &ROM{
		Size: 0x800, // 2k
	}
	charROM.Reset()
	charROM.Load("roms/char-901447-10.bin")

	video := &Video{
		Read:    bus.Read,
		VIA_CB2: via.CB2,
		PIA_CB1: pia1.CB1,
		ROM:     charROM,
	}

	// Configure "cassette"
	cas := &Cassette{}

	// Start GUI
	gui := GUI{
		Video: video,
	}
	err := gui.Init()
	if err != nil {
		panic(err)
	}
	defer gui.Stop()

	// Initialise the CPU & connect it to the bus
	cpu := mos6502.NewCPU(bus.Read, bus.Write, nil, writer)
	cpu.Reset()

	// Create a channel for GUI events
	events := make(chan Event, 10)

	pet := &PET{
		cpu:      cpu,
		bus:      &bus,
		ram:      ram,
		pia1:     pia1,
		pia2:     pia2,
		via:      via,
		cassette: cas,
		gui:      &gui,
	}
	pet.ReadWriter = &bus
	cpu.Trap = pet.HandleTrap

	// Run the CPU & pheripherals
	wg.Add(1)
	go func() {
		defer wg.Done()

		running := true
		for running {
			// Execute a single instruction
			err := cpu.Step()
			if err != nil {
				dumpAndExit(cpu, ram, fmt.Errorf("\nexecution stopped: %s", err))
			}

			// Handle any GUI events
			select {
			case event := <-events:
				switch e := event.(type) {
				case EventQuit:
					running = false
					break
				case EventKeypress:
					kbd.Scan(e.Key)
				}
			default:
			}

			// Check devices for interrupts
			if bus.CheckInterrupts() {
				cpu.Interrupt()
			}
		}

		// Cancel the context
		cancel()
	}()

	// Start the GUI event loop
	gui.EventLoop(ctx, events)

	wg.Wait()
	dump(cpu, ram)
}

// Trap handler
func (p *PET) HandleTrap(selector Byte) {
	switch selector {
	case TRAP_LOAD:
		filename, err := p.gui.LoadDialog("Load program", "PRG files", "prg")
		if err == dialog.Cancelled {
			break
		}

		err = p.cassette.Load(filename)
		if err != nil {
			panic(err)
		}
		addr := p.cassette.Addr()
		size := p.cassette.Size()
		fmt.Printf("load %d bytes to address $%04x\n", size, addr)
		for n := Word(0); n < size; n++ {
			b := p.cassette.FetchByte()
			p.bus.Write(addr+n, b)
		}
		// Set top of BASIC
		p.WriteWord(VARTAB, addr+size+1)
	case TRAP_SAVE:
		// Get start & top of BASIC, calculate size
		txttab := p.ReadWord(TXTTAB)
		vartab := p.ReadWord(VARTAB)
		size := vartab - txttab

		fmt.Printf("txttab: $%04x, vartab: $%04x, size=%d\n", txttab, vartab, size)

		// Copy BASIC data
		data := make([]Byte, size)
		for n := mos6502.Word(0); n < size; n++ {
			data[n] = p.bus.Read(txttab + n)
		}

		filename, err := p.gui.LoadDialog("Save program", "PRG files", "prg")
		if err == dialog.Cancelled {
			break
		}

		err = p.cassette.Save(filename, txttab, size, data)
		if err != nil {
			panic(err)
		}
	}
}

// Pheripheral Interface Adaptor #1
type PIA1 struct {
	ports [4]Byte // 4 8bit ports
	irq   bool    // Interrupt request

	Keyboard  *Keyboard  // Keyboard
	KbdBuffer chan (Key) // Keyboard "buffer"
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
	// Get & clear the current IRQ status
	i := p.irq
	p.irq = false

	select {
	case <-p.KbdBuffer:
		// Handle keypresses immediately
		return true
	default:
		return i
	}
}

func (p *PIA1) CB1(retrace bool) {
	// Set Retrace Interrupt flag
	if retrace {
		p.ports[3] |= 0x80
		p.irq = true
	} else {
		p.ports[3] ^= 0x80
		p.irq = false
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
